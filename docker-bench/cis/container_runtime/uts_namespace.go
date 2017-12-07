package containerruntime

import (
	"bitbucket.org/stack-rox/apollo/docker-bench/utils"
	"bitbucket.org/stack-rox/apollo/pkg/api/generated/api/v1"
)

type utsNamespaceBenchmark struct{}

func (c *utsNamespaceBenchmark) Definition() utils.Definition {
	return utils.Definition{
		CheckDefinition: v1.CheckDefinition{
			Name:        "CIS 5.20",
			Description: "Ensure the host's UTS namespace is not shared",
		}, Dependencies: []utils.Dependency{utils.InitContainers},
	}
}

func (c *utsNamespaceBenchmark) Run() (result v1.CheckResult) {
	utils.Pass(&result)
	for _, container := range utils.ContainersRunning {
		if container.HostConfig.UTSMode.IsHost() {
			utils.Warn(&result)
			utils.AddNotef(&result, "Container %v has UTS mode set to host", container.ID)
		}
	}
	return
}

// NewUTSNamespaceBenchmark implements CIS-5.20
func NewUTSNamespaceBenchmark() utils.Check {
	return &utsNamespaceBenchmark{}
}
