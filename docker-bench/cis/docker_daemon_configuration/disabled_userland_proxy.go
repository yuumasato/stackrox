package dockerdaemonconfiguration

import (
	"bitbucket.org/stack-rox/apollo/docker-bench/utils"
	"bitbucket.org/stack-rox/apollo/pkg/api/generated/api/v1"
)

type disableUserlandProxyBenchmark struct{}

func (c *disableUserlandProxyBenchmark) Definition() utils.Definition {
	return utils.Definition{
		CheckDefinition: v1.CheckDefinition{
			Name:        "CIS 2.15",
			Description: "Ensure Userland Proxy is Disabled",
		}, Dependencies: []utils.Dependency{utils.InitDockerConfig},
	}
}

func (c *disableUserlandProxyBenchmark) Run() (result v1.CheckResult) {
	opts, ok := utils.DockerConfig["userland-proxy"]
	if !ok {
		utils.Warn(&result)
		utils.AddNotes(&result, "userland proxy is enabled by default")
		return
	}
	if opts.Matches("false") {
		utils.Warn(&result)
		utils.AddNotes(&result, "userland proxy is enabled")
		return
	}
	utils.Pass(&result)
	return

}

// NewDisableUserlandProxyBenchmark implements CIS-2.15
func NewDisableUserlandProxyBenchmark() utils.Check {
	return &disableUserlandProxyBenchmark{}
}
