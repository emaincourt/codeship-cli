package providers

// Provider defines the structure of a provider
type Provider interface {
	GetHeader() (string, error)
	GetProjectsList() ([]string, error)
	GetProjectIDFromIndex(index int) (string, error)
	GetBuildsList(id string) ([]Build, error)
}
