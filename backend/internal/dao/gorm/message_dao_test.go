package gormdao

import (
	"strings"
	"testing"

	"agenthub/backend/internal/model"

	"gorm.io/gorm"
)

func TestApplyGroupMessageVisibilityUsesDirectTurnReplyPredicate(t *testing.T) {
	db, mock := setupTestDB(t)
	query := db.Session(&gorm.Session{DryRun: true}).
		Model(&model.Message{}).
		Where("task_id = ?", "task-123")

	stmt := applyGroupMessageVisibility(query, "task-123", "orch-session").Find(&[]model.Message{}).Statement
	sql := stmt.SQL.String()

	for _, want := range []string{
		"session_id = ?",
		"group_id <> ?",
		"SELECT 1 FROM messages user_msg",
		"SELECT 1 FROM messages agent_msg",
		"agent_msg.id > user_msg.id",
		"agent_msg.id < messages.id",
	} {
		if !strings.Contains(sql, want) {
			t.Fatalf("group visibility SQL = %q, want substring %q", sql, want)
		}
	}
	if strings.Contains(sql, "DISTINCT session_id") {
		t.Fatalf("group visibility should not expose all messages from directly-addressed sessions: %q", sql)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unfulfilled expectations: %v", err)
	}
}
