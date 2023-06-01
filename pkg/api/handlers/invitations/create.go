package invitations

import (
	"b2match_api/pkg/database"
	"b2match_api/pkg/database/models"
	"b2match_api/utils"
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func Create(ctx *gin.Context) {
	var meeting models.EventMeeting
	user := utils.GetUser(ctx)

	err := database.DB.QueryRow("SELECT id, end_time, event_id FROM event_meetings WHERE id = ?", ctx.Param("meeting_id")).Scan(&meeting.ID, &meeting.EndTime, &meeting.Event.ID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": "meeting not found",
			})
			return
		}
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "error occurred while doing the query",
		})
		return
	}

	if time.Now().After(meeting.EndTime) {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"message": "meeting is already finished",
		})
		return
	}

	var req struct {
		UserID int `json:"user_id"`
	}

	err = ctx.BindJSON(&req)

	if !utils.ContentChecker(ctx, utils.AcceptedTypes["JSON"]) || req.UserID == 0 {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "body is invalid",
		})
		return
	}

	var count int
	if err := database.DB.QueryRow("SELECT COUNT(*) FROM meeting_invitations where user_id = ? AND meeting_id = ?", req.UserID, meeting.ID).Scan(&count); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "error occurred while doing the query",
		})
		return
	}

	if count > 0 {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"message": "user is already invited to this meeting",
		})
		return
	}

	if err := database.DB.QueryRow("SELECT COUNT(*) FROM event_participants where user_id = ? AND event_id = ?", req.UserID, meeting.Event.ID).Scan(&count); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "error occurred while doing the query",
		})
		return
	}

	if count < 1 {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"message": "can not add user to the event because user is not taking part in it",
		})
		return
	}

	if _, err := database.DB.Exec("INSERT INTO meeting_invitations (user_id, invited_by_id, meeting_id) VALUES (?, ?, ?)", req.UserID, user.ID, meeting.ID); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "invitation created",
	})
}
