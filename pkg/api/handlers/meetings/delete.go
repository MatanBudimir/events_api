package meetings

import (
	"b2match_api/pkg/database"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Delete(ctx *gin.Context) {
	var count int

	err := database.DB.QueryRow("SELECT COUNT(*) FROM event_meetings WHERE id = ?", ctx.Param("meeting_id")).Scan(&count)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "error occurred while doing the query",
		})
		return
	}

	if count < 1 {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message": "meeting not found",
		})
		return
	}

	// delete all invitations
	if _, err = database.DB.Exec("DELETE FROM meeting_invitations WHERE meeting_id = ?", ctx.Param("meeting_id")); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "error occurred while doing the query",
		})
		return
	}

	// delete the meeting
	if _, err = database.DB.Exec("DELETE FROM event_meetings WHERE id = ?", ctx.Param("meeting_id")); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "error occurred while doing the query",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "meeting deleted",
	})
}
