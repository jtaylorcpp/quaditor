package quaditor

type Auditor interface {
	Publish([]Quad) error
	Query(...Query) error
	Close()
}
