package models

type EventParticipant struct {
	ID    int   `json:"id"`
	Event Event `json:"event,omitempty"`
	User  User  `json:"user,omitempty"`
}
