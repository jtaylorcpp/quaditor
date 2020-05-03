package quaditor

type Quaditor struct {
	authenticator Authenticator
	authorizer    Authorizer
	auditor       Auditor
}

func (q *Quaditor) Authenticate(metadata map[string]interface{}) (bool, error) {
	quads, approved, err := q.authenticator.Authenticate(metadata)
	quadErr := q.Publish(quads)

	if err != nil {
		return approved, err
	}

	if quadErr != nil {
		return approved, quadErr
	}

	return approved, nil
}

func (q *Quaditor) MintCredentials(metadata map[string]interface{}) (interface{}, error) {
	quads, creds, err := q.authenticator.MintCredentials(metadata)
	quadErr := q.Publish(quads)

	if err != nil {
		return creds, err
	}

	if quadErr != nil {
		return creds, quadErr
	}

	return creds, nil
}

func (q *Quaditor) Authorize(metadata map[string]interface{}) (bool, error) {
	quads, authorized, err := q.authorizer.Authorize(metadata)
	quadError := q.Publish(quads)
	if err != nil {
		return authorized, err
	}

	if quadError != nil {
		return authorized, quadError
	}

	return authorized, nil
}

func (q *Quaditor) Publish(quads []Quad) error {
	return q.auditor.Publish(quads)
}

func (q *Quaditor) Close() {
	q.authenticator.Close()
	q.authorizer.Close()
	q.auditor.Close()
}
