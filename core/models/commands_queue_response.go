package models




import (
	strfmt "github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"

	"github.com/go-openapi/errors"
)

// CommandsQueueResponse commands queue response
// swagger:model CommandsQueueResponse
type CommandsQueueResponse struct {

	// Commands to be executed.
	Commands []*Command `json:"commands"`
}

// Validate validates this commands queue response
func (m *CommandsQueueResponse) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateCommands(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *CommandsQueueResponse) validateCommands(formats strfmt.Registry) error {

	if swag.IsZero(m.Commands) { // not required
		return nil
	}

	for i := 0; i < len(m.Commands); i++ {

		if swag.IsZero(m.Commands[i]) { // not required
			continue
		}

		if m.Commands[i] != nil {

			if err := m.Commands[i].Validate(formats); err != nil {
				return err
			}
		}

	}

	return nil
}
