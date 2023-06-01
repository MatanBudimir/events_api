package events

import (
	"b2match_api/pkg/database"
	"b2match_api/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AddUserToEvent(ctx *gin.Context) {
	var count int
	// checking if event exists
	err := database.DB.QueryRow("SELECT COUNT(*) FROM events WHERE id = ?", ctx.Param("event_id")).Scan(&count)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "error occurred while doing the query",
		})
		return
	}

	if count < 1 {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message": "event not found",
		})
		return
	}

	user := utils.GetUser(ctx)

	// check if user is already taking part in this event
	err = database.DB.QueryRow("SELECT COUNT(*) FROM event_participants WHERE event_id = ? AND user_id = ?", ctx.Param("event_id"), user.ID).Scan(&count)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "error occurred while doing the query",
		})
		return
	}

	if count > 0 {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"message": "user is already taking part in the event",
		})
		return
	}

	// add user to event

	if _, err := database.DB.Exec("INSERT INTO event_participants (event_id, user_id) VALUES (?, ?)", ctx.Param("event_id"), user.ID); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "error occurred while doing the query",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "joined the event",
	})
}
