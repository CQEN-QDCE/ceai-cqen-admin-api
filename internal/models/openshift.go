package models

import "time"

type OpenshiftProject struct {
	Id          string `json:"id"`
	Description string `json:"description"`
	Displayname string `json:"displayname"`
}

type OpenshiftProjectWithLab struct {
	// Embedded struct due to allOf(#/components/schemas/OpenshiftProject)
	OpenshiftProject `yaml:",inline"`
	// Embedded fields due to inline allOf schema
	IdLab string `json:"idLab"`
}

type OpenshiftProjectWithMeta struct {
	// Embedded struct due to allOf(#/components/schemas/OpenshiftProjectWithLab)
	OpenshiftProjectWithLab `yaml:",inline"`
	// Embedded fields due to inline allOf schema
	Requester    *string    `json:"requester,omitempty"`
	CreationDate *time.Time `json:"creationDate,omitempty"`
}

type OpenshiftProjectUpdate struct {
	Description *string `json:"description,omitempty"`
	Displayname *string `json:"displayname,omitempty"`
	IdLab       *string `json:"idLab,omitempty"`
}
