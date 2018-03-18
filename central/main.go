package main

import (
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"

	_ "bitbucket.org/stack-rox/apollo/pkg/authproviders/all"
	_ "bitbucket.org/stack-rox/apollo/pkg/notifications/notifiers/all"
	_ "bitbucket.org/stack-rox/apollo/pkg/registries/all"
	_ "bitbucket.org/stack-rox/apollo/pkg/scanners/all"

	clustersZip "bitbucket.org/stack-rox/apollo/central/clusters/zip"
	"bitbucket.org/stack-rox/apollo/central/db"
	"bitbucket.org/stack-rox/apollo/central/db/boltdb"
	"bitbucket.org/stack-rox/apollo/central/db/inmem"
	"bitbucket.org/stack-rox/apollo/central/detection"
	"bitbucket.org/stack-rox/apollo/central/notifications"
	"bitbucket.org/stack-rox/apollo/central/search"
	"bitbucket.org/stack-rox/apollo/central/search/blevesearch"
	"bitbucket.org/stack-rox/apollo/central/service"
	"bitbucket.org/stack-rox/apollo/pkg/env"
	pkgGRPC "bitbucket.org/stack-rox/apollo/pkg/grpc"
	authnUser "bitbucket.org/stack-rox/apollo/pkg/grpc/authn/user"
	"bitbucket.org/stack-rox/apollo/pkg/grpc/authz/allow"
	authzUser "bitbucket.org/stack-rox/apollo/pkg/grpc/authz/user"
	"bitbucket.org/stack-rox/apollo/pkg/grpc/clusters"
	"bitbucket.org/stack-rox/apollo/pkg/grpc/routes"
	"bitbucket.org/stack-rox/apollo/pkg/logging"
	"bitbucket.org/stack-rox/apollo/pkg/mtls/verifier"
	"bitbucket.org/stack-rox/apollo/pkg/ui"
	"google.golang.org/grpc"
)

var (
	log = logging.LoggerForModule()
)

func main() {
	central := newCentral()

	indexer, err := blevesearch.NewIndexer()
	if err != nil {
		panic(err)
	}
	central.indexer = indexer

	persistence, err := boltdb.NewWithDefaults(env.DBPath.Setting(), indexer)
	if err != nil {
		panic(err)
	}
	central.database = inmem.New(persistence)

	central.notificationProcessor, err = notifications.NewNotificationProcessor(central.database)
	if err != nil {
		panic(err)
	}
	go central.notificationProcessor.Start()
	central.detector, err = detection.New(central.database, central.notificationProcessor)
	if err != nil {
		panic(err)
	}

	go central.startGRPCServer()

	central.processForever()
}

type central struct {
	signalsC              chan os.Signal
	detector              *detection.Detector
	notificationProcessor *notifications.Processor
	database              db.Storage
	indexer               search.Indexer
	server                pkgGRPC.API
}

func newCentral() *central {
	central := &central{}

	central.signalsC = make(chan os.Signal, 1)
	signal.Notify(central.signalsC, os.Interrupt)
	signal.Notify(central.signalsC, syscall.SIGINT, syscall.SIGTERM)

	return central
}

func (c *central) startGRPCServer() {
	idService := service.NewServiceIdentityService(c.database)
	clusterService := service.NewClusterService(c.database)
	clusterWatcher := clusters.NewClusterWatcher(c.database)
	userAuth := authnUser.NewAuthInterceptor(c.database)

	config := pkgGRPC.Config{
		CustomRoutes: c.customRoutes(userAuth, clusterService, idService),
		TLS:          verifier.CA{},
		UnaryInterceptors: []grpc.UnaryServerInterceptor{
			userAuth.UnaryInterceptor(),
			clusterWatcher.UnaryInterceptor(),
		},
		StreamInterceptors: []grpc.StreamServerInterceptor{
			userAuth.StreamInterceptor(),
			clusterWatcher.StreamInterceptor(),
		},
	}

	c.server = pkgGRPC.NewAPI(config)
	c.server.Register(service.NewAlertService(c.database))
	c.server.Register(service.NewAuthService())
	c.server.Register(service.NewAuthProviderService(c.database, userAuth))
	c.server.Register(service.NewBenchmarkService(c.database))
	c.server.Register(service.NewBenchmarkScansService(c.database))
	c.server.Register(service.NewBenchmarkScheduleService(c.database))
	c.server.Register(service.NewBenchmarkResultsService(c.database, c.notificationProcessor))
	c.server.Register(service.NewBenchmarkTriggerService(c.database))
	c.server.Register(clusterService)
	c.server.Register(service.NewDeploymentService(c.database))
	c.server.Register(service.NewImageService(c.database))
	c.server.Register(service.NewNotifierService(c.database, c.notificationProcessor, c.detector))
	c.server.Register(service.NewPingService())
	c.server.Register(service.NewPolicyService(c.database, c.detector))
	c.server.Register(service.NewRegistryService(c.database, c.detector))
	c.server.Register(service.NewScannerService(c.database, c.detector))
	c.server.Register(service.NewSearchService(c.indexer))
	c.server.Register(idService)
	c.server.Register(service.NewSensorEventService(c.detector, c.database))
	c.server.Start()
}

func (c *central) customRoutes(userAuth *authnUser.AuthInterceptor, clusterService *service.ClusterService, idService *service.IdentityService) (routeMap map[string]routes.CustomRoute) {
	routeMap = map[string]routes.CustomRoute{
		"/": {
			AuthInterceptor: userAuth.HTTPInterceptor,
			Authorizer:      allow.Anonymous(),
			ServerHandler:   ui.Mux(),
			Compression:     true,
		},
		"/api/extensions/clusters/zip": {
			AuthInterceptor: userAuth.HTTPInterceptor,
			Authorizer:      authzUser.Any(),
			ServerHandler:   clustersZip.Handler(clusterService, idService),
			Compression:     false,
		},
		"/db/backup": {
			AuthInterceptor: userAuth.HTTPInterceptor,
			Authorizer:      authzUser.Any(),
			ServerHandler:   c.database.BackupHandler(),
			Compression:     true,
		},
		"/db/export": {
			AuthInterceptor: userAuth.HTTPInterceptor,
			Authorizer:      authzUser.Any(),
			ServerHandler:   c.database.ExportHandler(),
			Compression:     true,
		},
	}

	c.addDebugRoutes(routeMap, userAuth)

	return
}

func (c *central) addDebugRoutes(routeMap map[string]routes.CustomRoute, userAuth *authnUser.AuthInterceptor) {
	rs := map[string]http.Handler{
		"/debug/pprof":         http.HandlerFunc(pprof.Index),
		"/debug/pprof/cmdline": http.HandlerFunc(pprof.Cmdline),
		"/debug/pprof/profile": http.HandlerFunc(pprof.Profile),
		"/debug/pprof/symbol":  http.HandlerFunc(pprof.Symbol),
		"/debug/pprof/trace":   http.HandlerFunc(pprof.Trace),
		"/debug/block":         pprof.Handler(`block`),
		"/debug/goroutine":     pprof.Handler(`goroutine`),
		"/debug/heap":          pprof.Handler(`heap`),
		"/debug/mutex":         pprof.Handler(`mutex`),
		"/debug/threadcreate":  pprof.Handler(`threadcreate`),
	}

	for r, h := range rs {
		routeMap[r] = routes.CustomRoute{
			AuthInterceptor: userAuth.HTTPInterceptor,
			Authorizer:      authzUser.Any(),
			ServerHandler:   h,
			Compression:     true,
		}
	}
}

func (c *central) processForever() {
	defer func() {
		if r := recover(); r != nil {
			log.Errorf("Caught panic in process loop; restarting. Stack: %s", string(debug.Stack()))
			c.processForever()
		}
	}()

	for {
		select {
		case sig := <-c.signalsC:
			log.Infof("Caught %s signal", sig)
			c.detector.Stop()
			c.database.Close()
			log.Infof("Central terminated")
			return
		}
	}
}
