package quaditor

type Quad struct {
	Subject   string `json:"subject"`
	Predicate string `json:"predicate"`
	Object    string `json:"object"`
	Start     uint64 `json:"start"`
	End       uint64 `json:"end"`
}

// Query holds 2 quads; the first is a constraint to query
//  the second Quad allows assignments of unknown value to be given a label;
//    this label will be used for follow on queries when re-used in an assignment block
type Query struct {
	Constraint Quad `json:"constraint"`
	Assignment Quad `json:"assign"`
}
