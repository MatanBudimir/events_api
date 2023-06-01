package models

type MeetingInvitation struct {
	ID      int          `json:"id"`
	User    User         `json:"user,omitempty"`
	Meeting EventMeeting `json:"-"`
	Status  byte         `json:"status"`
}
