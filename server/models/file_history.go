package models

import (
	"time"

	"github.com/ptcoffee/image-server/database"
)

// HistoryInfoRes is history of a file response to client
type HistoryInfoRes struct {
	ID             uint                `json:"id"`
	At             time.Time           `json:"at"`
	Fullname       string              `json:"fullname"`
	BackupFullname *string             `json:"backupFullname"`
	ActionType     database.FileAction `json:"actionType"`
}

// NewHistoryInfo returns HistoryInfo from FileHistory
func NewHistoryInfo(fileHistory database.FileHistory) HistoryInfoRes {
	return HistoryInfoRes{
		ID:             fileHistory.ID,
		At:             fileHistory.UpdatedAt,
		Fullname:       fileHistory.Fullname,
		BackupFullname: fileHistory.BackupFullname,
		ActionType:     fileHistory.ActionType,
	}
}

// NewHistoryInfos return array of HistoryInfo from FileHistory
func NewHistoryInfos(fileHistories []database.FileHistory) []HistoryInfoRes {
	historyInfos := make([]HistoryInfoRes, len(fileHistories))
	for index := range fileHistories {
		historyInfos[index] = NewHistoryInfo(fileHistories[index])
	}
	return historyInfos
}
