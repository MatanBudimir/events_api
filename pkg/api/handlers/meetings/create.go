package meetings

import (
	"b2match_api/pkg/database"
	"b2match_api/pkg/database/models"
	"b2match_api/utils"
	"context"
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type request struct {
	models.EventMeeting
	Invitations []struct {
		UserID int `json:"user_id"`
	}
	EventID int `json:"event_id"`
}

func CreateMeeting(ctx *gin.Context) {
	var count int
	user := utils.GetUser(ctx)

	var req request

	if !utils.ContentChecker(ctx, utils.AcceptedTypes["JSON"]) {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "invalid content type",
		})
		return
	}

	if err := ctx.BindJSON(&req); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	if err := req.validateRequest(); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	var event models.Event
	if err := database.DB.QueryRow("SELECT id, start_time, end_time FROM events where id = ?", req.EventID).Scan(&event.ID, &event.StartTime, &event.EndTime); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": "event not found",
			})
			return
		}

		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "error occurred while doing the query",
		})
		return
	}

	// check is meeting taking place between event start and end time
	if event.StartTime.After(req.StartTime) || event.EndTime.Before(req.EndTime) {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"message": "can not create a meeting that is not taking place during the event",
		})
		return
	}

	if err := database.DB.QueryRow("SELECT COUNT(*) FROM meeting_invitations invitations JOIN event_meetings meetings on invitations.meeting_id = meetings.id WHERE (invitations.user_id = ? AND status = '2') AND (meetings.start_time <= ? AND meetings.end_time >= ? OR meetings.start_time <= ? AND meetings.end_time >= ? OR meetings.start_time >= ? AND meetings.end_time <= ?)", user.ID, req.StartTime.Format("2006-01-02 15:04:05"), req.StartTime.Format("2006-01-02 15:04:05"), req.EndTime.Format("2006-01-02 15:04:05"), req.EndTime.Format("2006-01-02 15:04:05"), req.StartTime.Format("2006-01-02 15:04:05"), req.EndTime.Format("2006-01-02 15:04:05")).Scan(&count); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "error occurred while doing the query",
		})
		return
	}

	if count > 0 {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"message": "can not create meeting because user who tried to create a meeting has a meeting at the same time",
		})
		return
	}

	tx, err := database.DB.BeginTx(context.Background(), nil)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "error occurred while doing the query",
		})
		return
	}

	defer tx.Rollback()

	res, err := tx.Exec("INSERT INTO event_meetings (name, description, user_id, event_id, start_time, end_time) VALUES (?, ?, ?, ?, ?, ?)", req.Name, req.Description, user.ID, event.ID, req.StartTime.Format("2006-01-02 15:04:05"), req.EndTime.Format("2006-01-02 15:04:05"))
	meetingID, idErr := res.LastInsertId()

	if err != nil || idErr != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "error occurred while doing the query",
		})
		return
	}

	var failed []int

	if _, err := tx.Exec("INSERT INTO meeting_invitations (user_id, invited_by_id, meeting_id, status) VALUES (?, ?, ?, '2')", user.ID, user.ID, meetingID); err != nil {
		failed = append(failed, user.ID)
	}

	for _, data := range req.Invitations {
		if err := processInvitation(tx, data.UserID, user.ID, meetingID); err != nil {
			failed = append(failed, data.UserID)
		}
	}

	if err := tx.Commit(); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "error occurred while doing the query",
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message":    "meeting created",
		"failed_ids": failed,
	})
}

func (req request) validateRequest() error {
	if req.EventID < 1 {
		return errors.New("invalid event ID provided")
	} else if len(req.Name) < 1 {
		return errors.New("invalid name provided")
	} else if len(req.Description) < 1 {
		return errors.New("invalid description provided")
	} else if time.Now().After(req.StartTime) {
		return errors.New("invalid start time provided")
	} else if req.StartTime.After(req.EndTime) {
		return errors.New("invalid end time provided")
	}

	return nil
}

func processInvitation(tx *sql.Tx, userId, invitedByID int, meetingID int64) error {
	if userId < 1 {
		return nil
	}

	_, err := tx.Exec("INSERT INTO meeting_invitations (user_id, invited_by_id, meeting_id) VALUES (?, ?, ?)", userId, invitedByID, meetingID)

	return err
}
