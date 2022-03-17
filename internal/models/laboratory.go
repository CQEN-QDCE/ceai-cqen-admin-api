package models

type Laboratory struct {
	Id          string  `json:"id"`
	Displayname string  `json:"displayname"`
	Description string  `json:"description"`
	Type        string  `json:"type"`
	Gitrepo     *string `json:"gitrepo,omitempty"`
}

type LaboratoryRole struct {
	Laboratory `yaml:",inline"`
	Role       string `json:"role"`
}

type LaboratoryUpdate struct {
	Description *string `json:"description,omitempty"`
	Displayname *string `json:"displayname,omitempty"`
	Type        *string `json:"type,omitempty"`
	Gitrepo     *string `json:"gitrepo,omitempty"`
}

type LaboratoryWithUsers struct {
	// Embedded struct due to allOf(#/components/schemas/Laboratory)
	Laboratory `yaml:",inline"`
	// Embedded fields due to inline allOf schema
	Users *[]string `json:"users,omitempty"`
}

type LaboratoryWithResources struct {
	// Embedded struct due to allOf(#/components/schemas/LaboratoryWithUsers)
	LaboratoryWithUsers `yaml:",inline"`
	// Embedded fields due to inline allOf schema
	AWSAccounts       *[]AWSAccount       `json:"AWSAccounts,omitempty"`
	Openshiftprojects *[]OpenshiftProject `json:"openshiftprojects,omitempty"`
}
