package clients

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"
	"tracktora-backend/internal/models"
)

// AdzunaResponse matches the JSON structure from the Adzuna API documentation
type AdzunaResponse struct {
	Results []struct {
		Title       string  `json:"title"`
		Description string  `json:"description"`
		RedirectURL string  `json:"redirect_url"`
		Created     string  `json:"created"`
		SalaryMin   float64 `json:"salary_min"`
		SalaryMax   float64 `json:"salary_max"`
		Company     struct {
			DisplayName string `json:"display_name"`
		} `json:"company"`
		Location struct {
			DisplayName string `json:"display_name"`
		} `json:"location"`
	} `json:"results"`
}

// FetchLiveJobs handles fetching and paginating results from Adzuna India
func FetchLiveJobs(keyword, location string, page int, salaryMin int) ([]models.ExternalJob, error) {
	appID := os.Getenv("ADZUNA_APP_ID")
	appKey := os.Getenv("ADZUNA_APP_KEY")

	if keyword == "" {
		keyword = "internship"
	}

	if page < 1 {
		page = 1
	}

	// Build base URL with pagination and keyword
	baseUrl := fmt.Sprintf("https://api.adzuna.com/v1/api/jobs/in/search/%d?app_id=%s&app_key=%s&results_per_page=15&what=%s",
		page, appID, appKey, url.QueryEscape(keyword))

	// Append location if exists
	if location != "" {
		baseUrl = fmt.Sprintf("%s&where=%s", baseUrl, url.QueryEscape(location))
	}

	// Append salary filter if greater than 0
	if salaryMin > 0 {
		baseUrl = fmt.Sprintf("%s&salary_min=%d", baseUrl, salaryMin)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(baseUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to reach Adzuna: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("adzuna api error: status %d", resp.StatusCode)
	}

	var adzunaResp AdzunaResponse
	if err := json.NewDecoder(resp.Body).Decode(&adzunaResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	var jobs []models.ExternalJob
	for _, item := range adzunaResp.Results {
		jobs = append(jobs, models.ExternalJob{
			Title:       item.Title,
			Company:     item.Company.DisplayName,
			Location:    item.Location.DisplayName,
			Description: item.Description,
			ApplyURL:    item.RedirectURL,
			Source:      "Adzuna India",
			PublishedAt: item.Created,
			SalaryMin:   item.SalaryMin,
			SalaryMax:   item.SalaryMax,
		})
	}

	return jobs, nil
}
