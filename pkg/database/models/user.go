package models

type User struct {
	ID             int          `json:"id"`
	FirstName      string       `json:"first_name"`
	LastName       string       `json:"last_name"`
	Email          string       `json:"-"`
	OrganizationID int          `json:"-"`
	Organization   Organization `json:"organization"`
}
