package hostconfiguration

import (
	"context"

	"bitbucket.org/stack-rox/apollo/docker-bench/utils"
	"bitbucket.org/stack-rox/apollo/pkg/api/generated/api/v1"
)

type dockerUpdated struct{}

func (c *dockerUpdated) Definition() utils.Definition {
	return utils.Definition{
		CheckDefinition: v1.CheckDefinition{
			Name:        "CIS 1.3",
			Description: "Ensure Docker is up to date",
		}, Dependencies: []utils.Dependency{utils.InitDockerClient},
	}
}

func (c *dockerUpdated) Run() (result v1.CheckResult) {
	version, err := utils.DockerClient.ServerVersion(context.Background())
	if err != nil {
		utils.Note(&result)
		utils.AddNotef(&result, "Manual introspection will be req'd for docker version. Could not retrieve due to %+v", err)
		return
	}
	utils.Note(&result)
	utils.AddNotes(&result, "Docker server is currently running %v", version.Version)
	return
}

// NewDockerUpdated implements CIS-1.3
func NewDockerUpdated() utils.Check {
	return &dockerUpdated{}
}
