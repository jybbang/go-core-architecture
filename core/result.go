package core

type Result struct {
	V interface{}
	E error
}

func (r *Result) ToActionResult() {
	// TODO
}
