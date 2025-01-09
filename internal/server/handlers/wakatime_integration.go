package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/algrvvv/pomodoro/internal/config"
)

type integration = config.Integration

type resp struct {
	response string
	error    error
}

type wakatimeResponse struct {
	Total struct {
		Seconds float64 `json:"seconds"`
		Digital string  `json:"digital"`
		Text    string  `json:"text"`
	} `json:"cumulative_total"`
}

// type weekResponse struct {
//   Data struct {} `json:"data"`
// }

func WakatimeIntegration(w http.ResponseWriter, r *http.Request) {
	var i integration

	for _, integration := range config.Config.Intergations {
		if strings.ToLower(integration.Name) == "wakatime" && integration.Enabled {
			i = integration
			break
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if i.Name == "" {
		err := json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  http.StatusBadRequest,
			"message": "integration not found or not enabled",
		})
		if err != nil {
			fmt.Println("failed to encode message: ", err)
		}

		return
	}

	weekStats := make(chan resp)
	todayStats := make(chan resp)

	wg := &sync.WaitGroup{}
	context, cancel := context.WithTimeout(context.Background(), i.Timeout)
	defer cancel()

	start := time.Now()
	wg.Add(4) // 2 - делают запросы; 2 - их обрабатывают
	go getSummaries(context, todayStats, "today", i, wg)
	go getSummaries(context, weekStats, "last_7_days", i, wg)

	weekResp := make(map[string]interface{})
	todayResp := make(map[string]interface{})

	go func() {
		defer wg.Done()
		for c := range weekStats {
			if c.error != nil {
				weekResp = map[string]interface{}{
					"status": false,
					"data":   c.error.Error(),
				}
			} else {
				weekResp = map[string]interface{}{
					"status": true,
					"data":   c.response,
				}
			}
		}
	}()

	go func() {
		defer wg.Done()
		for c := range todayStats {
			if c.error != nil {
				todayResp = map[string]interface{}{
					"status": false,
					"data":   c.error.Error(),
				}
			} else {
				todayResp = map[string]interface{}{
					"status": true,
					"data":   c.response,
				}
			}
		}
	}()

	wg.Wait()

	err := json.NewEncoder(w).Encode(map[string]interface{}{
		"week":  weekResp,
		"today": todayResp,
		"time":  time.Since(start).String(),
	})
	if err != nil {
		fmt.Println("failed to encode message: ", err)
	}
}

func getSummaries(ctx context.Context, ch chan resp, sr string, i integration, wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(ch)

	const baseUrl = "https://api.wakatime.com/api/v1/users/current/summaries?range="
	url := fmt.Sprintf("%s%s", baseUrl, sr)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		ch <- resp{error: err}
		return
	}

	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", i.ApiKey))
	client := http.Client{}
	r, err := client.Do(req)
	if err != nil {
		ch <- resp{error: err}
		return
	}
	defer r.Body.Close()

	var answer wakatimeResponse

	err = json.NewDecoder(r.Body).Decode(&answer)
	if err != nil {
		ch <- resp{error: err}
		return
	}

	ch <- resp{response: answer.Total.Text}
}
