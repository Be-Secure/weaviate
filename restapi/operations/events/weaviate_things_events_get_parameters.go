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
 package events




import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime/middleware"

	strfmt "github.com/go-openapi/strfmt"
)

// NewWeaviateThingsEventsGetParams creates a new WeaviateThingsEventsGetParams object
// with the default values initialized.
func NewWeaviateThingsEventsGetParams() WeaviateThingsEventsGetParams {
	var ()
	return WeaviateThingsEventsGetParams{}
}

// WeaviateThingsEventsGetParams contains all the bound params for the weaviate things events get operation
// typically these are obtained from a http.Request
//
// swagger:parameters weaviate.things.events.get
type WeaviateThingsEventsGetParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request

	/*Unique ID of the event.
	  Required: true
	  In: path
	*/
	EventID string
	/*Unique ID of the thing.
	  Required: true
	  In: path
	*/
	ThingID string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls
func (o *WeaviateThingsEventsGetParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error
	o.HTTPRequest = r

	rEventID, rhkEventID, _ := route.Params.GetOK("eventId")
	if err := o.bindEventID(rEventID, rhkEventID, route.Formats); err != nil {
		res = append(res, err)
	}

	rThingID, rhkThingID, _ := route.Params.GetOK("thingId")
	if err := o.bindThingID(rThingID, rhkThingID, route.Formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *WeaviateThingsEventsGetParams) bindEventID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	o.EventID = raw

	return nil
}

func (o *WeaviateThingsEventsGetParams) bindThingID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	o.ThingID = raw

	return nil
}
