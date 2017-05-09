/*                          _       _
 *__      _____  __ ___   ___  __ _| |_ ___
 *\ \ /\ / / _ \/ _` \ \ / / |/ _` | __/ _ \
 * \ V  V /  __/ (_| |\ V /| | (_| | ||  __/
 *  \_/\_/ \___|\__,_| \_/ |_|\__,_|\__\___|
 *
 * Copyright © 2016 Weaviate. All rights reserved.
 * LICENSE: https://github.com/weaviate/weaviate/blob/master/LICENSE
 * AUTHOR: Bob van Luijt (bob@weaviate.com)
 * See www.weaviate.com for details
 * Contact: @weaviate_iot / yourfriends@weaviate.com
 */
 package commands




import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/weaviate/weaviate/models"
)

// WeaviateCommandsGetOKCode is the HTTP code returned for type WeaviateCommandsGetOK
const WeaviateCommandsGetOKCode int = 200

/*WeaviateCommandsGetOK Successful response.

swagger:response weaviateCommandsGetOK
*/
type WeaviateCommandsGetOK struct {

	/*
	  In: Body
	*/
	Payload *models.Command `json:"body,omitempty"`
}

// NewWeaviateCommandsGetOK creates WeaviateCommandsGetOK with default headers values
func NewWeaviateCommandsGetOK() *WeaviateCommandsGetOK {
	return &WeaviateCommandsGetOK{}
}

// WithPayload adds the payload to the weaviate commands get o k response
func (o *WeaviateCommandsGetOK) WithPayload(payload *models.Command) *WeaviateCommandsGetOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the weaviate commands get o k response
func (o *WeaviateCommandsGetOK) SetPayload(payload *models.Command) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *WeaviateCommandsGetOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// WeaviateCommandsGetNotFoundCode is the HTTP code returned for type WeaviateCommandsGetNotFound
const WeaviateCommandsGetNotFoundCode int = 404

/*WeaviateCommandsGetNotFound Successful query result but no resource was found.

swagger:response weaviateCommandsGetNotFound
*/
type WeaviateCommandsGetNotFound struct {
}

// NewWeaviateCommandsGetNotFound creates WeaviateCommandsGetNotFound with default headers values
func NewWeaviateCommandsGetNotFound() *WeaviateCommandsGetNotFound {
	return &WeaviateCommandsGetNotFound{}
}

// WriteResponse to the client
func (o *WeaviateCommandsGetNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
}

// WeaviateCommandsGetNotImplementedCode is the HTTP code returned for type WeaviateCommandsGetNotImplemented
const WeaviateCommandsGetNotImplementedCode int = 501

/*WeaviateCommandsGetNotImplemented Not (yet) implemented.

swagger:response weaviateCommandsGetNotImplemented
*/
type WeaviateCommandsGetNotImplemented struct {
}

// NewWeaviateCommandsGetNotImplemented creates WeaviateCommandsGetNotImplemented with default headers values
func NewWeaviateCommandsGetNotImplemented() *WeaviateCommandsGetNotImplemented {
	return &WeaviateCommandsGetNotImplemented{}
}

// WriteResponse to the client
func (o *WeaviateCommandsGetNotImplemented) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(501)
}
