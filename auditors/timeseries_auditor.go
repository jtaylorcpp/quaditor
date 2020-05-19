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

func (t *TimeSeriesAuditor) Query(queries ...quaditor.Query) ([]iterator.Path, error) {
	log.Println("running query")

	iteratorQueries := []iterator.Iterator{}
	for _, query := range queries {
		// check constraints
		constraints := []euler.Value{}
		if query.Constraint.Subject != "" {
			constraints = append(constraints, euler.Subject(query.Constraint.Subject))
		}
		if query.Constraint.Object != "" {
			constraints = append(constraints, euler.Object(query.Constraint.Object))
		}
		if query.Constraint.Predicate != "" {
			constraints = append(constraints, euler.Predicate(query.Constraint.Predicate))
		}
		if query.Constraint.Start != 0 {
			constraints = append(constraints, euler.Start(query.Constraint.Start))
		}
		if query.Constraint.End != 0 {
			constraints = append(constraints, euler.End(query.Constraint.End))
		}

		assign := []euler.Value{}
		if query.Assignment.Subject != "" {
			assign = append(assign, euler.Subject(query.Assignment.Subject))
		}
		if query.Assignment.Object != "" {
			assign = append(assign, euler.Object(query.Assignment.Object))
		}
		if query.Assignment.Predicate != "" {
			assign = append(assign, euler.Predicate(query.Assignment.Predicate))
		}
		if query.Assignment.Start != 0 {
			assign = append(assign, euler.Start(query.Assignment.Start))
		}
		if query.Assignment.End != 0 {
			assign = append(assign, euler.End(query.Assignment.End))
		}

		iteratorQueries = append(iteratorQueries, iterator.Constraint{constraints, assign})

	}

	paths, err := iterator.Run(t.graph,
		iteratorQueries...,
	)

	log.Println("number of returned paths: ", len(paths))
	if err != nil {
		log.Println("error running query: ", err.Error())
		return []iterator.Path{}, err
	} else {
		for idx, path := range paths {
			log.Printf("path %v: %#v\n", idx, path)
		}

		return paths, nil
	}
}
