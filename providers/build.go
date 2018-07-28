package providers

// Build defines the structure of a common build, whoever the provider is
type Build struct {
	StartedAt     string
	FinishedAt    string
	CommitMessage string
	Status        string
	Username      string
}
