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
 package model_manifests




import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/weaviate/weaviate/models"
)

// WeaviateModelManifestsCreateAcceptedCode is the HTTP code returned for type WeaviateModelManifestsCreateAccepted
const WeaviateModelManifestsCreateAcceptedCode int = 202

/*WeaviateModelManifestsCreateAccepted Successfully received.

swagger:response weaviateModelManifestsCreateAccepted
*/
type WeaviateModelManifestsCreateAccepted struct {

	/*
	  In: Body
	*/
	Payload *models.ModelManifest `json:"body,omitempty"`
}

// NewWeaviateModelManifestsCreateAccepted creates WeaviateModelManifestsCreateAccepted with default headers values
func NewWeaviateModelManifestsCreateAccepted() *WeaviateModelManifestsCreateAccepted {
	return &WeaviateModelManifestsCreateAccepted{}
}

// WithPayload adds the payload to the weaviate model manifests create accepted response
func (o *WeaviateModelManifestsCreateAccepted) WithPayload(payload *models.ModelManifest) *WeaviateModelManifestsCreateAccepted {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the weaviate model manifests create accepted response
func (o *WeaviateModelManifestsCreateAccepted) SetPayload(payload *models.ModelManifest) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *WeaviateModelManifestsCreateAccepted) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(202)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// WeaviateModelManifestsCreateNotImplementedCode is the HTTP code returned for type WeaviateModelManifestsCreateNotImplemented
const WeaviateModelManifestsCreateNotImplementedCode int = 501

/*WeaviateModelManifestsCreateNotImplemented Not (yet) implemented.

swagger:response weaviateModelManifestsCreateNotImplemented
*/
type WeaviateModelManifestsCreateNotImplemented struct {
}

// NewWeaviateModelManifestsCreateNotImplemented creates WeaviateModelManifestsCreateNotImplemented with default headers values
func NewWeaviateModelManifestsCreateNotImplemented() *WeaviateModelManifestsCreateNotImplemented {
	return &WeaviateModelManifestsCreateNotImplemented{}
}

// WriteResponse to the client
func (o *WeaviateModelManifestsCreateNotImplemented) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(501)
}
