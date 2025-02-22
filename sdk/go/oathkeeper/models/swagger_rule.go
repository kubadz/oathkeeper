// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"strconv"

	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
)

// SwaggerRule swaggerRule is a single rule that will get checked on every HTTP request.
// swagger:model swaggerRule
type SwaggerRule struct {

	// Authenticators is a list of authentication handlers that will try and authenticate the provided credentials.
	// Authenticators are checked iteratively from index 0 to n and if the first authenticator to return a positive
	// result will be the one used.
	//
	// If you want the rule to first check a specific authenticator  before "falling back" to others, have that authenticator
	// as the first item in the array.
	Authenticators []*SwaggerRuleHandler `json:"authenticators"`

	// Description is a human readable description of this rule.
	Description string `json:"description,omitempty"`

	// ID is the unique id of the rule. It can be at most 190 characters long, but the layout of the ID is up to you.
	// You will need this ID later on to update or delete the rule.
	ID string `json:"id,omitempty"`

	// authorizer
	Authorizer *SwaggerRuleHandler `json:"authorizer,omitempty"`

	// match
	Match *SwaggerRuleMatch `json:"match,omitempty"`

	// mutator
	Mutator *SwaggerRuleHandler `json:"mutator,omitempty"`

	// upstream
	Upstream *Upstream `json:"upstream,omitempty"`
}

// Validate validates this swagger rule
func (m *SwaggerRule) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateAuthenticators(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateAuthorizer(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateMatch(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateMutator(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateUpstream(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *SwaggerRule) validateAuthenticators(formats strfmt.Registry) error {

	if swag.IsZero(m.Authenticators) { // not required
		return nil
	}

	for i := 0; i < len(m.Authenticators); i++ {
		if swag.IsZero(m.Authenticators[i]) { // not required
			continue
		}

		if m.Authenticators[i] != nil {
			if err := m.Authenticators[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("authenticators" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (m *SwaggerRule) validateAuthorizer(formats strfmt.Registry) error {

	if swag.IsZero(m.Authorizer) { // not required
		return nil
	}

	if m.Authorizer != nil {
		if err := m.Authorizer.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("authorizer")
			}
			return err
		}
	}

	return nil
}

func (m *SwaggerRule) validateMatch(formats strfmt.Registry) error {

	if swag.IsZero(m.Match) { // not required
		return nil
	}

	if m.Match != nil {
		if err := m.Match.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("match")
			}
			return err
		}
	}

	return nil
}

func (m *SwaggerRule) validateMutator(formats strfmt.Registry) error {

	if swag.IsZero(m.Mutator) { // not required
		return nil
	}

	if m.Mutator != nil {
		if err := m.Mutator.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("mutator")
			}
			return err
		}
	}

	return nil
}

func (m *SwaggerRule) validateUpstream(formats strfmt.Registry) error {

	if swag.IsZero(m.Upstream) { // not required
		return nil
	}

	if m.Upstream != nil {
		if err := m.Upstream.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("upstream")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *SwaggerRule) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *SwaggerRule) UnmarshalBinary(b []byte) error {
	var res SwaggerRule
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
