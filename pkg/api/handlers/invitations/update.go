package invitations

import (
	"b2match_api/pkg/database"
	"b2match_api/utils"
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func UpdateInvitation(ctx *gin.Context) {
	var count int
	var startTime, endTime time.Time
	user := utils.GetUser(ctx)

	err := database.DB.QueryRow("SELECT COUNT(*), em.start_time, em.end_time FROM meeting_invitations INNER JOIN event_meetings em on meeting_invitations.meeting_id = em.id WHERE meeting_invitations.id = ? AND meeting_id = ? AND meeting_invitations.user_id = ?", ctx.Param("invitation_id"), ctx.Param("meeting_id"), user.ID).Scan(&count, &startTime, &endTime)

	if !errors.Is(sql.ErrNoRows, err) && err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "error occurred while doing the query",
		})
		return
	}

	if count < 1 {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message": "invitation not found",
		})
		return
	}

	var req struct {
		Status byte `json:"status"`
	}

	err = ctx.BindJSON(&req)

	if !utils.ContentChecker(ctx, utils.AcceptedTypes["JSON"]) || err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "body is invalid",
		})
		return
	}

	if req.Status != 1 && req.Status != 2 {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "invalid status provided. Valid statuses are: 1, 2",
		})
		return
	}

	if req.Status == 2 {
		if err := database.DB.QueryRow("SELECT COUNT(*) FROM meeting_invitations invitations JOIN event_meetings meetings on invitations.meeting_id = meetings.id WHERE (invitations.user_id = ? AND status = '2' AND invitations.meeting_id != ?) AND (meetings.start_time <= ? AND meetings.end_time >= ? OR meetings.start_time <= ? AND meetings.end_time >= ? OR meetings.start_time >= ? AND meetings.end_time <= ?)", user.ID, ctx.Param("meeting_id"), startTime, startTime, endTime, endTime, startTime, endTime).Scan(&count); err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "error occurred while doing the query",
			})
			return
		}

		if count > 0 {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"message": "you already have a meeting at that time",
			})
			return
		}
	}

	if _, err := database.DB.Exec("UPDATE meeting_invitations SET status = ? WHERE id = ?", string(req.Status+'0'), ctx.Param("invitation_id")); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "invitation updated",
	})
}
