package setup

type Deployment string

const (
	Local      Deployment = "local"
	Test       Deployment = "test"
	Staging    Deployment = "staging"
	Production Deployment = "prod"
)
