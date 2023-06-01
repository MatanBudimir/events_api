package meetings

import (
	"b2match_api/pkg/database"
	"b2match_api/pkg/database/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ListMeetings(ctx *gin.Context) {
	var meetings []models.EventMeeting

	rows, err := database.DB.Query("SELECT id, name, description, start_time, end_time, scheduled FROM event_meetings WHERE scheduled = TRUE AND event_id = ? AND end_time >= CURRENT_TIMESTAMP()", ctx.Param("event_id"))

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	for rows.Next() {
		var meeting models.EventMeeting
		if err := rows.Scan(&meeting.ID, &meeting.Name, &meeting.Description, &meeting.StartTime, &meeting.EndTime, &meeting.Scheduled); err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "error occurred while doing the query",
			})
			return
		}

		meetings = append(meetings, meeting)
	}

	ctx.JSON(http.StatusOK, meetings)
}
