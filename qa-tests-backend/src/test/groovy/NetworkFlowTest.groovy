import com.google.protobuf.Timestamp
import groups.BAT
import groups.NetworkFlowVisualization
import objects.Deployment
import objects.NetworkPolicy
import objects.NetworkPolicyTypes
import org.junit.experimental.categories.Category
import services.NetworkGraphService
import spock.lang.Unroll
import v1.NetworkEnums
import v1.NetworkGraphOuterClass
import v1.NetworkGraphOuterClass.NetworkGraph

class NetworkFlowTest extends BaseSpecification {
    static final private NETWORK_FLOW_UPDATE_CADENCE = 30000 // Network flow data is updated every 30 seconds

    // Deployment names
    static final private String UDPCONNECTIONTARGET = "udp-connection-target"
    static final private String TCPCONNECTIONTARGET = "tcp-connection-target"
    static final private String NGINXCONNECTIONTARGET = "nginx-connection-target"
    static final private String UDPCONNECTIONSOURCE = "udp-connection-source"
    static final private String TCPCONNECTIONSOURCE = "tcp-connection-source"
    //static final private String ICMPCONNECTIONSOURCE = "icmp-connection-source"
    static final private String NOCONNECTIONSOURCE = "no-connection-source"
    static final private String SHORTCONSISTENTSOURCE = "short-consistent-source"
    static final private String SINGLECONNECTIONSOURCE = "single-connection-source"
    static final private String EXTERNALCONNECTIONSOURCE = "external-connection-source"
    static final private String MULTIPLEPORTSCONNECTION = "two-ports-connect-source"

    static final private List<Deployment> DEPLOYMENTS = [
            //Target deployments
            new Deployment()
                    .setName(UDPCONNECTIONTARGET)
                    .setImage("ubuntu")
                    .addPort(8080, "UDP")
                    .addLabel("app", UDPCONNECTIONTARGET)
                    .setExposeAsService(true)
                    .setCommand(["/bin/sh", "-c",])
                    .setArgs(["apt-get update && " +
                                      "apt-get install socat -y && " +
                                      "socat -d -d -v UDP-RECV:8080 STDOUT" as String,]),
            new Deployment()
                    .setName(TCPCONNECTIONTARGET)
                    .setImage("ubuntu")
                    .addPort(80)
                    .addPort(8080)
                    .addLabel("app", TCPCONNECTIONTARGET)
                    .setExposeAsService(true)
                    .setCommand(["/bin/sh", "-c",])
                    .setArgs(["apt-get update && " +
                                      "apt-get install socat -y && " +
                                      "(socat -d -d -v TCP-LISTEN:80,fork STDOUT & " +
                                      "socat -d -d -v TCP-LISTEN:8080,fork STDOUT)" as String,]),
            new Deployment()
                    .setName(NGINXCONNECTIONTARGET)
                    .setImage("nginx")
                    .addPort(80)
                    .addLabel("app", NGINXCONNECTIONTARGET)
                    .setExposeAsService(true),

            //Source deployments
            new Deployment()
                    .setName(NOCONNECTIONSOURCE)
                    .setImage("nginx")
                    .addLabel("app", NOCONNECTIONSOURCE),
            new Deployment()
                    .setName(SHORTCONSISTENTSOURCE)
                    .setImage("nginx:1.15.4-alpine")
                    .addLabel("app", SHORTCONSISTENTSOURCE)
                    .setCommand(["/bin/sh", "-c",])
                    .setArgs(["while sleep ${NETWORK_FLOW_UPDATE_CADENCE / 1000}; " +
                                      "do wget -S http://${NGINXCONNECTIONTARGET}; " +
                                      "done" as String,]),
            new Deployment()
                    .setName(SINGLECONNECTIONSOURCE)
                    .setImage("nginx:1.15.4-alpine")
                    .addLabel("app", SINGLECONNECTIONSOURCE)
                    .setCommand(["/bin/sh", "-c",])
                    .setArgs(["wget -S -T 2 http://${NGINXCONNECTIONTARGET} && " +
                                      "while sleep 30; do echo hello; done" as String,]),
            new Deployment()
                    .setName(UDPCONNECTIONSOURCE)
                    .setImage("ubuntu")
                    .addLabel("app", UDPCONNECTIONSOURCE)
                    .setCommand(["/bin/sh", "-c",])
                    .setArgs(["apt-get update && " +
                                      "apt-get install socat -y && " +
                                      "while sleep 5; " +
                                      "do socat -s STDIN UDP:${UDPCONNECTIONTARGET}:8080; " +
                                      "done" as String,]),
            new Deployment()
                    .setName(TCPCONNECTIONSOURCE)
                    .setImage("ubuntu")
                    .addLabel("app", TCPCONNECTIONSOURCE)
                    .setCommand(["/bin/sh", "-c",])
                    .setArgs(["apt-get update && " +
                                      "apt-get install socat -y && " +
                                      "while sleep 5; " +
                                      "do socat -s STDIN TCP:${TCPCONNECTIONTARGET}:80; " +
                                      "done" as String,]),
            new Deployment()
                    .setName(EXTERNALCONNECTIONSOURCE)
                    .setImage("ubuntu")
                    .addLabel("app", EXTERNALCONNECTIONSOURCE)
                    .setCommand(["/bin/sh", "-c"])
                    .setArgs(["/usr/bin/apt-get update && " +
                                      "/usr/bin/apt-get install curl --assume-yes && " +
                                      "while sleep ${NETWORK_FLOW_UPDATE_CADENCE / 1000}; " +
                                      "do curl http://www.google.com; " +
                                      "done", ]),
            new Deployment()
                    .setName(MULTIPLEPORTSCONNECTION)
                    .setImage("ubuntu")
                    .addLabel("app", MULTIPLEPORTSCONNECTION)
                    .setCommand(["/bin/sh", "-c",])
                    .setArgs(["apt-get update && " +
                                      "apt-get install socat -y && " +
                                      "while sleep 5; " +
                                      "do socat -s STDIN TCP:${TCPCONNECTIONTARGET}:80; " +
                                      "socat -s STDIN TCP:${TCPCONNECTIONTARGET}:8080; " +
                                      "done" as String,]),
    ]

    def setupSpec() {
        orchestrator.batchCreateDeployments(DEPLOYMENTS)
        //
        // Commenting out ICMP test setup for now
        // See ROX-635
        //
        /*
        def nginxIp = DEPLOYMENTS.find { it.name == NGINXCONNECTIONTARGET }?.pods?.get(0)?.podIP
        Deployment icmp = new Deployment()
                .setName(ICMPCONNECTIONSOURCE)
                .setImage("ubuntu")
                .addLabel("app", ICMPCONNECTIONSOURCE)
                .setCommand(["/bin/sh", "-c",])
                .setArgs(["apt-get update && " +
                                  "apt-get install iputils-ping -y && " +
                                  "ping ${nginxIp}" as String,])
        orchestrator.createDeployment(icmp)
        DEPLOYMENTS.add(icmp)
        */
        for (Deployment d : DEPLOYMENTS) {
            assert Services.waitForDeployment(d)
        }
    }

    def cleanupSpec() {
        for (Deployment deployment : DEPLOYMENTS) {
            orchestrator.deleteDeployment(
                    deployment.getName(),
                    deployment.getNamespace(),
                    deployment.getExposeAsService()
            )
        }
    }

    @Unroll
    @Category([BAT, NetworkFlowVisualization])
    def "Verify connections can be detected: #protocol"() {
        given:
        "Two deployments, A and B, where B communicates to A via #protocol"
        String targetUid = DEPLOYMENTS.find { it.name == targetDeployment }?.deploymentUid
        assert targetUid != null
        String sourceUid = DEPLOYMENTS.find { it.name == sourceDeployment }?.deploymentUid
        assert sourceUid != null

        expect:
        "Check for edge in entwork graph"
        println "Checking for edge between ${sourceDeployment} and ${targetDeployment}"
        List<NetworkGraphOuterClass.NetworkEdge> edges = checkForEdge(sourceUid, targetUid)
        assert edges
        assert edges.get(0).protocol == protocol
        assert DEPLOYMENTS.find { it.name == targetDeployment }?.ports?.keySet()?.contains(edges.get(0).port)

        where:
        "Data is:"

        sourceDeployment     | targetDeployment      | protocol
        UDPCONNECTIONSOURCE  | UDPCONNECTIONTARGET   | NetworkEnums.L4Protocol.L4_PROTOCOL_UDP
        TCPCONNECTIONSOURCE  | TCPCONNECTIONTARGET   | NetworkEnums.L4Protocol.L4_PROTOCOL_TCP
        //ICMPCONNECTIONSOURCE | NGINXCONNECTIONTARGET | NetworkEnums.L4Protocol.L4_PROTOCOL_ICMP
    }

    @Category([NetworkFlowVisualization])
    def "Verify connections with short consistent intervals between 2 deployments"() {
        given:
        "Two deployments, A and B, where B communicates to A in short consistent intervals"
        String targetUid = DEPLOYMENTS.find { it.name == NGINXCONNECTIONTARGET }?.deploymentUid
        assert targetUid != null
        String sourceUid = DEPLOYMENTS.find { it.name == SHORTCONSISTENTSOURCE }?.deploymentUid
        assert sourceUid != null

        when:
        "Check for edge in entwork graph"
        println "Checking for edge between ${SHORTCONSISTENTSOURCE} and ${NGINXCONNECTIONTARGET}"
        List<NetworkGraphOuterClass.NetworkEdge> edges = checkForEdge(sourceUid, targetUid)
        assert edges

        then:
        "Wait for collector update and fetch graph again to confirm short interval connections remain"
        assert waitForEdgeUpdate(edges.get(0), 90)
    }

    @Category([NetworkFlowVisualization])
    def "Verify no connections between 2 deployments"() {
        given:
        "Two deployments, A and B, where neither communicates to the other"
        String targetUid = DEPLOYMENTS.find { it.name == NGINXCONNECTIONTARGET }?.deploymentUid
        assert targetUid != null
        String sourceUid = DEPLOYMENTS.find { it.name == NOCONNECTIONSOURCE }?.deploymentUid
        assert sourceUid != null

        expect:
        "Assert connection states"
        println "Checking for NO edge between ${NOCONNECTIONSOURCE} and ${NGINXCONNECTIONTARGET}"
        assert !checkForEdge(sourceUid, targetUid, null, 30)
    }

    @Category([NetworkFlowVisualization])
    def "Verify one-time connections show at first, but do not appear again"() {
        given:
        "Two deployments, A and B, where B communicates to A a single time during initial deployment"
        String targetUid = DEPLOYMENTS.find { it.name == NGINXCONNECTIONTARGET }?.deploymentUid
        assert targetUid != null
        String sourceUid = DEPLOYMENTS.find { it.name == SINGLECONNECTIONSOURCE }?.deploymentUid
        assert sourceUid != null

        when:
        "Check for edge in entwork graph"
        println "Checking for edge between ${SINGLECONNECTIONSOURCE} and ${NGINXCONNECTIONTARGET}"
        List<NetworkGraphOuterClass.NetworkEdge> edges = checkForEdge(sourceUid, targetUid)
        assert edges

        then:
        "Wait for collector update and fetch graph again to confirm connection dropped"
        assert !waitForEdgeUpdate(edges.get(0))
    }

    @Category([NetworkFlowVisualization])
    def "Verify connections between two deployments on 2 separate ports shows both edges in the graph"() {
        given:
        "Two deployments, A and B, where B communicates to A on 2 different ports"
        String targetUid = DEPLOYMENTS.find { it.name == TCPCONNECTIONTARGET }?.deploymentUid
        assert targetUid != null
        String sourceUid = DEPLOYMENTS.find { it.name == MULTIPLEPORTSCONNECTION }?.deploymentUid
        assert sourceUid != null

        when:
        "Check for edge in entwork graph"
        List<NetworkGraphOuterClass.NetworkEdge> edges = checkForEdge(sourceUid, targetUid)
        assert edges

        then:
        "Assert that there are 2 connection edges"
        assert edges.size() == 2
    }

    @Category([NetworkFlowVisualization])
    def "Verify cluster updates can block flow connections from showing"() {
        given:
        "Two deployments, A and B, where B communicates to A"
        String targetUid = DEPLOYMENTS.find { it.name == NGINXCONNECTIONTARGET }?.deploymentUid
        assert targetUid != null
        String sourceUid = DEPLOYMENTS.find { it.name == SHORTCONSISTENTSOURCE }?.deploymentUid
        assert sourceUid != null

        when:
        "apply network policy to block ingress to A"
        NetworkPolicy policy = new NetworkPolicy("deny-all-traffic-to-a")
                .setNamespace("qa")
                .addPodSelector(["app":NGINXCONNECTIONTARGET])
                .addPolicyType(NetworkPolicyTypes.INGRESS)
        def policyId = orchestrator.applyNetworkPolicy(policy)
        println "Sleeping 60s to allow policy to propagate and flows to update after propagation"
        sleep 60000

        and:
        "Check for original edge in network graph"
        println "Checking for edge between ${SHORTCONSISTENTSOURCE} and ${NGINXCONNECTIONTARGET}"
        List<NetworkGraphOuterClass.NetworkEdge> edges = checkForEdge(sourceUid, targetUid)
        assert edges

        then:
        "make sure edge does not get updated"
        assert !waitForEdgeUpdate(edges.get(0))

        cleanup:
        "remove policy"
        if (policyId != null) {
            orchestrator.deleteNetworkPolicy(policy)
        }
    }

    @Category([NetworkFlowVisualization])
    def "Verify deployment connecting to external does not generate edge"() {
        given:
        "One deployments, A, where A communicates to google.com"
        String sourceUid = DEPLOYMENTS.find { it.name == EXTERNALCONNECTIONSOURCE }?.deploymentUid
        assert sourceUid != null

        expect:
        "Check for edge in network graph"
        println "Checking for NO edge between ${EXTERNALCONNECTIONSOURCE} and www.google.com"
        assert checkForEdge(sourceUid, null, null, 30)?.size() <= 1
    }

    @Category([NetworkFlowVisualization])
    def "Verify edge timestamps are never in the future, or before start of flow tests"() {
        given:
        "Get current state of edges and current timestamp"
        NetworkGraph currentGraph = NetworkGraphService.getNetworkGraph()
        Long time = System.currentTimeMillis()
        Timestamp currentTimestamp =  Timestamp.newBuilder().setSeconds(time / 1000 as Long)
                .setNanos((int) ((time % 1000) * 1000000)).build()

        expect:
        "Check timestamp for each edge"
        for (NetworkGraphOuterClass.NetworkEdge edge : currentGraph.edgesList) {
            assert edge.lastActiveTimestamp.seconds <= currentTimestamp.seconds
            assert edge.lastActiveTimestamp.seconds >= testStartTime.seconds
        }
    }

    private checkForEdge(String sourceId, String targetId, Timestamp since = null, int timeoutSeconds = 90) {
        int intervalSeconds = 1
        int waitTime
        def startTime = System.currentTimeMillis()
        for (waitTime = 0; waitTime <= timeoutSeconds / intervalSeconds; waitTime++) {
            NetworkGraphOuterClass.NetworkGraph graph = NetworkGraphService.getNetworkGraph(since)
            List<NetworkGraphOuterClass.NetworkEdge> edges = graph.edgesList.findAll {
                targetId != null ?
                        it.source == sourceId && it.target == targetId :
                        it.source == sourceId
            }
            if (edges != null && edges.size() > 0) {
                println "Found source -> target in graph after ${(System.currentTimeMillis() - startTime) / 1000}s"
                return edges
            }
            sleep intervalSeconds * 1000
        }
        println "SR did not detect the edge in Network Flow graph"
        return null
    }

    private waitForEdgeUpdate(NetworkGraphOuterClass.NetworkEdge edge, int timeoutSeconds = 60) {
        int intervalSeconds = 1
        int waitTime
        def startTime = System.currentTimeMillis()
        for (waitTime = 0; waitTime <= timeoutSeconds / intervalSeconds; waitTime++) {
            NetworkGraph graph = NetworkGraphService.getNetworkGraph()
            NetworkGraphOuterClass.NetworkEdge newEdge = graph.edgesList.find {
                it.source == edge.source && it.target == edge.target
            }
            if (newEdge != null && newEdge.lastActiveTimestamp?.seconds > edge.lastActiveTimestamp.seconds) {
                println "Found updated edge in graph after ${(System.currentTimeMillis() - startTime) / 1000}s"
                return newEdge
            }
            sleep intervalSeconds * 1000
        }
        println "SR did not detect updated edge in Network Flow graph"
        return null
    }
}
