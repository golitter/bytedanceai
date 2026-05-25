package model

import "time"

type DiffSnapshot struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	SnapshotID  string    `gorm:"uniqueIndex;size:36" json:"snapshot_id"`
	SessionID   string    `gorm:"index;size:128" json:"session_id"`
	DiffContent string    `gorm:"type:longtext" json:"diff_content"`
	Status      string    `gorm:"size:16;default:pending" json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
