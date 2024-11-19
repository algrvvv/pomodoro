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

func GetData(repo repositories.SessionRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessions, err := repo.GetAll()
		if err != nil {
			fmt.Println("failed to get all sessions: ", err)
			http.Error(w, "failed to get all sessions", http.StatusInternalServerError)
			return
		}

		sCount := make(map[string]int, len(sessions))
		sDuration := make(map[string]int, len(sessions))
		totalMinutes := 0
		for _, s := range sessions {
			if s.SessionType != types.WorkSession {
				continue
			}

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
			"calendar":     datesAchivedGoal,
			"chart":        sCount,
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
