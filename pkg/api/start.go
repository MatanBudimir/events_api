package api

import (
	"b2match_api/pkg/api/handlers/events"
	"b2match_api/pkg/api/handlers/invitations"
	"b2match_api/pkg/api/handlers/meetings"
	"b2match_api/utils"
	"github.com/gin-gonic/gin"
	"log"
)

func StartWebServer() {
	gin.SetMode(mode())
	r := gin.Default()
	// add authentication middleware
	r.Use(Auth())

	// add user to specific event middleware
	r.POST("/api/events/:event_id/users", events.AddUserToEvent)

	// delete the meeting
	r.DELETE("/api/meetings/:meeting_id", meetings.Delete)

	// create a meeting
	r.POST("/api/meetings", meetings.CreateMeeting)

	// get a list of scheduled meetings
	r.GET("/api/events/:event_id/meetings", meetings.ListMeetings)

	// get a list of all meeting invitations
	r.GET("/api/meetings/:meeting_id/invitations", invitations.ListInvitations)

	// create a meeting invitation
	r.POST("/api/meetings/:meeting_id/invitations", invitations.Create)

	// delete invitation
	r.DELETE("/api/meetings/:meeting_id/invitations/:invitation_id", invitations.Delete)

	// update the meeting invitation
	r.PUT("/api/meetings/:meeting_id/invitations/:invitation_id", invitations.UpdateInvitation)

	// schedule the meeting
	r.POST("/api/meetings/:meeting_id/schedule", meetings.Schedule)

	if err := r.Run(utils.GetEnv("HTTP_PORT", ":8000")); err != nil {
		log.Fatalln(err)
	}
}

func mode() string {
	if mode := utils.GetEnv("DEBUG_MODE", "true"); mode == "false" {
		return gin.ReleaseMode
	}

	return gin.DebugMode
}
