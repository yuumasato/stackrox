package notifiers

import (
	"fmt"
	"testing"
	"text/template"

	types2 "github.com/gogo/protobuf/types"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/fixtures"
	"github.com/stackrox/rox/pkg/fixtures/fixtureconsts"
	"github.com/stackrox/rox/pkg/images/types"
	mitreDataStore "github.com/stackrox/rox/pkg/mitre/datastore"
	"github.com/stackrox/rox/pkg/timeutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	expectedFormattedDeploymentAlert = `Alert ID: ` + fixtureconsts.Alert1 + `
Alert URL: https://localhost:8080/main/violations/` + fixtureconsts.Alert1 + `
Time (UTC): 2021-01-20 22:42:02
Severity: Low

Violations:
	 - Deployment is affected by 'CVE-2017-15804'
	 - Deployment is affected by 'CVE-2017-15670'
	 - This is a kube event violation
		 - pod : nginx
		 - container : nginx
	 - This is a process violation

Policy Definition:

	Description:
	 - Alert if the container contains vulnerabilities

	Rationale:
	 - This is the rationale

	Remediation:
	 - This is the remediation

	Policy Criteria:

		Section Unnamed :

			- Image Registry: docker.io
			- Image Remote: r/.*stackrox/nginx.*
			- Image Tag: 1.10
			- Image Age: 30
			- Dockerfile Line: VOLUME=/etc/*
			- CVE: CVE-1234
			- Image Component: berkeley*=.*
			- Image Scan Age: 10
			- Environment Variable: UNSET=key=value
			- Volume Name: name
			- Volume Type: nfs
			- Volume Destination: /etc/network
			- Volume Source: 10.0.0.1/export
			- Writable Mounted Volume: false
			- Port: 8080
			- Protocol: tcp
			- Privileged: true
			- CVSS: >= 5.000000
			- Drop Capabilities: DROP1 OR DROP2
			- Add Capabilities: ADD1 OR ADD2

Deployment:
	 - ID: ` + fixtureconsts.Deployment1 + `
	 - Name: nginx_server
	 - Cluster: prod cluster
	 - ClusterId: ` + fixtureconsts.Cluster1 + `
	 - Namespace: stackrox
	 - Images: docker.io/library/nginx:1.10@sha256:SHA1
`
	expectedFormattedDeploymentAlertWithMitre = `Alert ID: ` + fixtureconsts.Alert1 + `
Alert URL: https://localhost:8080/main/violations/` + fixtureconsts.Alert1 + `
Time (UTC): 2021-01-20 22:42:02
Severity: Low

Violations:
	 - Deployment is affected by 'CVE-2017-15804'
	 - Deployment is affected by 'CVE-2017-15670'
	 - This is a kube event violation
		 - pod : nginx
		 - container : nginx
	 - This is a process violation

Policy Definition:

	Description:
	 - Alert if the container contains vulnerabilities

	Rationale:
	 - This is the rationale

	Remediation:
	 - This is the remediation

	MITRE ATT&CK:
	 - Tactic: Initial Access ( TA0001 )
		 - Techniques:
			 - Valid Accounts ( T1078 )
			 - Valid Accounts: Default Accounts ( T1078.001 )
	 - Tactic: Persistence ( TA0003 )

	Policy Criteria:

		Section Unnamed :

			- Image Registry: docker.io
			- Image Remote: r/.*stackrox/nginx.*
			- Image Tag: 1.10
			- Image Age: 30
			- Dockerfile Line: VOLUME=/etc/*
			- CVE: CVE-1234
			- Image Component: berkeley*=.*
			- Image Scan Age: 10
			- Environment Variable: UNSET=key=value
			- Volume Name: name
			- Volume Type: nfs
			- Volume Destination: /etc/network
			- Volume Source: 10.0.0.1/export
			- Writable Mounted Volume: false
			- Port: 8080
			- Protocol: tcp
			- Privileged: true
			- CVSS: >= 5.000000
			- Drop Capabilities: DROP1 OR DROP2
			- Add Capabilities: ADD1 OR ADD2

Deployment:
	 - ID: ` + fixtureconsts.Deployment1 + `
	 - Name: nginx_server
	 - Cluster: prod cluster
	 - ClusterId: ` + fixtureconsts.Cluster1 + `
	 - Namespace: stackrox
	 - Images: docker.io/library/nginx:1.10@sha256:SHA1
`
	expectedFormatImageAlert = `Alert ID: ` + fixtureconsts.Alert1 + `
Alert URL: https://localhost:8080/main/vulnerability-management/image/sha256:SHA2
Time (UTC): 2021-01-20 22:42:02
Severity: Low

Violations:
	 - Deployment is affected by 'CVE-2017-15804'
	 - Deployment is affected by 'CVE-2017-15670'
	 - This is a kube event violation
		 - pod : nginx
		 - container : nginx
	 - This is a process violation

Policy Definition:

	Description:
	 - Alert if the container contains vulnerabilities

	Rationale:
	 - This is the rationale

	Remediation:
	 - This is the remediation

	Policy Criteria:

		Section Unnamed :

			- Image Registry: docker.io
			- Image Remote: r/.*stackrox/nginx.*
			- Image Tag: 1.10
			- Image Age: 30
			- Dockerfile Line: VOLUME=/etc/*
			- CVE: CVE-1234
			- Image Component: berkeley*=.*
			- Image Scan Age: 10
			- Environment Variable: UNSET=key=value
			- Volume Name: name
			- Volume Type: nfs
			- Volume Destination: /etc/network
			- Volume Source: 10.0.0.1/export
			- Writable Mounted Volume: false
			- Port: 8080
			- Protocol: tcp
			- Privileged: true
			- CVSS: >= 5.000000
			- Drop Capabilities: DROP1 OR DROP2
			- Add Capabilities: ADD1 OR ADD2

Image:
	 - Name: stackrox.io/srox/mongo:latest
`
	expectedFormatNetworkAlert = `Alert ID: ` + fixtureconsts.Alert1 + `
Alert URL: https://localhost:8080/main/violations/` + fixtureconsts.Alert1 + `
Time (UTC): 2021-01-20 22:42:02
Severity: High

Violations:
	 - Unexpected network flow found in deployment. Source name: 'central'. Destination name: 'External Entities'. Destination port: '9'. Protocol: 'L4_PROTOCOL_UDP'.
	 - Unexpected network flow found in deployment. Source name: 'central'. Destination name: 'scanner'. Destination port: '8080'. Protocol: 'L4_PROTOCOL_TCP'.

Policy Definition:

	Description:
	 - This policy generates a violation for the network flows that fall outside baselines for which 'alert on anomalous violations' is set.

	Rationale:
	 - The network baseline is a list of flows that are allowed, and once it is frozen, any flow outside that is a concern.

	Remediation:
	 - Evaluate this network flow. If deemed to be okay, add it to the baseline. If not, investigate further as required.

	Policy Criteria:

		Section Unnamed :

			- Unexpected Network Flow Detected: true

Deployment:
	 - ID: ` + fixtureconsts.Deployment1 + `
	 - Name: central
	 - Cluster: remote
	 - ClusterId: ` + fixtureconsts.Cluster1 + `
	 - Namespace: stackrox
	 - Images: docker.io/library/nginx:1.10@sha256:SHA1
`
	expectedFormattedResourceAlert = `Alert ID: ` + fixtureconsts.Alert1 + `
Alert URL: https://localhost:8080/main/violations/` + fixtureconsts.Alert1 + `
Time (UTC): 2021-01-20 22:42:02
Severity: Low

Violations:
	 - Access to secret "my-secret" in "cluster-id / stackrox"
		 - Kubernetes API Verb : CREATE
		 - username : test-user
		 - user groups : groupA, groupB
		 - resource : /api/v1/namespace/stackrox/secrets/my-secret
		 - user agent : oc/4.7.0 (darwin/amd64) kubernetes/c66c03f
		 - IP address : 192.168.0.1, 127.0.0.1
		 - impersonated username : central-service-account
		 - impersonated user groups : service-accounts, groupB

Policy Definition:

	Description:
	 - Alert if the container contains vulnerabilities

	Rationale:
	 - This is the rationale

	Remediation:
	 - This is the remediation

	Policy Criteria:

		Section Unnamed :

			- Kubernetes Resource: SECRETS
			- Kubernetes API Verb: CREATE

Resource:
	 - Name: my-secret
	 - Type: SECRETS
	 - Cluster: prod cluster
	 - ClusterId: ` + fixtureconsts.Cluster1 + `
	 - Namespace: stackrox
`
	expectedFormattedClusterRoleResourceAlert = `Alert ID: ` + fixtureconsts.Alert1 + `
Alert URL: https://localhost:8080/main/violations/` + fixtureconsts.Alert1 + `
Time (UTC): 2021-01-20 22:42:02
Severity: Low

Violations:
	 - Access to cluster role "my-cluster-role"
		 - Kubernetes API Verb : CREATE
		 - username : test-user
		 - user groups : groupA, groupB
		 - resource : /apis/rbac.authorization.k8s.io/v1/clusterroles/my-cluster-role
		 - user agent : oc/4.7.0 (darwin/amd64) kubernetes/c66c03f
		 - IP address : 192.168.0.1, 127.0.0.1
		 - impersonated username : central-service-account
		 - impersonated user groups : service-accounts, groupB

Policy Definition:

	Description:
	 - Alert if the container contains vulnerabilities

	Rationale:
	 - This is the rationale

	Remediation:
	 - This is the remediation

	Policy Criteria:

		Section Unnamed :

			- Kubernetes Resource: CLUSTER_ROLES
			- Kubernetes API Verb: CREATE

Resource:
	 - Name: my-cluster-role
	 - Type: CLUSTER_ROLES
	 - Cluster: prod cluster
	 - ClusterId: ` + fixtureconsts.Cluster3 + `
`
)

func TestFormatAlert(t *testing.T) {
	runFormatTest(t, fixtures.GetAlert(), expectedFormattedDeploymentAlert)

	imageAlert := fixtures.GetAlert()
	imageAlert.Entity = &storage.Alert_Image{Image: types.ToContainerImage(fixtures.GetImage())}
	runFormatTest(t, imageAlert, expectedFormatImageAlert)
}

func TestNetworkAlert(t *testing.T) {
	runFormatTest(t, fixtures.GetNetworkAlert(), expectedFormatNetworkAlert)
}

func TestResourceAlert(t *testing.T) {
	runFormatTest(t, fixtures.GetResourceAlert(), expectedFormattedResourceAlert)
	runFormatTest(t, fixtures.GetClusterResourceAlert(), expectedFormattedClusterRoleResourceAlert)
}

func TestFormatAlertWithMitre(t *testing.T) {
	runFormatTest(t, fixtures.GetAlertWithMitre(), expectedFormattedDeploymentAlertWithMitre)
}

func runFormatTest(t *testing.T, alert *storage.Alert, expectedFormattedAlert string) {
	funcMap := template.FuncMap{
		"header": func(s string) string {
			return fmt.Sprintf("\n%v\n", s)
		},
		"subheader": func(s string) string {
			return fmt.Sprintf("\n\t%v\n", s)
		},
		"line": func(s string) string {
			return fmt.Sprintf("%v\n", s)
		},
		"list": func(s string) string {
			return fmt.Sprintf("\t - %v\n", s)
		},
		"nestedList": func(s string) string {
			return fmt.Sprintf("\t\t - %v\n", s)
		},
		"section": func(s string) string {
			return fmt.Sprintf("\n\t\t%v\n", s)
		},
		"group": func(s string) string {
			return fmt.Sprintf("\n\t\t\t- %v", s)
		},
	}

	testFormat := func(alert *storage.Alert, expected string) {
		var err error
		alert.Time, err = types2.TimestampProto(timeutil.MustParse("2006-01-02 15:04:05", "2021-01-20 22:42:02"))
		require.NoError(t, err)
		formatted, err := FormatAlert(alert, AlertLink("https://localhost:8080", alert), funcMap, mitreDataStore.Singleton())
		require.NoError(t, err)
		assert.Equal(t, expected, formatted)
	}

	testFormat(alert, expectedFormattedAlert)
}

func TestSummaryForAlert(t *testing.T) {
	alertWithNoEntity := fixtures.GetAlert()
	alertWithNoEntity.Entity = nil

	cases := []struct {
		name            string
		alert           *storage.Alert
		expectedSummary string
	}{
		{
			name:            "Deployment alert summary",
			alert:           fixtures.GetAlert(),
			expectedSummary: "Deployment nginx_server (in cluster prod cluster) violates 'Vulnerable Container' Policy",
		},
		{
			name:            "Image alert summary",
			alert:           fixtures.GetImageAlert(),
			expectedSummary: "Image stackrox.io/srox/mongo:latest violates 'Vulnerable Container' Policy",
		},
		{
			name:            "Resource alert summary",
			alert:           fixtures.GetResourceAlert(),
			expectedSummary: "Policy 'Vulnerable Container' violated in cluster prod cluster",
		},
		{
			name:            "Network alert summary",
			alert:           fixtures.GetNetworkAlert(),
			expectedSummary: "Deployment central (in cluster remote) violates 'Unauthorized Network Flow' Policy",
		},
		{
			name:            "Unexpected entity alert summary",
			alert:           alertWithNoEntity,
			expectedSummary: "Policy 'Vulnerable Container' violated",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			summary := SummaryForAlert(c.alert)

			assert.Equal(t, c.expectedSummary, summary)
		})
	}
}
