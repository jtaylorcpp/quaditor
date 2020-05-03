package quaditor

type Auditor interface {
	Publish([]Quad) error
	Close()
}
