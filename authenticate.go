package quaditor

type Authenticator interface {
	Authenticate(metadata map[string]interface{}) ([]Quad, bool, error)
	MintCredentials(metadata map[string]interface{}) ([]Quad, interface{}, error)
	Close()
}
