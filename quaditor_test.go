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

func (a *TestAuditor) Close() {
	close(a.c)
}

type TestAuthx struct{}

func (a *TestAuthx) Authenticate(metadata map[string]interface{}) ([]Quad, bool, error) {
	user, userOk := metadata["user"]
	password, passOk := metadata["password"]

	if !userOk || !passOk {
		return []Quad{}, false, nil
	}

	userString, userOk := user.(string)
	passString, passOk := password.(string)

	if !userOk || !passOk {
		return []Quad{}, false, nil
	}

	if userString == "root" && passString == "root" {
		return []Quad{}, true, nil
	}

	return []Quad{}, false, nil

}
func (a *TestAuthx) MintCredentials(metadata map[string]interface{}) ([]Quad, interface{}, error) {
	user, userOk := metadata["user"]
	password, passOk := metadata["password"]

	if !userOk || !passOk {
		return []Quad{}, nil, nil
	}

	userString, userOk := user.(string)
	passString, passOk := password.(string)

	if !userOk || !passOk {
		return []Quad{}, false, nil
	}

	return []Quad{}, userString + passString, nil
}

func (a *TestAuthx) Authorize(metadata map[string]interface{}) ([]Quad, bool, error) {
	creds, credsOk := metadata["creds"]
	if !credsOk {
		return []Quad{}, false, nil
	}

	credsString, credsOk := creds.(string)
	if !credsOk {
		return []Quad{}, false, nil
	}

	if credsString != "rootroot" {
		return []Quad{}, false, nil
	}

	return []Quad{}, true, nil
}

func (a *TestAuthx) Close() {}

func TestQuaditor(t *testing.T) {
	auditor := &TestAuditor{
		c: make(chan Quad, 100),
	}

	authx := &TestAuthx{}

	q := &Quaditor{
		auditor:       auditor,
		authorizer:    authx,
		authenticator: authx,
	}

	defer q.Close()
}

func TestQuaditorLogin(t *testing.T) {
	auditor := &TestAuditor{
		c: make(chan Quad, 100),
	}

	authx := &TestAuthx{}

	q := &Quaditor{
		auditor:       auditor,
		authorizer:    authx,
		authenticator: authx,
	}

	defer q.Close()

	authd, err := q.Authenticate(map[string]interface{}{"user": "root", "password": "root"})
	if err != nil {
		t.Fatal(err.Error())
	}

	if !authd {
		t.Fatal("root:root should login")
	}

	authd, err = q.Authenticate(map[string]interface{}{"user": "l33t", "password": "password"})
	if err != nil {
		t.Fatal(err.Error())
	}

	if authd {
		t.Fatal("l33t:password shouldnt login")
	}
}

func TestQuaditorCreds(t *testing.T) {
	auditor := &TestAuditor{
		c: make(chan Quad, 100),
	}

	authx := &TestAuthx{}

	q := &Quaditor{
		auditor:       auditor,
		authorizer:    authx,
		authenticator: authx,
	}

	defer q.Close()

	authd, err := q.Authenticate(map[string]interface{}{"user": "root", "password": "root"})
	if err != nil {
		t.Fatal(err.Error())
	}

	if !authd {
		t.Fatal("root:root should login")
	}

	authd, err = q.Authenticate(map[string]interface{}{"user": "l33t", "password": "password"})
	if err != nil {
		t.Fatal(err.Error())
	}

	if authd {
		t.Fatal("l33t:password shouldnt login")
	}
}

func TestPublish(t *testing.T) {
	auditor := &TestAuditor{
		c: make(chan Quad, 100),
	}

	authx := &TestAuthx{}

	q := &Quaditor{
		auditor:       auditor,
		authorizer:    authx,
		authenticator: authx,
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
