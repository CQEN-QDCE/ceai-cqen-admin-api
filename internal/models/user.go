package models

type User struct {
	Disabled     *bool  `json:"disabled,omitempty"`
	Email        string `json:"email"`
	Firstname    string `json:"firstname"`
	Infrarole    string `json:"infrarole"`
	Lastname     string `json:"lastname"`
	Organisation string `json:"organisation"`
}

type UserUpdate struct {
	Disabled     *bool   `json:"disabled,omitempty"`
	Firstname    *string `json:"firstname,omitempty"`
	Infrarole    *string `json:"infrarole,omitempty"`
	Lastname     *string `json:"lastname,omitempty"`
	Organisation *string `json:"organisation,omitempty"`
}

type UserWithLabs struct {
	// Embedded struct due to allOf(#/components/schemas/User)
	User `yaml:",inline"`
	// Embedded fields due to inline allOf schema
	Laboratories *[]LaboratoryRole `json:"laboratories,omitempty"`
}

type LaboratoryRole struct {
	Laboratory string `json:"laboratory"`
	Role       string `json:"role"`
}
