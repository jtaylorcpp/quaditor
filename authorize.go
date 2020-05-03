package quaditor

type Authorizer interface {
	Authorize(metadata map[string]interface{}) ([]Quad, bool, error)
	Close()
}
