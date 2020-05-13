package quaditor

type Quaditor struct {
	auditor Auditor
}

func (q *Quaditor) Publish(quads []Quad) error {
	return q.auditor.Publish(quads)
}

func (q *Quaditor) Close() {
	q.auditor.Close()
}
