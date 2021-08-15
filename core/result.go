package core

type Result struct {
	V interface{}
	E error
}

var EmptyResult = Result{
	V: nil,
	E: ErrInternalServerError,
}

func (r *Result) ToActionResult() {
	// TODO
}
