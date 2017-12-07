package containerimagesandbuild

import (
	"bitbucket.org/stack-rox/apollo/docker-bench/utils"
	"bitbucket.org/stack-rox/apollo/pkg/api/generated/api/v1"
)

type setuidSetGidPermissionsBenchmark struct{}

func (c *setuidSetGidPermissionsBenchmark) Definition() utils.Definition {
	return utils.Definition{
		CheckDefinition: v1.CheckDefinition{
			Name:        "CIS 4.8",
			Description: "Ensure setuid and setgid permissions are removed in the images",
		}, Dependencies: []utils.Dependency{utils.InitImages},
	}
}

func (c *setuidSetGidPermissionsBenchmark) Run() (result v1.CheckResult) {
	utils.Note(&result)
	utils.AddNotes(&result, "Checking if setuid and setgid permissions are removed in the images is invasive and requires running every image")
	return
}

// NewSetuidSetGidPermissionsBenchmark implements CIS-4.8
func NewSetuidSetGidPermissionsBenchmark() utils.Check {
	return &setuidSetGidPermissionsBenchmark{}
}
