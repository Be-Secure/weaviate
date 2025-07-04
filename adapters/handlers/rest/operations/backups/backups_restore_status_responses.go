//                           _       _
// __      _____  __ ___   ___  __ _| |_ ___
// \ \ /\ / / _ \/ _` \ \ / / |/ _` | __/ _ \
//  \ V  V /  __/ (_| |\ V /| | (_| | ||  __/
//   \_/\_/ \___|\__,_| \_/ |_|\__,_|\__\___|
//
//  Copyright © 2016 - 2025 Weaviate B.V. All rights reserved.
//
//  CONTACT: hello@weaviate.io
//

// Code generated by go-swagger; DO NOT EDIT.

package backups

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/weaviate/weaviate/entities/models"
)

// BackupsRestoreStatusOKCode is the HTTP code returned for type BackupsRestoreStatusOK
const BackupsRestoreStatusOKCode int = 200

/*
BackupsRestoreStatusOK Backup restoration status successfully returned

swagger:response backupsRestoreStatusOK
*/
type BackupsRestoreStatusOK struct {

	/*
	  In: Body
	*/
	Payload *models.BackupRestoreStatusResponse `json:"body,omitempty"`
}

// NewBackupsRestoreStatusOK creates BackupsRestoreStatusOK with default headers values
func NewBackupsRestoreStatusOK() *BackupsRestoreStatusOK {

	return &BackupsRestoreStatusOK{}
}

// WithPayload adds the payload to the backups restore status o k response
func (o *BackupsRestoreStatusOK) WithPayload(payload *models.BackupRestoreStatusResponse) *BackupsRestoreStatusOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the backups restore status o k response
func (o *BackupsRestoreStatusOK) SetPayload(payload *models.BackupRestoreStatusResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *BackupsRestoreStatusOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// BackupsRestoreStatusUnauthorizedCode is the HTTP code returned for type BackupsRestoreStatusUnauthorized
const BackupsRestoreStatusUnauthorizedCode int = 401

/*
BackupsRestoreStatusUnauthorized Unauthorized or invalid credentials.

swagger:response backupsRestoreStatusUnauthorized
*/
type BackupsRestoreStatusUnauthorized struct {
}

// NewBackupsRestoreStatusUnauthorized creates BackupsRestoreStatusUnauthorized with default headers values
func NewBackupsRestoreStatusUnauthorized() *BackupsRestoreStatusUnauthorized {

	return &BackupsRestoreStatusUnauthorized{}
}

// WriteResponse to the client
func (o *BackupsRestoreStatusUnauthorized) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(401)
}

// BackupsRestoreStatusForbiddenCode is the HTTP code returned for type BackupsRestoreStatusForbidden
const BackupsRestoreStatusForbiddenCode int = 403

/*
BackupsRestoreStatusForbidden Forbidden

swagger:response backupsRestoreStatusForbidden
*/
type BackupsRestoreStatusForbidden struct {

	/*
	  In: Body
	*/
	Payload *models.ErrorResponse `json:"body,omitempty"`
}

// NewBackupsRestoreStatusForbidden creates BackupsRestoreStatusForbidden with default headers values
func NewBackupsRestoreStatusForbidden() *BackupsRestoreStatusForbidden {

	return &BackupsRestoreStatusForbidden{}
}

// WithPayload adds the payload to the backups restore status forbidden response
func (o *BackupsRestoreStatusForbidden) WithPayload(payload *models.ErrorResponse) *BackupsRestoreStatusForbidden {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the backups restore status forbidden response
func (o *BackupsRestoreStatusForbidden) SetPayload(payload *models.ErrorResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *BackupsRestoreStatusForbidden) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(403)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// BackupsRestoreStatusNotFoundCode is the HTTP code returned for type BackupsRestoreStatusNotFound
const BackupsRestoreStatusNotFoundCode int = 404

/*
BackupsRestoreStatusNotFound Not Found - Backup does not exist

swagger:response backupsRestoreStatusNotFound
*/
type BackupsRestoreStatusNotFound struct {

	/*
	  In: Body
	*/
	Payload *models.ErrorResponse `json:"body,omitempty"`
}

// NewBackupsRestoreStatusNotFound creates BackupsRestoreStatusNotFound with default headers values
func NewBackupsRestoreStatusNotFound() *BackupsRestoreStatusNotFound {

	return &BackupsRestoreStatusNotFound{}
}

// WithPayload adds the payload to the backups restore status not found response
func (o *BackupsRestoreStatusNotFound) WithPayload(payload *models.ErrorResponse) *BackupsRestoreStatusNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the backups restore status not found response
func (o *BackupsRestoreStatusNotFound) SetPayload(payload *models.ErrorResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *BackupsRestoreStatusNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// BackupsRestoreStatusInternalServerErrorCode is the HTTP code returned for type BackupsRestoreStatusInternalServerError
const BackupsRestoreStatusInternalServerErrorCode int = 500

/*
BackupsRestoreStatusInternalServerError An error has occurred while trying to fulfill the request. Most likely the ErrorResponse will contain more information about the error.

swagger:response backupsRestoreStatusInternalServerError
*/
type BackupsRestoreStatusInternalServerError struct {

	/*
	  In: Body
	*/
	Payload *models.ErrorResponse `json:"body,omitempty"`
}

// NewBackupsRestoreStatusInternalServerError creates BackupsRestoreStatusInternalServerError with default headers values
func NewBackupsRestoreStatusInternalServerError() *BackupsRestoreStatusInternalServerError {

	return &BackupsRestoreStatusInternalServerError{}
}

// WithPayload adds the payload to the backups restore status internal server error response
func (o *BackupsRestoreStatusInternalServerError) WithPayload(payload *models.ErrorResponse) *BackupsRestoreStatusInternalServerError {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the backups restore status internal server error response
func (o *BackupsRestoreStatusInternalServerError) SetPayload(payload *models.ErrorResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *BackupsRestoreStatusInternalServerError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
