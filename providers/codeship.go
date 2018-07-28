package providers

import (
	"context"
	"fmt"
	"strings"

	codeship "github.com/codeship/codeship-go"
)

// CodeShipAPI interfaces the Codeship API
type CodeShipAPI interface {
	ListProjects(ctx context.Context, opts ...codeship.PaginationOption) (codeship.ProjectList, codeship.Response, error)
	ListBuilds(ctx context.Context, projectUUID string, opts ...codeship.PaginationOption) (codeship.BuildList, codeship.Response, error)
}

// CodeShipProvider refers to a codeship provider
type CodeShipProvider struct {
	Provider

	Context context.Context
	API     CodeShipAPI

	Projects []codeship.Project
}

// NewCodeShipProviderFromCredentials instantiates a new CodeShipProvider from password, username and organization
func NewCodeShipProviderFromCredentials(user string, password string, organization string) (*CodeShipProvider, error) {
	ctx := context.Background()
	auth := codeship.NewBasicAuth(user, password)

	client, err := codeship.New(auth)
	if err != nil {
		return nil, err
	}

	org, err := client.Organization(ctx, organization)
	if err != nil {
		return nil, err
	}

	return &CodeShipProvider{
		Context: ctx,
		API:     org,
	}, nil
}

// GetHeader gets the header to display for Codeship
func (c *CodeShipProvider) GetHeader() (string, error) {
	return "Input the index of the project :", nil
}

// GetProjectsList gets the projects list
func (c *CodeShipProvider) GetProjectsList() ([]string, error) {
	res, _, err := c.API.ListProjects(c.Context)
	if err != nil {
		return nil, err
	}

	c.Projects = res.Projects

	var strs []string
	for index, project := range res.Projects {
		strs = append(strs, fmt.Sprintf("[%d] %s", index, project.Name))
	}

	return strs, nil
}

// GetProjectIDFromIndex gets the uuid of a project from a given index out of the total projects array length
func (c *CodeShipProvider) GetProjectIDFromIndex(index int) (string, error) {
	return c.Projects[index].UUID, nil
}

// GetBuildsList gets the builds list
func (c *CodeShipProvider) GetBuildsList(uuid string) ([]Build, error) {
	res, _, err := c.API.ListBuilds(c.Context, uuid)
	if err != nil {
		return nil, err
	}

	var builds []Build
	for _, build := range res.Builds {
		commitMessage := build.CommitMessage
		if len(commitMessage) > 70 {
			commitMessage = build.CommitMessage[:70]
		}

		status := fmt.Sprintf("[%s](fg-cyan)", build.Status)
		if strings.Contains(build.Status, "success") {
			status = fmt.Sprintf("[%s](fg-green)", build.Status)
		}
		if strings.Contains(build.Status, "error") {
			status = fmt.Sprintf("[%s](fg-red)", build.Status)
		}

		builds = append(builds, Build{
			StartedAt:     build.AllocatedAt.Format("02/01/06 03:04:05"),
			FinishedAt:    build.FinishedAt.Format("02/01/06 03:04:05"),
			CommitMessage: commitMessage,
			Status:        status,
			Username:      build.Username,
		})
	}

	return builds, nil
}
