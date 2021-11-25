package models

type AWSAccount struct {
	Email *string `json:"email,omitempty"`
	Id    *string `json:"id,omitempty"`
	Name  *string `json:"name,omitempty"`
}
