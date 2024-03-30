package model

// ImplResponse defines an implementation response with error code and the associated body
type ImplResponse struct {
	Code int
	Body interface{}
}
