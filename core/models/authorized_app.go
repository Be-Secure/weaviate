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

package models

import (
	strfmt "github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"

	"github.com/go-openapi/errors"
)

// AuthorizedApp authorized app
// swagger:model AuthorizedApp
type AuthorizedApp struct {

	// Android apps authorized under this project ID.
	AndroidApps []*AuthorizedAppAndroidAppsItems0 `json:"androidApps"`

	// The display name of the app.
	DisplayName string `json:"displayName,omitempty"`

	// An icon for the app.
	IconURL string `json:"iconUrl,omitempty"`

	// Identifies what kind of resource this is. Value: the fixed string "weave#authorizedApp".
	Kind *string `json:"kind,omitempty"`

	// Project ID.
	ProjectID string `json:"projectId,omitempty"`
}

// Validate validates this authorized app
func (m *AuthorizedApp) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateAndroidApps(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *AuthorizedApp) validateAndroidApps(formats strfmt.Registry) error {

	if swag.IsZero(m.AndroidApps) { // not required
		return nil
	}

	for i := 0; i < len(m.AndroidApps); i++ {

		if swag.IsZero(m.AndroidApps[i]) { // not required
			continue
		}

		if m.AndroidApps[i] != nil {

			if err := m.AndroidApps[i].Validate(formats); err != nil {
				return err
			}
		}

	}

	return nil
}

// AuthorizedAppAndroidAppsItems0 authorized app android apps items0
// swagger:model AuthorizedAppAndroidAppsItems0
type AuthorizedAppAndroidAppsItems0 struct {

	// Android certificate hash.
	CertificateHash string `json:"certificate_hash,omitempty"`

	// Android package name.
	PackageName string `json:"package_name,omitempty"`
}

// Validate validates this authorized app android apps items0
func (m *AuthorizedAppAndroidAppsItems0) Validate(formats strfmt.Registry) error {
	var res []error

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
