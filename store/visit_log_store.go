package store

import (
	"time"

	"cms/models"

	"github.com/google/uuid"
)

var visitLogStore = NewStore("visit_logs.json")

type VisitLogList struct {
	VisitLogs []models.VisitLog `json:"visit_logs"`
}

func GetAllVisitLogs() ([]models.VisitLog, error) {
	var list VisitLogList
	if err := visitLogStore.ReadAll(&list); err != nil {
		return nil, err
	}
	return list.VisitLogs, nil
}

func CreateVisitLog(log *models.VisitLog) error {
	logs, err := GetAllVisitLogs()
	if err != nil {
		return err
	}

	log.ID = uuid.New().String()
	log.CreatedAt = time.Now()

	logs = append(logs, *log)

	cutoff := time.Now().AddDate(0, 0, -30)
	var recentLogs []models.VisitLog
	for _, l := range logs {
		if l.CreatedAt.After(cutoff) {
			recentLogs = append(recentLogs, l)
		}
	}

	return visitLogStore.WriteAll(VisitLogList{VisitLogs: recentLogs})
}

func GetTodayUV() (int, error) {
	logs, err := GetAllVisitLogs()
	if err != nil {
		return 0, err
	}

	today := time.Now().Format("2006-01-02")
	ipSet := make(map[string]bool)

	for _, l := range logs {
		if l.CreatedAt.Format("2006-01-02") == today {
			ipSet[l.IP] = true
		}
	}

	return len(ipSet), nil
}

func GetTodayPV() (int, error) {
	logs, err := GetAllVisitLogs()
	if err != nil {
		return 0, err
	}

	today := time.Now().Format("2006-01-02")
	count := 0

	for _, l := range logs {
		if l.CreatedAt.Format("2006-01-02") == today {
			count++
		}
	}

	return count, nil
}

func GetVisitStatsByDay(days int) (map[string]map[string]int, error) {
	logs, err := GetAllVisitLogs()
	if err != nil {
		return nil, err
	}

	result := make(map[string]map[string]int)
	now := time.Now()

	for i := 0; i < days; i++ {
		date := now.AddDate(0, 0, -i)
		dateStr := date.Format("2006-01-02")
		result[dateStr] = map[string]int{
			"pv": 0,
			"uv": 0,
		}
	}

	uvIPs := make(map[string]map[string]bool)
	for _, l := range logs {
		dateStr := l.CreatedAt.Format("2006-01-02")
		if _, exists := result[dateStr]; exists {
			result[dateStr]["pv"]++
			if uvIPs[dateStr] == nil {
				uvIPs[dateStr] = make(map[string]bool)
			}
			uvIPs[dateStr][l.IP] = true
		}
	}

	for dateStr, ips := range uvIPs {
		if _, exists := result[dateStr]; exists {
			result[dateStr]["uv"] = len(ips)
		}
	}

	return result, nil
}
