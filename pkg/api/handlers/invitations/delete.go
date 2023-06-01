package invitations

import (
	"b2match_api/pkg/database"
	"b2match_api/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Delete(ctx *gin.Context) {
	var count int

	err := database.DB.QueryRow("SELECT COUNT(*) FROM meeting_invitations WHERE user_id = ? AND id = ? AND meeting_id = ?", utils.GetUser(ctx).ID, ctx.Param("invitation_id"), ctx.Param("meeting_id")).Scan(&count)

	if err != nil {
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

	if _, err := database.DB.Exec("DELETE FROM meeting_invitations WHERE id = ?", ctx.Param("invitation_id")); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "error occurred while doing the query",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "invitation deleted",
	})
}
