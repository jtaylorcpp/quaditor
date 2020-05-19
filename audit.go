package quaditor

import (
	"pault.ag/go/euler/iterator"
)

type Auditor interface {
	Publish([]Quad) error
	Query(...Query) ([]iterator.Path, error)
	Close()
}
