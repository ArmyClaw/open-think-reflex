package models

import (
	"time"

	"github.com/google/uuid"
)

// ThoughtSession groups a sequence of thought nodes.
type ThoughtSession struct {
	ID        string
	Title     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// ThoughtNode is a single thought entry in a session.
type ThoughtNode struct {
	ID        string
	SessionID string
	ParentID  string
	Text      string
	CreatedAt time.Time
}

func NewThoughtSession(title string) *ThoughtSession {
	now := time.Now()
	return &ThoughtSession{
		ID:        uuid.New().String(),
		Title:     title,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func NewThoughtNode(sessionID, parentID, text string) *ThoughtNode {
	now := time.Now()
	return &ThoughtNode{
		ID:        uuid.New().String(),
		SessionID: sessionID,
		ParentID:  parentID,
		Text:      text,
		CreatedAt: now,
	}
}
