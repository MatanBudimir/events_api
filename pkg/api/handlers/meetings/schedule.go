package meetings

import (
	"b2match_api/pkg/database"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Schedule(ctx *gin.Context) {
	var count int

	if err := database.DB.QueryRow("SELECT COUNT(*) FROM event_meetings WHERE id = ? AND scheduled = FALSE", ctx.Param("meeting_id")).Scan(&count); err != nil {
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

	if err := database.DB.QueryRow("SELECT COUNT(*) FROM meeting_invitations WHERE meeting_id = ? AND status != '2'", ctx.Param("meeting_id")).Scan(&count); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "error occurred while doing the query",
		})
		return
	}

	if count > 0 {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message": "can not schedule the meeting because some people did not accept the invite",
		})
		return
	}

	if _, err := database.DB.Exec("UPDATE event_meetings SET scheduled = TRUE WHERE id = ?", ctx.Param("meeting_id")); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "error occurred while doing the query",
		})
		return
	}

	ctx.AbortWithStatusJSON(http.StatusOK, gin.H{
		"message": "meeting scheduled",
	})
	return
}
