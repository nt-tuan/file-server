package models

import (
	"time"

	"github.com/ptcoffee/image-server/database"
)

type HistoryInfo struct {
	ID             uint                `json:"id"`
	At             time.Time           `json:"at"`
	Fullname       string              `json:"fullname"`
	BackupFullname *string             `json:"backupFullname"`
	ActionType     database.FileAction `json:"actionType"`
}

// NewHistoryInfo return response body of history file
func NewHistoryInfo(fileHistory database.FileHistory) HistoryInfo {
	return HistoryInfo{
		ID:             fileHistory.ID,
		At:             fileHistory.UpdatedAt,
		Fullname:       fileHistory.Fullname,
		BackupFullname: fileHistory.BackupFullname,
		ActionType:     fileHistory.ActionType,
	}
}

func NewHistoryInfos(fileHistories []database.FileHistory) []HistoryInfo {
	historyInfos := make([]HistoryInfo, len(fileHistories))
	for index := range fileHistories {
		historyInfos[index] = NewHistoryInfo(fileHistories[index])
	}
	return historyInfos
}
