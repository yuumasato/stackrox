package containerimagesandbuild

import (
	"bitbucket.org/stack-rox/apollo/docker-bench/utils"
	"bitbucket.org/stack-rox/apollo/pkg/api/generated/api/v1"
)

type trustedBaseImagesBenchmark struct{}

func (c *trustedBaseImagesBenchmark) Definition() utils.Definition {
	return utils.Definition{
		CheckDefinition: v1.CheckDefinition{
			Name:        "CIS 4.2",
			Description: "Ensure that containers use trusted base images",
		}, Dependencies: []utils.Dependency{utils.InitContainers},
	}
}

func (c *trustedBaseImagesBenchmark) Run() (result v1.CheckResult) {
	utils.Note(&result)
	utils.AddNotes(&result, "Verification of trusted base images requires user specification")
	return
}

// NewTrustedBaseImagesBenchmark implements CIS-4.2
func NewTrustedBaseImagesBenchmark() utils.Check {
	return &trustedBaseImagesBenchmark{}
}
