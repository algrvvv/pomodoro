package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/algrvvv/pomodoro/internal/config"
	"github.com/algrvvv/pomodoro/internal/repositories"
	"github.com/algrvvv/pomodoro/internal/types"
)

var displayZero bool

func GetData(repo repositories.SessionRepository) http.HandlerFunc {
	displayZero = config.Config.App.DisplayZeroDays

	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		var endDate string
		startDate := r.URL.Query().Get("date")
		if startDate == "" {
			startDate = time.Now().Format(time.DateOnly)
			endDate = time.Now().AddDate(0, 1, 0).Format(time.DateOnly)
		} else {
			e, err := time.Parse(time.DateOnly, startDate)
			if err != nil {
				fmt.Println("failed to parse start date")
				return
			}

			endDate = e.AddDate(0, 1, 0).Format(time.DateOnly)
		}

		sessions, err := repo.GetAll(startDate, endDate)
		if err != nil {
			fmt.Println("failed to get all sessions: ", err)
			http.Error(w, "failed to get all sessions", http.StatusInternalServerError)
			return
		}

		// данные для подсказок ['day' => 1: minutes; 2: minutes]
		dataForTooltip := make(map[int]map[int]int)
		// колво сессий для графика ['date' => 'count of session']
		sCount := make(map[string]int, len(sessions))
		// колво минут для графика ['date' => 'minutes of session']
		sDuration := make(map[string]int, len(sessions))
		totalMinutes := 0
		var lastDate time.Time
		for _, s := range sessions {
			day := s.CreatedAt.Day()
			if _, ok := dataForTooltip[day]; !ok {
				dataForTooltip[day] = make(map[int]int)
			}

			if s.SessionType != types.WorkSession {
				if _, ok := dataForTooltip[day][types.BreakSession]; !ok {
					dataForTooltip[day][types.BreakSession] = s.Minutes
				} else {
					dataForTooltip[day][types.BreakSession] += s.Minutes
				}

				continue
			}

			if _, ok := dataForTooltip[day][types.WorkSession]; !ok {
				dataForTooltip[day][types.WorkSession] = s.Minutes
			} else {
				dataForTooltip[day][types.WorkSession] += s.Minutes
			}

			// проверка на пропещенные дни, котрые мы должны забить нулями
			if d := lastDate.Day() - s.CreatedAt.Day(); displayZero && d > 1 {
				for i := lastDate.Add(-24 * time.Hour); i.Day() != s.CreatedAt.Day(); i = i.Add(-24 * time.Hour) {
					iFormat := i.Format(time.DateOnly)
					sCount[iFormat] = 0
					sDuration[iFormat] = 0
				}
			}
			lastDate = s.CreatedAt

			t := s.CreatedAt.Format(time.DateOnly)
			if _, ok := sCount[t]; !ok {
				sCount[t] = 1
				sDuration[t] = s.Minutes
			} else {
				sCount[t] += 1
				sDuration[t] += s.Minutes
			}
			totalMinutes += s.Minutes
		}

		if displayZero {
			// если есть дни с начала месяца, которых нет в общем списке сессий - заполняем их нулями
			date := lastDate
			for d := lastDate.Day() - 1; d > 0; d-- {
				date = date.Add(-24 * time.Hour)
				sCount[date.Format(time.DateOnly)] = 0
				sDuration[date.Format(time.DateOnly)] = 0
			}
		}

		datesAchivedGoal := []string{}
		for d, m := range sDuration {
			if m >= config.Config.Pomodoro.SessionGoalMinutes {
				datesAchivedGoal = append(datesAchivedGoal, d)
			}
		}

		todaySessions, err := repo.GetTodaySessions()
		if err != nil {
			fmt.Println("failed to get today sessions: ", err)
			http.Error(w, "failed to get sessions", http.StatusInternalServerError)
			return
		}

		data := map[string]interface{}{
			"status":       true,
			"minutesGoal":  config.Config.Pomodoro.SessionGoalMinutes,
			"totalMinutes": totalMinutes,
			"sessions":     todaySessions,
			"tooltips":     dataForTooltip,
			"calendar":     datesAchivedGoal,
			"chartCount":   sCount,
			"chartMinutes": sDuration,
			"time":         time.Since(start).String(),
		}

		// jsonData, err := json.Marshal(data)
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(data)
		if err != nil {
			fmt.Println("failed to marshal data: ", err)
			http.Error(w, "failed to prepare data", http.StatusInternalServerError)
			return
		}

		// if _, err := w.Write(jsonData); err != nil {
		// 	fmt.Println("failed to write data to response")
		// }
	}
}
