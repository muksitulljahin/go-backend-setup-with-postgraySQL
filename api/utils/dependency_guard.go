package utils

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-backend-setup-with-postgray/response"
	"gorm.io/gorm"
)

type DependencyRule struct {
	Table       string
	ForeignKey  string
	EntityLabel string
}

func EnsureNoDependencies(db *gorm.DB, id string, rules ...DependencyRule) error {
	trimmedID := strings.TrimSpace(id)
	if trimmedID == "" {
		return response.NewAppError(http.StatusBadRequest, "ID is required")
	}
	for _, rule := range rules {
		var count int64
		err := db.Table(rule.Table).
			Where(fmt.Sprintf("%s = ?", rule.ForeignKey), trimmedID).
			Count(&count).Error
		if err != nil {
			return response.NewAppError(http.StatusInternalServerError, "Failed to check record dependencies", err)
		}
		if count > 0 {
			msg := fmt.Sprintf("Cannot delete because %d %s currently depend on it. Please reassign or delete them first.", count, rule.EntityLabel)
			return response.NewAppError(http.StatusBadRequest, msg)
		}
	}
	return nil
}
