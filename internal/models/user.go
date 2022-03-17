package models

type User struct {
	Email        string `json:"email"`
	Firstname    string `json:"firstname"`
	Lastname     string `json:"lastname"`
	Infrarole    string `json:"infrarole"`
	Organisation string `json:"organisation"`
	Disabled     *bool  `json:"disabled,omitempty"`
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

type AuthenticatedUser struct {
	Server   *string `json:"server"`
	Username *string `json:"username"`
	Roles    *string `json:"roles"`
}
