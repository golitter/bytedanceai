package impl

import (
	"log/slog"
	"unicode/utf8"

	"agenthub/backend/internal/dao"
	"agenthub/backend/internal/model"
)

const maxGroupChatMsgLen = 2000

func fetchGroupChatWindow(messageDao dao.MessageDao, taskID, sessionID string) []map[string]interface{} {
	lastMessage, err := messageDao.FindLatestCompletedAgentMessage(taskID, sessionID)
	if err != nil {
		slog.Warn("group chat window query failed to load latest agent message, degrading to empty", "task_id", taskID, "session_id", sessionID, "error", err)
		return []map[string]interface{}{}
	}

	messages, err := messageDao.ListGroupChatWindowMessages(taskID, sessionID, lastMessage)
	if err != nil {
		slog.Warn("group chat window query failed, degrading to empty", "task_id", taskID, "session_id", sessionID, "error", err)
		return []map[string]interface{}{}
	}

	result := make([]map[string]interface{}, 0, len(messages))
	seen := make(map[string]bool, len(messages))
	for _, message := range messages {
		content := truncateWindowMessage(message)
		dedupeKey := message.AgentName + "\x00" + message.AgentType + "\x00" + content
		if seen[dedupeKey] {
			continue
		}
		seen[dedupeKey] = true
		result = append(result, map[string]interface{}{
			"role":       message.Role,
			"agent_name": message.AgentName,
			"content":    content,
		})
	}
	return result
}

func truncateWindowMessage(message model.Message) string {
	content := message.Content
	if utf8.RuneCountInString(content) <= maxGroupChatMsgLen {
		return content
	}
	runes := []rune(content)
	return string(runes[:maxGroupChatMsgLen]) + "\n...[截断]"
}
