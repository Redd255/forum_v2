package server

import (
	"fmt"
	"net/http"

	"github.com/gofrs/uuid"
)

func convertime(seconds int64) string {
	timeUnits := []struct {
		duration int64
		singular string
	}{
		{31557600, "Year"},
		{2629744, "Month"},
		{86400, "Day"},
		{3600, "Hour"},
		{60, "Minute"},
	}

	for _, unit := range timeUnits {
		if seconds >= unit.duration {
			count := seconds / unit.duration
			return pluralize(count, unit.singular)
		}
	}

	if seconds == 0 {
		return "Now"
	}
	return fmt.Sprintf("%d s", seconds)
}

func pluralize(count int64, unit string) string {
	if count == 1 {
		return fmt.Sprintf("%d %s", count, unit)
	}
	return fmt.Sprintf("%d %ss", count, unit)
}

func GenerateUUID() (string, error) {
	newUUID, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	return newUUID.String(), nil
}
func errorPage(w http.ResponseWriter, message, templateName string) {
	data := map[string]string{
		"ErrorMessage": message,
	}
	err := tmpl.ExecuteTemplate(w, templateName, data)
	if err != nil {
		http.Error(w, "Failed to execute template: "+err.Error(), http.StatusInternalServerError)
	}
}
