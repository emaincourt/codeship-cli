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
)

func TestGetHeader(t *testing.T) {
	codeshipProvider := providers.CodeShipProvider{}
	header, _ := codeshipProvider.GetHeader()

	if header != "Input the index of the project :" {
		t.Errorf("Expected another output. Received : <%s>", header)
	}
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

	if len(projects) != 2 || projects[0] != "[0] Project#0" || projects[1] != "[1] Project#1" {
		t.Errorf("Expected another output. Received : <%v>", projects)
	}
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

	if len(builds) != 2 {
		t.Errorf("Builds length should equal 2. Received %d", len(builds))
	}
	if builds[0].StartedAt != time.Date(2018, 7, 1, 0, 0, 0, 0, time.UTC).Format("02/01/06 03:04:05") {
		t.Errorf("Dates format does not match the expected one. Received : %s", builds[0].StartedAt)
	}
	if builds[0].Status != fmt.Sprintf("[%s](fg-green)", "success") {
		t.Errorf("Status format does not match the expected one. Received : %s", builds[0].Status)
	}
	if builds[1].Status != fmt.Sprintf("[%s](fg-red)", "error") {
		t.Errorf("Status format does not match the expected one. Received : %s", builds[0].Status)
	}
}
