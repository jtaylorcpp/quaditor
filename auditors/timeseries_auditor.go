package auditors

import (
	"errors"
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/jtaylorcpp/quaditor"
	log "github.com/sirupsen/logrus"
	"pault.ag/go/euler"
	"pault.ag/go/euler/iterator"
	"pault.ag/go/euler/sql"
)

func NewTimeSeriesAuditor(backend, username, password, host, port string) (*TimeSeriesAuditor, error) {
	db, err := gorm.Open(
		backend,
		fmt.Sprintf("host=%s port=%s user=%s dbname=euler password=%s", host, port, username, password),
	)
	if err != nil {
		return nil, err
	}

	db = db.Debug()

	graph, err := sql.NewGraph(db)

	if err != nil {
		return nil, err
	}

	if err = graph.Init(); err != nil {
		return nil, err
	}

	return &TimeSeriesAuditor{
		graph: graph,
	}, nil
}

type TimeSeriesAuditor struct {
	graph *sql.Graph
}

func (t *TimeSeriesAuditor) Publish(quads []quaditor.Quad) error {
	errored := false
	for idx, quad := range quads {
		err := t.graph.Add(euler.NewQuad(quad.Subject, quad.Predicate, quad.Object, quad.Start, quad.End))
		if err != nil {
			errored = true
			log.Warningf("quad (%v, %#v) was not written with error: %#v\n", idx, quad, err.Error())
		} else {
			log.Infof("quad(%v, %#v) written\n", idx, quad)
		}
	}

	if errored {
		return errors.New("some quads not written")
	}
	return nil
}

func (t *TimeSeriesAuditor) Close() {}

func (t *TimeSeriesAuditor) Query() {
	log.Println("running query")
	paths, err := iterator.Run(
		t.graph,
		iterator.Constraint{
			[]euler.Value{
				//euler.Predicate("recieves-pub-key"),
				euler.Object("bob"),
			},
			[]euler.Value{
				//euler.Subject("bobs-friend"),
			},
		},
		iterator.Constraint{
			[]euler.Value{
				//euler.Predicate("sends-pub-key"),
				euler.Object("alice"),
			},
			[]euler.Value{
				//euler.Subject("bobs-friend"),
			},
		},
	)

	log.Println("number of returned paths: ", len(paths))
	if err != nil {
		log.Println("error running query: ", err.Error())
	} else {
		for idx, path := range paths {
			log.Printf("path %v: %#v\n", idx, path)
		}
	}
}
