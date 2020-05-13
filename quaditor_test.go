package quaditor

import (
	"testing"
)

type TestAuditor struct {
	c chan Quad
}

func (a *TestAuditor) Publish(q []Quad) error {
	for _, quad := range q {
		a.c <- quad
	}
	return nil
}

func (a *TestAuditor) Query(q ...Query) error {
	return nil
}

func (a *TestAuditor) Close() {
	close(a.c)
}

func TestPublish(t *testing.T) {
	auditor := &TestAuditor{
		c: make(chan Quad, 100),
	}

	q := &Quaditor{
		auditor: auditor,
	}

	defer q.Close()

	q.Publish([]Quad{Quad{"s", "p", "o", 0, 0}})

	if len(q.auditor.(*TestAuditor).c) != 1 {
		t.Fatal("should be message in queue")
	}

	quad := <-q.auditor.(*TestAuditor).c

	t.Logf("quad: %#v\n", quad)
	if quad.Subject != "s" || quad.Object != "o" ||
		quad.Predicate != "p" || quad.Start != 0 || quad.End != 0 {
		t.Fatal("did not get right quad back")
	}

}
