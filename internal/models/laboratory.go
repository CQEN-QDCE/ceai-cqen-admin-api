package models

type AWSAccount struct {
	Email *string `json:"email,omitempty"`
	Name  *string `json:"name,omitempty"`
	Id    *string `json:"number,omitempty"`
}

// Laboratory defines model for Laboratory.
type Laboratory struct {
	Id          string  `json:"id"`
	Description string  `json:"description"`
	Displayname string  `json:"displayname"`
	Type        string  `json:"type"`
	Gitrepo     *string `json:"gitrepo,omitempty"`
}

type LaboratoryUpdate struct {
	Description *string `json:"description,omitempty"`
	Displayname *string `json:"displayname,omitempty"`
	Type        *string `json:"type,omitempty"`
	Gitrepo     *string `json:"gitrepo,omitempty"`
}

// LaboratoryWithResources defines model for LaboratoryWithResources.
type LaboratoryWithResources struct {
	// Embedded struct due to allOf(#/components/schemas/LaboratoryWithUsers)
	*LaboratoryWithUsers `yaml:",inline"`
	// Embedded fields due to inline allOf schema
	AWSAccounts       *[]AWSAccount       `json:"AWSAccounts,omitempty"`
	Openshiftprojects *[]OpenshiftProject `json:"openshiftprojects,omitempty"`
}

// LaboratoryWithUsers defines model for LaboratoryWithUsers.
type LaboratoryWithUsers struct {
	// Embedded struct due to allOf(#/components/schemas/Laboratory)
	*Laboratory `yaml:",inline"`
	// Embedded fields due to inline allOf schema
	Users *[]string `json:"users,omitempty"`
}
