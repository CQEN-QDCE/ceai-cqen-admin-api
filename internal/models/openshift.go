package models

import "time"

type OpenshiftProject struct {
	Description string `json:"description"`
	Displayname string `json:"displayname"`
	Id          string `json:"id"`
}

type OpenshiftProjectWithLab struct {
	// Embedded struct due to allOf(#/components/schemas/OpenshiftProject)
	*OpenshiftProject `yaml:",inline"`
	// Embedded fields due to inline allOf schema
	IdLab string `json:"idLab,omitempty"`
}

type OpenshiftProjectWithMeta struct {
	// Embedded struct due to allOf(#/components/schemas/OpenshiftProjectWithLab)
	*OpenshiftProjectWithLab `yaml:",inline"`
	// Embedded fields due to inline allOf schema
	CreationDate *time.Time `json:"creationDate,omitempty"`
	Requester    *string    `json:"requester,omitempty"`
}
