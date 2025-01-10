package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/algrvvv/pomodoro/internal/config"
)

// NOTE: есть возможность для оптимизации и уменьшении кода - вынести основную логику создания запроса в отдельную
// функцию, но мне так лень уже это делать. может как то в будущем это все таки сделаю... Помечу как туду

// TODO:  можно потом вынести общую часть запросов в отдульный метод

// получение тестов за сегодня
// https://api.monkeytype.com/users/currentTestActivity - ссылка на получение последней активности
// дока - https://api.monkeytype.com/docs/#tag/users/operation/users.getCurrentTestActivity

// получение общего колва тестов и времени за тестами
// https://api.monkeytype.com/users/stats - ссылка для получения данных
// дока - https://api.monkeytype.com/docs/#tag/users/operation/users.getStats

type todayResp struct {
	Message string `json:"message"`
	Data    struct {
		TestsByDays []int `json:"testsByDays"`
		LastDay     int64 `json:"lastDay"`
	} `json:"data"`

	todayTests int
	error      error
}

type totalResp struct {
	Message string `json:"message"`
	Data    struct {
		StartedTests int     `json:"startedTests"`
		TimeTyping   float64 `json:"timeTyping"`
	} `json:"data"`

	time  string
	error error
}

func MonkeyTypeIntegration(w http.ResponseWriter, r *http.Request) {
	var i integration

	for _, integration := range config.Config.Intergations {
		if strings.ToLower(integration.Name) == "monkeytype" && integration.Enabled {
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

	todayData := make(chan todayResp)
	totalData := make(chan totalResp)

	wg := &sync.WaitGroup{}
	context, cancel := context.WithTimeout(context.Background(), i.Timeout)
	defer cancel()

	wg.Add(2)
	start := time.Now()
	go getTodayTests(context, todayData, i, wg)
	go getTotalInfo(context, totalData, i, wg)

	var getDataError string

	today := <-todayData
	total := <-totalData

	if today.error != nil {
		getDataError = "ошибка получения сегодняшних данных: " + today.error.Error()
	}
	if total.error != nil {
		getDataError += "\nошибка получения данных профиля: " + total.error.Error()
	}

	err := json.NewEncoder(w).Encode(map[string]interface{}{
		"todayTests":   today.todayTests,
		"startedTests": total.Data.StartedTests,
		"timeTyping":   total.time,
		"lastDay":      today.Data.LastDay,
		"time":         time.Since(start).String(),
		"error":        getDataError,
	})
	if err != nil {
		fmt.Println("failed to encode message: ", err)
	}
}

func getBodyFromReq(ctx context.Context, url string, key string) (io.ReadCloser, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("ApeKey %s", key))
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}

func getTodayTests(ctx context.Context, ch chan todayResp, i integration, wg *sync.WaitGroup) {
	defer wg.Done()

	const url = "https://api.monkeytype.com/users/currentTestActivity"
	body, err := getBodyFromReq(ctx, url, i.ApiKey)
	if err != nil {
		ch <- todayResp{error: err}
		return
	}
	defer body.Close()

	var answer todayResp
	err = json.NewDecoder(body).Decode(&answer)
	if err != nil {
		ch <- todayResp{error: err}
		return
	}

	l := len(answer.Data.TestsByDays)
	answer.todayTests = answer.Data.TestsByDays[l-1]

	ch <- answer
}

func getTotalInfo(ctx context.Context, ch chan totalResp, i integration, wg *sync.WaitGroup) {
	defer wg.Done()

	const url = "https://api.monkeytype.com/users/stats"
	body, err := getBodyFromReq(ctx, url, i.ApiKey)
	if err != nil {
		ch <- totalResp{error: err}
	}
	defer body.Close()

	var answer totalResp

	err = json.NewDecoder(body).Decode(&answer)
	if err != nil {
		ch <- totalResp{error: err}
	}

	dur := time.Duration(answer.Data.TimeTyping) * time.Second
	t := time.Time{}.Add(dur)
	answer.time = t.Format(time.TimeOnly)

	ch <- answer
}
