package store

import (
	"time"

	"cms/models"

	"github.com/google/uuid"
)

var auditLogStore = NewStore("audit_logs.json")

type AuditLogList struct {
	AuditLogs []models.AuditLog `json:"audit_logs"`
}

func GetAllAuditLogs() ([]models.AuditLog, error) {
	var list AuditLogList
	if err := auditLogStore.ReadAll(&list); err != nil {
		return nil, err
	}
	return list.AuditLogs, nil
}

func CreateAuditLog(log *models.AuditLog) error {
	logs, err := GetAllAuditLogs()
	if err != nil {
		return err
	}

	log.ID = uuid.New().String()
	log.CreatedAt = time.Now()

	logs = append(logs, *log)

	if len(logs) > 1000 {
		logs = logs[len(logs)-1000:]
	}

	return auditLogStore.WriteAll(AuditLogList{AuditLogs: logs})
}

func GetAuditLogsByUser(userID string) ([]models.AuditLog, error) {
	logs, err := GetAllAuditLogs()
	if err != nil {
		return nil, err
	}

	var result []models.AuditLog
	for _, l := range logs {
		if l.UserID == userID {
			result = append(result, l)
		}
	}

	return result, nil
}

func GetAuditLogsByAction(action string) ([]models.AuditLog, error) {
	logs, err := GetAllAuditLogs()
	if err != nil {
		return nil, err
	}

	var result []models.AuditLog
	for _, l := range logs {
		if l.Action == action {
			result = append(result, l)
		}
	}

	return result, nil
}
