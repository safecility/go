package setup

type Deployment string

const (
	LocalTest  Deployment = "local-test"
	Local      Deployment = "local"
	Test       Deployment = "test"
	Staging    Deployment = "staging"
	Production Deployment = "prod"
)
