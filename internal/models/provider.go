package models

type Provider struct {
	Name        *string `json:"name,omitempty"`
	Type        *string `json:"type,omitempty"`
	Description *string `json:"description,omitempty"`
	ClientUIID  *string `json:"clientUIID,omitempty"`
	Healthy     *bool   `json:"healthy,omitempty"`
}
