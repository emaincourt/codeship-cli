package providers_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	codeship "github.com/codeship/codeship-go"
	"github.com/emaincourt/codeship-cli/providers"
	"github.com/emaincourt/codeship-cli/tests/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGetHeader(t *testing.T) {
	codeshipProvider := providers.CodeShipProvider{}
	header, _ := codeshipProvider.GetHeader()

	assert.Equal(t, header, "Input the index of the project :")
}

func TestGetProjectsList(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockCodeshipAPI := mocks.NewMockCodeShipAPI(mockCtrl)

	codeshipProvider := providers.CodeShipProvider{
		Context: context.Background(),
		API:     mockCodeshipAPI,
	}
	mockCodeshipAPI.EXPECT().ListProjects(codeshipProvider.Context).Return(codeship.ProjectList{
		Projects: []codeship.Project{
			codeship.Project{Name: "Project#0"},
			codeship.Project{Name: "Project#1"},
		},
	}, codeship.Response{}, nil).Times(1)

	projects, _ := codeshipProvider.GetProjectsList()

	assert.Equal(t, len(projects), 2)
	assert.Equal(t, projects[0], "[0] Project#0")
	assert.Equal(t, projects[1], "[1] Project#1")
}

func TestGetGetBuildsList(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockCodeshipAPI := mocks.NewMockCodeShipAPI(mockCtrl)

	codeshipProvider := providers.CodeShipProvider{
		Context: context.Background(),
		API:     mockCodeshipAPI,
	}
	mockCodeshipAPI.EXPECT().ListBuilds(codeshipProvider.Context, "any-uuid").Return(codeship.BuildList{
		Builds: []codeship.Build{
			codeship.Build{
				AllocatedAt:   time.Date(2018, 7, 1, 0, 0, 0, 0, time.UTC),
				FinishedAt:    time.Date(2018, 7, 1, 1, 0, 0, 0, time.UTC),
				Status:        "success",
				CommitMessage: "commit message",
				Username:      "emaincourt",
			},
			codeship.Build{
				AllocatedAt:   time.Date(2018, 8, 1, 0, 0, 0, 0, time.UTC),
				FinishedAt:    time.Date(2018, 8, 1, 1, 0, 0, 0, time.UTC),
				Status:        "error",
				CommitMessage: "commit message",
				Username:      "emaincourt",
			},
		},
	}, codeship.Response{}, nil).Times(1)

	builds, _ := codeshipProvider.GetBuildsList("any-uuid")

	assert.Equal(t, len(builds), 2)
	assert.Equal(t, builds[0].StartedAt, time.Date(2018, 7, 1, 0, 0, 0, 0, time.UTC).Format("02/01/06 03:04:05"))
	assert.Equal(t, builds[0].Status, fmt.Sprintf("[%s](fg-green)", "success"))
	assert.Equal(t, builds[1].Status, fmt.Sprintf("[%s](fg-red)", "failed"))
}
