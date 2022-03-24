package models

type AWSAccount struct {
	Id    *string `json:"id,omitempty"`
	Name  *string `json:"name,omitempty"`
	Email *string `json:"email,omitempty"`
}
