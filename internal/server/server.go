package server

import (
	"net/http"

	"github.com/algrvvv/pomodoro/internal/config"
	"github.com/algrvvv/pomodoro/internal/repositories"
	"github.com/algrvvv/pomodoro/internal/server/handlers"
)

func NewServer(r repositories.SessionRepository) *http.Server {
	s := http.NewServeMux()

	s.HandleFunc("GET /api/v1/data", handlers.GetData(r))
	s.HandleFunc("/ws", handlers.GetPassedSessionTime)

	// статический роут для клиентских файлов
	s.Handle("/", http.FileServer(http.Dir("./static")))

	return &http.Server{
		Addr:    ":" + config.Config.App.Port,
		Handler: s,
	}
}
