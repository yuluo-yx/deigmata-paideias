package report

import (
	"embed"
	"encoding/json"
	"html/template"
	"os"
	"time"

	"github.com/deigmata-paideias/github-contributor/internal/ghapi"
)

//go:embed templates/report.html
var templateFS embed.FS

type Data struct {
	GeneratedAt string                     `json:"generated_at"`
	Stats       ghapi.OrgContributionStats `json:"stats"`
}

func WriteHTML(path string, data Data) error {
	if data.GeneratedAt == "" {
		data.GeneratedAt = time.Now().Format("2006-01-02 15:04:05 MST")
	}

	tmpl, err := template.ParseFS(templateFS, "templates/report.html")
	if err != nil {
		return err
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return tmpl.Execute(f, data)
}

func WriteJSON(path string, data Data) error {
	if data.GeneratedAt == "" {
		data.GeneratedAt = time.Now().Format("2006-01-02 15:04:05 MST")
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(data)
}
