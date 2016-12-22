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
 * See package.json for author and maintainer info
 * Contact: @weaviate_iot / yourfriends@weaviate.com
 */
 package authorized_apps




import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/weaviate/weaviate/core/models"
)

/*WeaveAuthorizedAppsListOK Successful response

swagger:response weaveAuthorizedAppsListOK
*/
type WeaveAuthorizedAppsListOK struct {

	// In: body
	Payload *models.AuthorizedAppsListResponse `json:"body,omitempty"`
}

// NewWeaveAuthorizedAppsListOK creates WeaveAuthorizedAppsListOK with default headers values
func NewWeaveAuthorizedAppsListOK() *WeaveAuthorizedAppsListOK {
	return &WeaveAuthorizedAppsListOK{}
}

// WithPayload adds the payload to the weave authorized apps list o k response
func (o *WeaveAuthorizedAppsListOK) WithPayload(payload *models.AuthorizedAppsListResponse) *WeaveAuthorizedAppsListOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the weave authorized apps list o k response
func (o *WeaveAuthorizedAppsListOK) SetPayload(payload *models.AuthorizedAppsListResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *WeaveAuthorizedAppsListOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		if err := producer.Produce(rw, o.Payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
