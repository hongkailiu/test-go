// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"
)

// GetUsersOKCode is the HTTP code returned for type GetUsersOK
const GetUsersOKCode int = 200

/*GetUsersOK A JSON array of user names

swagger:response getUsersOK
*/
type GetUsersOK struct {

	/*
	  In: Body
	*/
	Payload []string `json:"body,omitempty"`
}

// NewGetUsersOK creates GetUsersOK with default headers values
func NewGetUsersOK() *GetUsersOK {

	return &GetUsersOK{}
}

// WithPayload adds the payload to the get users o k response
func (o *GetUsersOK) WithPayload(payload []string) *GetUsersOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get users o k response
func (o *GetUsersOK) SetPayload(payload []string) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetUsersOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	payload := o.Payload
	if payload == nil {
		payload = make([]string, 0, 50)
	}

	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}
