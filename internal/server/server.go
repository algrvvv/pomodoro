package server

import (
	"net/http"

	"github.com/algrvvv/pomodoro/internal/config"
	"github.com/algrvvv/pomodoro/internal/repositories"
	"github.com/algrvvv/pomodoro/internal/server/handlers"
)

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Добавляем CORS-заголовки
		w.Header().Set("Access-Control-Allow-Origin", "*") // Разрешаем все источники
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Если это preflight-запрос, отправляем пустой ответ
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Передаем управление следующему обработчику
		next.ServeHTTP(w, r)
	})
}

func NewServer(r repositories.SessionRepository) *http.Server {
	s := http.NewServeMux()

	// роуты
	s.HandleFunc("GET /api/v1/data", handlers.GetData(r))
	s.HandleFunc("/ws", handlers.GetPassedSessionTime)
	// интеграции
	s.HandleFunc("GET /api/v1/integrations/wakatime", handlers.WakatimeIntegration)

	// статический роут для клиентских файлов
	s.Handle("/", http.FileServer(http.Dir("./static")))

	return &http.Server{
		Addr:    ":" + config.Config.App.Port,
		Handler: corsMiddleware(s),
	}
}
