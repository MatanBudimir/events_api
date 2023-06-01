package meetings

import (
	"b2match_api/pkg/database"
	"b2match_api/pkg/database/models"
	"encoding/json"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestListMeetings(t *testing.T) {
	r := gin.Default()

	r.GET("/api/events/:event_id/meetings", ListMeetings)

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create SQL mock: %s", err)
	}
	defer db.Close()

	database.DB = db

	mock.ExpectQuery("SELECT id, name, description, start_time, end_time, scheduled FROM event_meetings WHERE scheduled = TRUE AND event_id = \\? AND end_time >= CURRENT_TIMESTAMP()").
		WithArgs("1").
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description", "start_time", "end_time", "scheduled"}).AddRow(1, "Testing event", "I'm a test", time.Now().Add(time.Hour*4), time.Now().Add(time.Hour*5), true).
			AddRow(2, "Another Testing event", "I'm a test 1", time.Now(), time.Now().Add(time.Hour), true).AddRow(3, "Yes Another Testing event", "I'm a test 2", time.Now().Add(time.Hour*2), time.Now().Add(time.Hour*3), true))

	req, err := http.NewRequest("GET", "/api/events/1/meetings", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %s", err)
	}

	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status code %d but got %d", http.StatusOK, rec.Code)
	}

	var response []models.EventMeeting
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse response body: %s", err)
	}

	for i, meeting := range response {
		if meeting.ID != i+1 {
			t.Errorf("Expected event with ID - %d, got - %d", i+1, meeting.ID)
		}
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}
