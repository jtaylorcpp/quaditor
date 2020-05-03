package quaditor

type Quad struct {
	Subject   string `json:"subject"`
	Predicate string `json:"predicate"`
	Object    string `json:"object"`
	Start     uint64 `json:"start"`
	End       uint64 `json:"end"`
}
