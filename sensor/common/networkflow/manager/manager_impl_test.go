package manager

import (
	"fmt"
	"testing"
	"time"

	"github.com/stackrox/rox/generated/internalapi/central"
	"github.com/stackrox/rox/generated/internalapi/sensor"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/env"
	"github.com/stackrox/rox/pkg/features"
	"github.com/stackrox/rox/pkg/net"
	"github.com/stackrox/rox/pkg/networkgraph"
	"github.com/stackrox/rox/pkg/protoconv"
	"github.com/stackrox/rox/pkg/timestamp"
	"github.com/stackrox/rox/sensor/common"
	"github.com/stackrox/rox/sensor/common/clusterentities"
	mocksDetector "github.com/stackrox/rox/sensor/common/detector/mocks"
	mocksManager "github.com/stackrox/rox/sensor/common/networkflow/manager/mocks"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

var (
	openNetworkEndpoint = &sensor.NetworkEndpoint{
		SocketFamily: sensor.SocketFamily_SOCKET_FAMILY_IPV4,
		Protocol:     storage.L4Protocol_L4_PROTOCOL_TCP,
		ContainerId:  "FakeContainerId",
		ListenAddress: &sensor.NetworkAddress{
			Port: 80,
		},
		Originator: &storage.NetworkProcessUniqueKey{
			ProcessName:         "socat",
			ProcessExecFilePath: "/usr/bin/socat",
			ProcessArgs:         "port: 80",
		},
	}
	openNetworkEndpoint81 = &sensor.NetworkEndpoint{
		SocketFamily: sensor.SocketFamily_SOCKET_FAMILY_IPV4,
		Protocol:     storage.L4Protocol_L4_PROTOCOL_TCP,
		ContainerId:  "FakeContainerId",
		ListenAddress: &sensor.NetworkAddress{
			Port: 81,
		},
		Originator: &storage.NetworkProcessUniqueKey{
			ProcessName:         "socat",
			ProcessExecFilePath: "/usr/bin/socat",
			ProcessArgs:         "port: 81",
		},
	}
	openNetworkEndpointNoOriginator = &sensor.NetworkEndpoint{
		SocketFamily: sensor.SocketFamily_SOCKET_FAMILY_IPV4,
		Protocol:     storage.L4Protocol_L4_PROTOCOL_TCP,
		ContainerId:  "FakeContainerId",
		ListenAddress: &sensor.NetworkAddress{
			Port: 80,
		},
	}
	closedNetworkEndpoint = &sensor.NetworkEndpoint{
		SocketFamily:   sensor.SocketFamily_SOCKET_FAMILY_IPV4,
		Protocol:       storage.L4Protocol_L4_PROTOCOL_TCP,
		ContainerId:    "FakeContainerId",
		CloseTimestamp: protoconv.ConvertTimeToTimestamp(time.Now()),
		ListenAddress: &sensor.NetworkAddress{
			Port: 80,
		},
		Originator: &storage.NetworkProcessUniqueKey{
			ProcessName:         "socat",
			ProcessExecFilePath: "/usr/bin/socat",
			ProcessArgs:         "port: 80",
		},
	}
)

const (
	waitTimeout = 20 * time.Millisecond
)

func TestNetworkFlowManager(t *testing.T) {
	suite.Run(t, new(NetworkFlowManagerTestSuite))
}

type NetworkFlowManagerTestSuite struct {
	suite.Suite
}

// region hostConnection.Process tests

func (s *NetworkFlowManagerTestSuite) TestAddNothing() {
	h := hostConnections{}
	h.endpoints = make(map[containerEndpoint]*connStatus)

	networkInfo := &sensor.NetworkConnectionInfo{}
	nowTimestamp := timestamp.Now()
	var sequenceID int64
	err := h.Process(networkInfo, nowTimestamp, sequenceID)
	s.NoError(err)
	s.Len(h.endpoints, 0)
}

func (s *NetworkFlowManagerTestSuite) TestAddOpen() {
	h := hostConnections{}
	h.endpoints = make(map[containerEndpoint]*connStatus)

	networkInfo := &sensor.NetworkConnectionInfo{
		UpdatedEndpoints: []*sensor.NetworkEndpoint{openNetworkEndpoint},
	}

	nowTimestamp := timestamp.Now()
	var sequenceID int64
	h.connectionsSequenceID = sequenceID
	err := h.Process(networkInfo, nowTimestamp, sequenceID)
	s.NoError(err)
	s.Len(h.endpoints, 1)
}

func (s *NetworkFlowManagerTestSuite) TestAddOpenAndClosed() {
	h := hostConnections{}
	h.endpoints = make(map[containerEndpoint]*connStatus)

	networkInfoOpen := &sensor.NetworkConnectionInfo{
		UpdatedEndpoints: []*sensor.NetworkEndpoint{openNetworkEndpoint},
	}

	networkInfoClosed := &sensor.NetworkConnectionInfo{
		UpdatedEndpoints: []*sensor.NetworkEndpoint{closedNetworkEndpoint},
	}

	nowTimestamp := timestamp.Now()
	var sequenceID int64
	h.connectionsSequenceID = sequenceID

	err := h.Process(networkInfoOpen, nowTimestamp, sequenceID)
	s.NoError(err)

	err = h.Process(networkInfoClosed, nowTimestamp, sequenceID)
	s.NoError(err)

	s.Len(h.endpoints, 1)
}

func (s *NetworkFlowManagerTestSuite) TestAddTwoDifferent() {
	h := hostConnections{}
	h.endpoints = make(map[containerEndpoint]*connStatus)

	networkInfoOpen := &sensor.NetworkConnectionInfo{
		UpdatedEndpoints: []*sensor.NetworkEndpoint{openNetworkEndpoint},
	}

	networkInfoOpen81 := &sensor.NetworkConnectionInfo{
		UpdatedEndpoints: []*sensor.NetworkEndpoint{openNetworkEndpoint81},
	}

	nowTimestamp := timestamp.Now()
	var sequenceID int64
	h.connectionsSequenceID = sequenceID

	err := h.Process(networkInfoOpen, nowTimestamp, sequenceID)
	s.NoError(err)

	err = h.Process(networkInfoOpen81, nowTimestamp, sequenceID)
	s.NoError(err)

	s.Len(h.endpoints, 2)
}

func (s *NetworkFlowManagerTestSuite) TestAddTwoDifferentSameBatch() {
	h := hostConnections{}
	h.endpoints = make(map[containerEndpoint]*connStatus)

	networkInfoOpen := &sensor.NetworkConnectionInfo{
		UpdatedEndpoints: []*sensor.NetworkEndpoint{openNetworkEndpoint, openNetworkEndpoint81},
	}

	nowTimestamp := timestamp.Now()
	var sequenceID int64
	h.connectionsSequenceID = sequenceID

	err := h.Process(networkInfoOpen, nowTimestamp, sequenceID)
	s.NoError(err)

	s.Len(h.endpoints, 2)
}

func (s *NetworkFlowManagerTestSuite) TestAddNoOriginator() {
	h := hostConnections{}
	h.endpoints = make(map[containerEndpoint]*connStatus)

	networkInfoOpen := &sensor.NetworkConnectionInfo{
		UpdatedEndpoints: []*sensor.NetworkEndpoint{openNetworkEndpointNoOriginator},
	}

	nowTimestamp := timestamp.Now()
	var sequenceID int64
	h.connectionsSequenceID = sequenceID

	err := h.Process(networkInfoOpen, nowTimestamp, sequenceID)
	s.NoError(err)

	s.Len(h.endpoints, 1)
}

// endregion

// region networkFlowManager tests

func (s *NetworkFlowManagerTestSuite) TestEnrichConnection() {
	mockCtrl := gomock.NewController(s.T())
	m, mockEntityStore, mockExternalSrc, _ := createManager(mockCtrl)
	srcID := "src-id"
	dstID := "dst-id"
	cases := map[string]struct {
		connPair                    *connectionPair
		enrichedConnections         map[networkConnIndicator]timestamp.MicroTS
		expectEntityLookupContainer expectFn
		expectEntityLookupEndpoint  expectFn
		expectExternalLookup        expectFn
		expectedIndicator           *networkConnIndicator
		expectedConnection          *connection
		expectedStatus              *connStatus
	}{
		"Rotten connection should return rotten status": {
			connPair: createConnectionPair().incoming().external().firstSeen(timestamp.Now().Add(-maxContainerResolutionWaitPeriod * 2)),
			expectEntityLookupContainer: func() {
				mockEntityStore.EXPECT().LookupByContainerID(gomock.Any()).Times(1).DoAndReturn(func(_ any) (clusterentities.ContainerMetadata, bool) {
					return clusterentities.ContainerMetadata{}, false
				})
			},
			expectedStatus: &connStatus{
				rotten: true,
			},
		},
		"Incoming external connection with unsuccessful lookup should return internet entity": {
			connPair:            createConnectionPair().incoming().external(),
			enrichedConnections: make(map[networkConnIndicator]timestamp.MicroTS),
			expectEntityLookupContainer: func() {
				mockEntityStore.EXPECT().LookupByContainerID(gomock.Any()).Times(1).DoAndReturn(func(_ any) (clusterentities.ContainerMetadata, bool) {
					return clusterentities.ContainerMetadata{
						DeploymentID: dstID,
					}, true
				})
			},
			expectExternalLookup: func() {
				mockExternalSrc.EXPECT().LookupByNetwork(gomock.Any()).Times(1).DoAndReturn(func(_ any) *storage.NetworkEntityInfo {
					return nil
				})
			},
			expectedStatus: &connStatus{
				used: true,
			},
			expectedIndicator: &networkConnIndicator{
				dstPort:   80,
				protocol:  net.TCP.ToProtobuf(),
				srcEntity: networkgraph.InternetEntity(),
				dstEntity: networkgraph.EntityForDeployment(dstID),
			},
		},
		"Outgoing external connection with successful external lookup should return the correct id": {
			connPair:            createConnectionPair().external(),
			enrichedConnections: make(map[networkConnIndicator]timestamp.MicroTS),
			expectEntityLookupContainer: func() {
				mockEntityStore.EXPECT().LookupByContainerID(gomock.Any()).Times(1).DoAndReturn(func(_ any) (clusterentities.ContainerMetadata, bool) {
					return clusterentities.ContainerMetadata{
						DeploymentID: srcID,
					}, true
				})
			},
			expectExternalLookup: func() {
				mockExternalSrc.EXPECT().LookupByNetwork(gomock.Any()).Times(1).DoAndReturn(func(_ any) *storage.NetworkEntityInfo {
					return &storage.NetworkEntityInfo{
						Id: dstID,
					}
				})

			},
			expectedStatus: &connStatus{
				used: true,
			},
			expectedIndicator: &networkConnIndicator{
				dstPort:  80,
				protocol: net.TCP.ToProtobuf(),
				dstEntity: networkgraph.EntityFromProto(&storage.NetworkEntityInfo{
					Id: dstID,
				}),
				srcEntity: networkgraph.EntityForDeployment(srcID),
			},
		},
		"Incoming connection with successful lookup should not return a networkConnIndicator": {
			connPair:            createConnectionPair().incoming(),
			enrichedConnections: make(map[networkConnIndicator]timestamp.MicroTS),
			expectEntityLookupContainer: func() {
				mockEntityStore.EXPECT().LookupByContainerID(gomock.Any()).Times(1).DoAndReturn(func(_ any) (clusterentities.ContainerMetadata, bool) {
					return clusterentities.ContainerMetadata{
						DeploymentID: srcID,
					}, true
				})
			},
			expectEntityLookupEndpoint: func() {
				mockEntityStore.EXPECT().LookupByEndpoint(gomock.Any()).Times(1).DoAndReturn(func(_ any) []clusterentities.LookupResult {
					return []clusterentities.LookupResult{
						{
							Entity: networkgraph.Entity{
								ID: dstID,
							},
						},
					}
				})
			},
			expectedStatus: &connStatus{
				used: true,
			},
		},
		"Incoming fresh connection with valid address should not return anything": {
			connPair:            createConnectionPair().incoming(),
			enrichedConnections: make(map[networkConnIndicator]timestamp.MicroTS),
			expectEntityLookupContainer: func() {
				mockEntityStore.EXPECT().LookupByContainerID(gomock.Any()).Times(1).DoAndReturn(func(_ any) (clusterentities.ContainerMetadata, bool) {
					return clusterentities.ContainerMetadata{
						DeploymentID: dstID,
					}, true
				})
			},
			expectEntityLookupEndpoint: func() {
				mockEntityStore.EXPECT().LookupByEndpoint(gomock.Any()).Times(1).DoAndReturn(func(_ any) []clusterentities.LookupResult {
					return nil
				})
			},
			expectedStatus: &connStatus{},
		},
		"Incoming fresh connection with invalid address should not return anything": {
			connPair:            createConnectionPair().incoming().invalidAddress(),
			enrichedConnections: make(map[networkConnIndicator]timestamp.MicroTS),
			expectEntityLookupContainer: func() {
				mockEntityStore.EXPECT().LookupByContainerID(gomock.Any()).Times(1).DoAndReturn(func(_ any) (clusterentities.ContainerMetadata, bool) {
					return clusterentities.ContainerMetadata{
						DeploymentID: dstID,
					}, true
				})
			},
			expectEntityLookupEndpoint: func() {
				mockEntityStore.EXPECT().LookupByEndpoint(gomock.Any()).Times(1).DoAndReturn(func(_ any) []clusterentities.LookupResult {
					return nil
				})
			},
			expectExternalLookup: func() {
				mockExternalSrc.EXPECT().LookupByNetwork(gomock.Any()).Times(1).DoAndReturn(func(_ any) *storage.NetworkEntityInfo {
					return nil
				})
			},
			expectedStatus: &connStatus{},
		},
		"Outgoing connection with successful internal lookup should return the correct id": {
			connPair:            createConnectionPair(),
			enrichedConnections: make(map[networkConnIndicator]timestamp.MicroTS),
			expectEntityLookupContainer: func() {
				mockEntityStore.EXPECT().LookupByContainerID(gomock.Any()).Times(1).DoAndReturn(func(_ any) (clusterentities.ContainerMetadata, bool) {
					return clusterentities.ContainerMetadata{
						DeploymentID: srcID,
					}, true
				})
			},
			expectEntityLookupEndpoint: func() {
				mockEntityStore.EXPECT().LookupByEndpoint(gomock.Any()).Times(1).DoAndReturn(func(_ any) []clusterentities.LookupResult {
					return []clusterentities.LookupResult{
						{
							Entity: networkgraph.Entity{
								ID: dstID,
							},
							ContainerPorts: []uint16{
								80,
							},
						},
					}
				})
			},
			expectedStatus: &connStatus{
				used: true,
			},
			expectedIndicator: &networkConnIndicator{
				dstPort:  80,
				protocol: net.TCP.ToProtobuf(),
				dstEntity: networkgraph.EntityFromProto(&storage.NetworkEntityInfo{
					Id: dstID,
				}),
				srcEntity: networkgraph.EntityForDeployment(srcID),
			},
		},
	}
	for name, tCase := range cases {
		s.Run(name, func() {
			tCase.expectEntityLookupContainer.runIfSet()
			tCase.expectEntityLookupEndpoint.runIfSet()
			tCase.expectExternalLookup.runIfSet()
			m.enrichConnection(tCase.connPair.conn, tCase.connPair.status, tCase.enrichedConnections)
			s.Assert().Equal(tCase.expectedStatus.used, tCase.connPair.status.used)
			s.Assert().Equal(tCase.expectedStatus.rotten, tCase.connPair.status.rotten)
			if tCase.expectedIndicator != nil {
				_, ok := tCase.enrichedConnections[*tCase.expectedIndicator]
				s.Assert().True(ok)
			} else {
				s.Assert().Len(tCase.enrichedConnections, 0)
			}
		})
	}
}

func (s *NetworkFlowManagerTestSuite) TestEnrichContainerEndpoint() {
	mockCtrl := gomock.NewController(s.T())
	m, mockEntityStore, _, _ := createManager(mockCtrl)
	id := "id"
	cases := map[string]struct {
		endpointPair                *endpointPair
		enrichedConnections         map[containerEndpointIndicator]timestamp.MicroTS
		expectEntityLookupContainer expectFn
		expectedStatus              *connStatus
		expectedEndpoint            *containerEndpointIndicator
	}{
		"Rotten connection should return rotten status": {
			endpointPair:                createEndpointPair(timestamp.Now().Add(-maxContainerResolutionWaitPeriod * 2)),
			expectEntityLookupContainer: expectEntityLookupContainerHelper(mockEntityStore, 1, clusterentities.ContainerMetadata{}, false),
			expectedStatus: &connStatus{
				rotten: true,
				used:   true,
			},
		},
		"Container endpoint should return an containerEndpointIndicator with the correct id": {
			endpointPair:        createEndpointPair(timestamp.Now()),
			enrichedConnections: make(map[containerEndpointIndicator]timestamp.MicroTS),
			expectEntityLookupContainer: expectEntityLookupContainerHelper(mockEntityStore, 1, clusterentities.ContainerMetadata{
				DeploymentID: id,
			}, true),
			expectedStatus: &connStatus{used: true},
			expectedEndpoint: &containerEndpointIndicator{
				entity:   networkgraph.EntityForDeployment(id),
				port:     80,
				protocol: net.TCP.ToProtobuf(),
			},
		},
	}
	for name, tCase := range cases {
		s.Run(name, func() {
			tCase.expectEntityLookupContainer.runIfSet()
			m.enrichContainerEndpoint(tCase.endpointPair.endpoint, tCase.endpointPair.status, tCase.enrichedConnections)
			s.Assert().Equal(tCase.expectedStatus.rotten, tCase.endpointPair.status.rotten)
			s.Assert().Equal(tCase.expectedStatus.used, tCase.endpointPair.status.used)
			if tCase.expectedEndpoint != nil {
				_, ok := tCase.enrichedConnections[*tCase.expectedEndpoint]
				s.Assert().True(ok)
			}
		})
	}
}

func (s *NetworkFlowManagerTestSuite) TestEnrichProcessListening() {
	mockCtrl := gomock.NewController(s.T())
	m, mockEntityStore, _, _ := createManager(mockCtrl)
	deploymentID := "deployment-id"
	podID := "pod-id"
	cases := map[string]struct {
		containerPair               *containerPair
		enrichedConnections         map[processListeningIndicator]timestamp.MicroTS
		expectEntityLookupContainer expectFn
		expectedStatus              *connStatus
		expectedListeningIndicator  *processListeningIndicator
	}{
		"Rotten connection should return rotten status": {
			containerPair:               createContainerPair(timestamp.Now().Add(-maxContainerResolutionWaitPeriod * 2)),
			expectEntityLookupContainer: expectEntityLookupContainerHelper(mockEntityStore, 1, clusterentities.ContainerMetadata{}, false),
			expectedStatus: &connStatus{
				rotten:      true,
				usedProcess: true,
			},
		},
		"Container endpoint should return a processListeningIndicator with the correct id": {
			containerPair:       createContainerPair(timestamp.Now()),
			enrichedConnections: make(map[processListeningIndicator]timestamp.MicroTS),
			expectEntityLookupContainer: expectEntityLookupContainerHelper(mockEntityStore, 1, clusterentities.ContainerMetadata{
				DeploymentID:  deploymentID,
				ContainerName: "container-name",
				PodID:         podID,
			}, true),
			expectedStatus: &connStatus{
				usedProcess: true,
			},
			expectedListeningIndicator: &processListeningIndicator{
				key: processUniqueKey{
					podID:         podID,
					containerName: "container-name",
					deploymentID:  deploymentID,
					process:       defaultProcessKey(),
				},
				port:     80,
				protocol: net.TCP.ToProtobuf(),
			},
		},
	}
	for name, tCase := range cases {
		s.Run(name, func() {
			tCase.expectEntityLookupContainer.runIfSet()
			m.enrichProcessListening(tCase.containerPair.endpoint, tCase.containerPair.status, tCase.enrichedConnections)
			s.Assert().Equal(tCase.expectedStatus.rotten, tCase.containerPair.status.rotten)
			s.Assert().Equal(tCase.expectedStatus.usedProcess, tCase.containerPair.status.usedProcess)
			if tCase.expectedListeningIndicator != nil {
				_, ok := tCase.enrichedConnections[*tCase.expectedListeningIndicator]
				s.Assert().True(ok)
			}
		})
	}
}

func (s *NetworkFlowManagerTestSuite) TestManagerOfflineMode() {
	s.T().Setenv(env.ProcessesListeningOnPort.EnvVar(), "false")
	srcID := "src-id"
	dstID := "dst-id"
	mockCtrl := gomock.NewController(s.T())
	hostname := "hostname"
	containerID := "container-id"
	m, mockEntity, _, mockDetector := createManager(mockCtrl)
	states := []struct {
		testName                    string
		notify                      common.SensorComponentEvent
		connections                 []*HostnameAndConnections
		expectEntityLookupContainer expectFn
		expectEntityLookupEndpoint  expectFn
		expectDetector              expectFn
		expectedSensorMessage       *central.MsgFromSensor
	}{
		{
			testName:    "In offline mode we should not send any messages upon receiving a connection",
			notify:      common.SensorComponentEventOfflineMode,
			connections: []*HostnameAndConnections{createHostnameConnections(hostname).withConnectionPair(createConnectionPair())},
		},
		{
			testName: "In online mode we should enrich and send the previously received connection",
			notify:   common.SensorComponentEventCentralReachable,
			expectEntityLookupContainer: expectEntityLookupContainerHelper(mockEntity, 1, clusterentities.ContainerMetadata{
				DeploymentID: srcID,
			}, true),
			expectEntityLookupEndpoint: expectEntityLookupEndpointHelper(mockEntity, 1, []clusterentities.LookupResult{
				{
					Entity:         networkgraph.Entity{ID: dstID},
					ContainerPorts: []uint16{80},
				},
			}),
			expectDetector:        expectDetectorHelper(mockDetector, 1),
			expectedSensorMessage: createExpectedSensorMessageWithConnections(&expectedEntitiesPair{srcID: srcID, dstID: dstID}),
		},
		{
			testName: "In offline mode we should not send any messages upon receiving multiple connections",
			notify:   common.SensorComponentEventOfflineMode,
			connections: []*HostnameAndConnections{
				createHostnameConnections(hostname).withConnectionPair(createConnectionPair().containerID(fmt.Sprintf("%s-1", containerID))),
				createHostnameConnections(hostname).withConnectionPair(createConnectionPair().containerID(fmt.Sprintf("%s-2", containerID))),
			},
		},
		{
			testName: "In online mode we should enrich and send the previously received connections",
			notify:   common.SensorComponentEventCentralReachable,
			expectEntityLookupContainer: func() {
				gomock.InOrder(
					mockEntity.EXPECT().LookupByContainerID(gomock.Any()).Times(1).DoAndReturn(func(_ any) (clusterentities.ContainerMetadata, bool) {
						return clusterentities.ContainerMetadata{DeploymentID: fmt.Sprintf("%s-1", srcID)}, true
					}),
					mockEntity.EXPECT().LookupByContainerID(gomock.Any()).Times(1).DoAndReturn(func(_ any) (clusterentities.ContainerMetadata, bool) {
						return clusterentities.ContainerMetadata{DeploymentID: fmt.Sprintf("%s-2", srcID)}, true
					}),
				)
			},
			expectEntityLookupEndpoint: func() {
				gomock.InOrder(
					mockEntity.EXPECT().LookupByEndpoint(gomock.Any()).Times(1).DoAndReturn(func(_ any) []clusterentities.LookupResult {
						return []clusterentities.LookupResult{
							{
								Entity:         networkgraph.Entity{ID: fmt.Sprintf("%s-1", dstID)},
								ContainerPorts: []uint16{80},
							},
						}
					}),
					mockEntity.EXPECT().LookupByEndpoint(gomock.Any()).Times(1).DoAndReturn(func(_ any) []clusterentities.LookupResult {
						return []clusterentities.LookupResult{
							{
								Entity:         networkgraph.Entity{ID: fmt.Sprintf("%s-2", dstID)},
								ContainerPorts: []uint16{80},
							},
						}
					}),
				)
			},
			expectDetector: expectDetectorHelper(mockDetector, 2),
			expectedSensorMessage: createExpectedSensorMessageWithConnections(
				&expectedEntitiesPair{srcID: fmt.Sprintf("%s-1", srcID), dstID: fmt.Sprintf("%s-1", dstID)},
				&expectedEntitiesPair{srcID: fmt.Sprintf("%s-2", srcID), dstID: fmt.Sprintf("%s-2", dstID)},
			),
		},
		{
			testName: "In offline mode we should not send any messages upon receiving multiple endpoints",
			notify:   common.SensorComponentEventOfflineMode,
			connections: []*HostnameAndConnections{
				createHostnameConnections(hostname).withEndpointPair(createEndpointPair(timestamp.Now()).containerID(fmt.Sprintf("%s-1", containerID))),
				createHostnameConnections(hostname).withEndpointPair(createEndpointPair(timestamp.Now()).containerID(fmt.Sprintf("%s-2", containerID))),
			},
		},
		{
			testName: "In online mode we should enrich and send the previously received endpoints",
			notify:   common.SensorComponentEventCentralReachable,
			expectEntityLookupContainer: func() {
				gomock.InOrder(
					mockEntity.EXPECT().LookupByContainerID(gomock.Any()).Times(1).DoAndReturn(func(_ any) (clusterentities.ContainerMetadata, bool) {
						return clusterentities.ContainerMetadata{DeploymentID: fmt.Sprintf("%s-1", srcID)}, true
					}),
					mockEntity.EXPECT().LookupByContainerID(gomock.Any()).Times(1).DoAndReturn(func(_ any) (clusterentities.ContainerMetadata, bool) {
						return clusterentities.ContainerMetadata{DeploymentID: fmt.Sprintf("%s-2", srcID)}, true
					}),
				)
			},
			expectedSensorMessage: createExpectedSensorMessageWithEndpoints(
				fmt.Sprintf("%s-1", srcID),
				fmt.Sprintf("%s-2", srcID),
			),
		},
	}
	fakeTicker := make(chan time.Time)
	defer close(fakeTicker)
	go m.enrichConnections(fakeTicker)
	for _, state := range states {
		for _, cnn := range state.connections {
			addHostConnection(m, cnn)
		}
		s.Run(state.testName, func() {
			state.expectEntityLookupContainer.runIfSet()
			state.expectEntityLookupEndpoint.runIfSet()
			state.expectDetector.runIfSet()
			// We do not test ticking here, but without this line, the test would deadlock.
			mockEntity.EXPECT().RecordTick().AnyTimes()
			m.Notify(state.notify)
			fakeTicker <- time.Now()
			if state.expectedSensorMessage != nil {
				select {
				case <-time.After(10 * time.Second):
					s.Fail("timeout waiting for sensor message")
				case msg, ok := <-m.sensorUpdates:
					s.Require().True(ok, "the sensorUpdates channel should not be closed")
					s.Assert().NotNil(msg)
					msgFromSensor, ok := msg.Msg.(*central.MsgFromSensor_NetworkFlowUpdate)
					s.Require().True(ok, "the message received is not a NetworkFlowUpdate message")
					expectedMsg, ok := state.expectedSensorMessage.Msg.(*central.MsgFromSensor_NetworkFlowUpdate)
					s.Require().True(ok, "the message expected is not a NetworkFlowUpdate message")
					s.Assert().Len(msgFromSensor.NetworkFlowUpdate.GetUpdated(), len(expectedMsg.NetworkFlowUpdate.GetUpdated()))
					s.assertSensorMessageConnectionIDs(expectedMsg.NetworkFlowUpdate.GetUpdated(), msgFromSensor.NetworkFlowUpdate.GetUpdated())
					s.Assert().Len(msgFromSensor.NetworkFlowUpdate.GetUpdatedEndpoints(), len(expectedMsg.NetworkFlowUpdate.GetUpdatedEndpoints()))
					s.assertSensorMessageEndpointIDs(expectedMsg.NetworkFlowUpdate.GetUpdatedEndpoints(), msgFromSensor.NetworkFlowUpdate.GetUpdatedEndpoints())
				}
			} else {
				select {
				case _, ok := <-m.sensorUpdates:
					s.Require().True(ok, "the sensorUpdates channel should not be closed")
					s.Fail("should not received message in sensorUpdates channel")
				case <-time.After(time.Second):
					break
				}
			}
		})
	}
	m.Stop(nil)
}

func (s *NetworkFlowManagerTestSuite) TestExpireMessage() {
	s.T().Setenv(env.ProcessesListeningOnPort.EnvVar(), "false")
	mockCtrl := gomock.NewController(s.T())
	hostname := "hostname"
	containerID := "container-id"
	m, mockEntity, _, mockDetector := createManager(mockCtrl)
	fakeTicker := make(chan time.Time)
	defer close(fakeTicker)
	go m.enrichConnections(fakeTicker)
	mockEntity.EXPECT().LookupByContainerID(gomock.Any()).Times(1).DoAndReturn(func(_ any) (clusterentities.ContainerMetadata, bool) {
		return clusterentities.ContainerMetadata{
			DeploymentID: containerID,
		}, true
	})
	mockEntity.EXPECT().LookupByEndpoint(gomock.Any()).Times(1).DoAndReturn(func(_ any) []clusterentities.LookupResult {
		return []clusterentities.LookupResult{
			{
				Entity:         networkgraph.Entity{ID: containerID},
				ContainerPorts: []uint16{80},
			},
		}
	})
	mockDetector.EXPECT().ProcessNetworkFlow(gomock.Any(), gomock.Any()).Times(1)
	mockEntity.EXPECT().RecordTick().AnyTimes()
	addHostConnection(m, createHostnameConnections(hostname).withConnectionPair(createConnectionPair()))
	m.Notify(common.SensorComponentEventCentralReachable)
	fakeTicker <- time.Now()
	select {
	case <-time.After(10 * time.Second):
		s.Fail("timeout waiting for sensor message")
	case msg, ok := <-m.sensorUpdates:
		s.Require().True(ok, "the sensorUpdates channel should not be closed")
		m.Notify(common.SensorComponentEventOfflineMode)
		m.Notify(common.SensorComponentEventCentralReachable)
		s.Assert().True(msg.IsExpired(), "the message should be expired")
	}
	m.Stop(nil)
}

func TestSendNetworkFlows(t *testing.T) {
	t.Setenv(features.SensorCapturesIntermediateEvents.EnvVar(), "true")
	suite.Run(t, new(sendNetflowsSuite))
}

type sendNetflowsSuite struct {
	suite.Suite
	mockCtrl     *gomock.Controller
	mockEntity   *mocksManager.MockEntityStore
	m            *networkFlowManager
	mockDetector *mocksDetector.MockDetector
	fakeTicker   chan time.Time
}

const (
	srcID = "src-id"
	dstID = "dst-id"
)

func (b *sendNetflowsSuite) SetupTest() {
	b.mockCtrl = gomock.NewController(b.T())
	b.m, b.mockEntity, _, b.mockDetector = createManager(b.mockCtrl)

	b.fakeTicker = make(chan time.Time)
	go b.m.enrichConnections(b.fakeTicker)
}

func (b *sendNetflowsSuite) TeardownTest() {
	b.m.done.Signal()
}

func (b *sendNetflowsSuite) updateConn(pair *connectionPair) {
	addHostConnection(b.m, createHostnameConnections("hostname").withConnectionPair(pair))
}

func (b *sendNetflowsSuite) expectLookups(n int) {
	b.mockEntity.EXPECT().RecordTick().AnyTimes()
	expectEntityLookupContainerHelper(b.mockEntity, n, clusterentities.ContainerMetadata{
		DeploymentID: srcID,
	}, true)()
	expectEntityLookupEndpointHelper(b.mockEntity, n, []clusterentities.LookupResult{
		{
			Entity:         networkgraph.Entity{ID: dstID},
			ContainerPorts: []uint16{80},
		},
	})()
}

func (b *sendNetflowsSuite) expectDetections(n int) {
	expectDetectorHelper(b.mockDetector, n)()
}

func (b *sendNetflowsSuite) TestUpdateConnectionGeneratesNetflow() {
	b.expectLookups(1)
	b.expectDetections(1)

	b.updateConn(createConnectionPair())
	b.thenTickerTicks()
	b.assertOneUpdatedConnection()
}

func (b *sendNetflowsSuite) TestUnchangedConnection() {
	b.expectLookups(2)
	b.expectDetections(1)

	pair := createConnectionPair()
	pair.status.lastSeen = timestamp.InfiniteFuture
	b.updateConn(pair)
	b.thenTickerTicks()
	b.assertOneUpdatedConnection()

	// There should be no second update, the connection did not change
	b.thenTickerTicks()
	mustNotRead(b.T(), b.m.sensorUpdates)
}

func (b *sendNetflowsSuite) TestSendTwoUpdatesOnConnectionChanged() {
	b.expectLookups(2)
	b.expectDetections(2)

	pair := createConnectionPair()
	oneHourAgo := timestamp.NowMinus(time.Hour)
	pair.status.lastSeen = timestamp.FromProtobuf(oneHourAgo)
	b.updateConn(pair)
	b.thenTickerTicks()
	b.assertOneUpdatedConnection()

	pair.status.lastSeen = timestamp.Now()
	b.updateConn(pair)
	b.thenTickerTicks()
	b.assertOneUpdatedConnection()
}

func (b *sendNetflowsSuite) TestUpdatesGetBufferedWhenUnread() {
	b.expectLookups(4)
	b.expectDetections(4)

	// four times without reading
	for i := 0; i < 4; i++ {
		ts := timestamp.NowMinus(time.Duration(4-i) * time.Hour)
		pair := createConnectionPair()
		pair.status.lastSeen = timestamp.FromProtobuf(ts)
		b.updateConn(pair)
		b.thenTickerTicks()
		time.Sleep(100 * time.Millisecond) // Immediately ticking without waiting causes unexpected behavior
	}

	// should be able to read four buffered updates in sequence
	for i := 0; i < 4; i++ {
		b.assertOneUpdatedConnection()
	}
}

func (b *sendNetflowsSuite) TestCallsDetectionEvenOnFullBuffer() {
	b.expectLookups(6)
	b.expectDetections(6)

	for i := 0; i < 6; i++ {
		ts := timestamp.NowMinus(time.Duration(6-i) * time.Hour)
		pair := createConnectionPair()
		pair.status.lastSeen = timestamp.FromProtobuf(ts)
		b.updateConn(pair)
		b.thenTickerTicks()
		time.Sleep(100 * time.Millisecond)
	}

	// Will only store 5 network flow updates, as it's the maximum buffer size in the test
	for i := 0; i < 5; i++ {
		b.assertOneUpdatedConnection()
	}

	mustNotRead(b.T(), b.m.sensorUpdates)
}

func (b *sendNetflowsSuite) thenTickerTicks() {
	mustSendWithoutBlock(b.T(), b.fakeTicker, time.Now())
}

func (b *sendNetflowsSuite) assertOneUpdatedConnection() {
	msg := mustReadTimeout(b.T(), b.m.sensorUpdates)
	netflowUpdate, ok := msg.Msg.(*central.MsgFromSensor_NetworkFlowUpdate)
	b.Require().True(ok, "message is NetworkFlowUpdate")
	b.Assert().Len(netflowUpdate.NetworkFlowUpdate.GetUpdated(), 1, "one updated connection")
}

func mustNotRead[T any](t *testing.T, ch chan T) {
	select {
	case <-ch:
		t.Fatal("should not receive in channel")
	case <-time.After(waitTimeout):
	}
}

func mustReadTimeout[T any](t *testing.T, ch chan T) T {
	var result T
	select {
	case v, more := <-ch:
		if !more {
			require.True(t, more, "channel should never close")
		}
		result = v
	case <-time.After(waitTimeout):
		t.Fatal("blocked on reading from channel")
	}
	return result
}

func mustSendWithoutBlock[T any](t *testing.T, ch chan T, v T) {
	select {
	case ch <- v:
		return
	case <-time.After(waitTimeout):
		t.Fatal("blocked on sending to channel")
	}
}

// endregion
