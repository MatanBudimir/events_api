package invitations

import (
	"b2match_api/pkg/database"
	"b2match_api/pkg/database/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ListInvitations(ctx *gin.Context) {
	var invitations []models.MeetingInvitation

	rows, err := database.DB.Query("SELECT meeting_invitations.id, meeting_invitations.status, meeting_invitations.user_id, users.first_name, users.last_name, o.id, o.name FROM meeting_invitations INNER JOIN users ON meeting_invitations.user_id = users.id INNER JOIN organizations o on users.organization_id = o.id WHERE meeting_invitations.meeting_id = ?", ctx.Param("meeting_id"))

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "error occurred while doing the query",
		})
		return
	}

	for rows.Next() {
		var invitation models.MeetingInvitation
		if err := rows.Scan(&invitation.ID, &invitation.Status, &invitation.User.ID, &invitation.User.FirstName, &invitation.User.LastName, &invitation.User.Organization.ID, &invitation.User.Organization.Name); err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "error occurred while doing the query",
			})
			return
		}

		invitations = append(invitations, invitation)
	}

	ctx.JSON(http.StatusOK, invitations)
}
