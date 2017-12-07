package dockersecurityoperations

import (
	"bitbucket.org/stack-rox/apollo/docker-bench/utils"
	"bitbucket.org/stack-rox/apollo/pkg/api/generated/api/v1"
)

type containerSprawlBenchmark struct{}

func (c *containerSprawlBenchmark) Definition() utils.Definition {
	return utils.Definition{
		CheckDefinition: v1.CheckDefinition{
			Name:        "CIS 6.2",
			Description: "Ensure container sprawl is avoided",
		}, Dependencies: []utils.Dependency{utils.InitContainers},
	}
}

func (c *containerSprawlBenchmark) Run() (result v1.CheckResult) {
	utils.Info(&result)
	utils.AddNotef(&result, "There are %v containers in use out of %v", len(utils.ContainersRunning), len(utils.ContainersAll))
	return
}

// NewContainerSprawlBenchmark implements CIS-6.2
func NewContainerSprawlBenchmark() utils.Check {
	return &containerSprawlBenchmark{}
}
