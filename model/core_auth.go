package model

import validation "github.com/go-ozzo/ozzo-validation"

type LoginRequest struct {
	// ClientID string `json:"client_id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

// Validate the LoginRequest
func (req LoginRequest) Validate() error {
	return validation.ValidateStruct(&req,
		// validation.Field(&req.ClientID, validation.Required, validation.Length(8, 64)),
		validation.Field(&req.Name, validation.Required, validation.Length(3, 16)),
		validation.Field(&req.Color, validation.Required),
	)
}
