package handler

import (
	"testing"

	"agenthub/backend/internal/model"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// setupTestDB creates an in-memory mock database for testing.
func setupTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	t.Helper()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open gorm: %v", err)
	}
	return gormDB, mock
}

func TestCascadeDeleteBySessionIDs_EmptySlice(t *testing.T) {
	db, mock := setupTestDB(t)
	// With empty slice, no queries should be executed
	cascadeDeleteBySessionIDs(db, []string{})
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestCascadeDeleteBySessionIDs_DeletesAllTables(t *testing.T) {
	db, mock := setupTestDB(t)

	sessionIDs := []string{"sess-1", "sess-2"}

	// Expect DELETE for Message
	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM `messages`").WithArgs(sessionIDs[0], sessionIDs[1]).WillReturnResult(sqlmock.NewResult(0, 5))
	mock.ExpectCommit()

	// Expect DELETE for SessionAgent
	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM `session_agents`").WithArgs(sessionIDs[0], sessionIDs[1]).WillReturnResult(sqlmock.NewResult(0, 2))
	mock.ExpectCommit()

	// Expect DELETE for DiffSnapshot
	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM `diff_snapshots`").WithArgs(sessionIDs[0], sessionIDs[1]).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	cascadeDeleteBySessionIDs(db, sessionIDs)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestCascadeDeleteByTaskID_DeletesAllTables(t *testing.T) {
	db, mock := setupTestDB(t)

	taskID := "task-123"
	sessionIDs := []string{"sess-1"}

	// Expect SELECT session_ids
	mock.ExpectQuery("SELECT `session_id` FROM `sessions`").WithArgs(taskID).WillReturnRows(
		sqlmock.NewRows([]string{"session_id"}).AddRow("sess-1"),
	)

	// cascadeDeleteBySessionIDs queries
	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM `messages`").WithArgs(sessionIDs[0]).WillReturnResult(sqlmock.NewResult(0, 3))
	mock.ExpectCommit()
	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM `session_agents`").WithArgs(sessionIDs[0]).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()
	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM `diff_snapshots`").WithArgs(sessionIDs[0]).WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	// Expect DELETE for Session
	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM `sessions`").WithArgs(taskID).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	// Expect DELETE for Announcement
	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM `announcements`").WithArgs(taskID).WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	cascadeDeleteByTaskID(db, taskID)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

// Verify the cascade helper covers all expected models.
func TestCascadeModels(t *testing.T) {
	// This test just verifies the model types exist and are referenced correctly.
	// The actual delete logic is tested via sqlmock above.
	var _ model.Message
	var _ model.SessionAgent
	var _ model.DiffSnapshot
	var _ model.Session
	var _ model.Announcement
}
