package containerruntime

import (
	"strings"

	"bitbucket.org/stack-rox/apollo/docker-bench/utils"
	"bitbucket.org/stack-rox/apollo/pkg/api/generated/api/v1"
)

type privilegedDockerExecBenchmark struct{}

func (c *privilegedDockerExecBenchmark) Definition() utils.Definition {
	return utils.Definition{
		CheckDefinition: v1.CheckDefinition{
			Name:        "CIS 5.22",
			Description: "Ensure docker exec commands are not used with privileged option",
		}, Dependencies: []utils.Dependency{utils.InitContainers},
	}
}

func (c *privilegedDockerExecBenchmark) Run() (result v1.CheckResult) {
	utils.Pass(&result)
	auditLog, err := utils.ReadFile("/var/log/audit/audit.log")
	if err != nil {
		utils.Warn(&result)
		utils.AddNotef(&result, "Error reading /var/log/audit/audit.log: %+v", err)
		return
	}
	lines := strings.Split(auditLog, "\n")
	for _, line := range lines {
		if strings.Contains(line, "exec") && strings.Contains(line, "privileged") {
			utils.Warn(&result)
			utils.AddNotef(&result, "docker exec was used with the --privileged option: %v", line)
		}
	}
	return
}

// NewPrivilegedDockerExecBenchmark implements CIS-5.22
func NewPrivilegedDockerExecBenchmark() utils.Check {
	return &privilegedDockerExecBenchmark{}
}
