package reconciler

import (
	pkgReconciler "github.com/operator-framework/helm-operator-plugins/pkg/reconciler"
	"github.com/stackrox/rox/image"
	platform "github.com/stackrox/rox/operator/apis/platform/v1alpha1"
	commonExtensions "github.com/stackrox/rox/operator/pkg/common/extensions"
	"github.com/stackrox/rox/operator/pkg/proxy"
	"github.com/stackrox/rox/operator/pkg/reconciler"
	"github.com/stackrox/rox/operator/pkg/securedcluster/extensions"
	scTranslation "github.com/stackrox/rox/operator/pkg/securedcluster/values/translation"
	"github.com/stackrox/rox/operator/pkg/utils"
	"github.com/stackrox/rox/operator/pkg/values/translation"
	"github.com/stackrox/rox/pkg/features"
	"github.com/stackrox/rox/pkg/version"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

// RegisterNewReconciler registers a new helm reconciler in the given k8s controller manager
func RegisterNewReconciler(mgr ctrl.Manager, selector string) error {
	proxyEnv := proxy.GetProxyEnvVars() // fix at startup time

	opts := []pkgReconciler.Option{
		pkgReconciler.WithExtraWatch(
			source.Kind(mgr.GetCache(), &platform.Central{}),
			reconciler.HandleSiblings(platform.SecuredClusterGVK, mgr),
			// Only appearance and disappearance of a Central resource can influence whether
			// a local scanner should be deployed by the SecuredCluster controller.
			utils.CreateAndDeleteOnlyPredicate{}),
		pkgReconciler.WithPreExtension(extensions.CheckClusterNameExtension(nil)),
		pkgReconciler.WithPreExtension(proxy.ReconcileProxySecretExtension(mgr.GetClient(), proxyEnv)),
		pkgReconciler.WithPreExtension(commonExtensions.CheckForbiddenNamespacesExtension(commonExtensions.IsSystemNamespace)),
		pkgReconciler.WithPreExtension(commonExtensions.ReconcileProductVersionStatusExtension(version.GetMainVersion())),
		pkgReconciler.WithPreExtension(extensions.ReconcileLocalScannerDBPasswordExtension(mgr.GetClient())),
	}

	if features.ScannerV4Support.Enabled() {
		opts = append(opts, pkgReconciler.WithPreExtension(extensions.ReconcileLocalScannerV4DBPasswordExtension(mgr.GetClient())))
	}

	opts, err := commonExtensions.AddSelectorOptionIfNeeded(selector, opts)
	if err != nil {
		return err
	}

	opts = commonExtensions.AddMapKubeAPIsExtensionIfMapFileExists(opts)

	return reconciler.SetupReconcilerWithManager(
		mgr, platform.SecuredClusterGVK,
		image.SecuredClusterServicesChartPrefix,
		translation.WithEnrichment(
			scTranslation.New(mgr.GetClient()),
			proxy.NewProxyEnvVarsInjector(proxyEnv, mgr.GetLogger())),
		opts...,
	)
}
