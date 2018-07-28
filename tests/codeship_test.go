package providers_test

import (
	"context"
	"testing"

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
	}, codeship.Response{}, nil)

	projects, _ := codeshipProvider.GetProjectsList()

	if len(projects) != 2 || projects[0] != "[0] Project#0" || projects[1] != "[1] Project#1" {
		t.Errorf("Expected another output. Received : <%v>", projects)
	}
}
