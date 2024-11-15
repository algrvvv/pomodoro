package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/algrvvv/pomodoro/internal/repositories"
)

func GetData(repo repositories.SessionRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessions, err := repo.GetAll()
		if err != nil {
			fmt.Println("failed to get all sessions: ", err)
			http.Error(w, "failed to get all sessions", http.StatusInternalServerError)
			return
		}

		// var datesAchivedGoal []string

		data := map[string]interface{}{
			"status":   true,
			"sessions": sessions,
			"calendar": []string{"2024-11-05", "2024-11-10", "2024-11-15"},
			"chart": map[string]int{
				"2024-11-05": 3,
				"2024-11-10": 5,
				"2024-11-15": 8,
			},
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
