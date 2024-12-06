package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"

	"github.com/algrvvv/pomodoro/internal/app"
	"github.com/algrvvv/pomodoro/internal/config"
	"github.com/algrvvv/pomodoro/internal/database"
	"github.com/algrvvv/pomodoro/internal/notify"
	repos "github.com/algrvvv/pomodoro/internal/repositories/postgres"
	"github.com/algrvvv/pomodoro/internal/server"
)

var (
	onlyWebServer    bool
	startSessionType string
)

func main() {
	flag.BoolVar(&onlyWebServer, "only-web", false, "only web server")
	flag.StringVar(&startSessionType, "type", "work", "first session type (can be work or break)")
	flag.Parse()

	if err := config.Parse("config.yml"); err != nil {
		fmt.Println("failed to load config")
		os.Exit(1)
	}
	config.Config.Pomodoro.SessionGoalMinutes = config.Config.Pomodoro.WorkMinutes * config.Config.Pomodoro.SessionsGoal

	if err := database.Connect(); err != nil {
		fmt.Println("database connection failed: ", err)
	}
	defer database.Close()

	wg := sync.WaitGroup{}
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)

	sessionRepo := repos.NewPostgresRepo()
	serv := server.NewServer(sessionRepo)

	if !onlyWebServer {
		wg.Add(1)
		notifier := notify.GetTerminalNotifier()

		go func() {
			if err := app.Start(ctx, notifier, sessionRepo, startSessionType, &wg); err != nil {
				fmt.Println("app error: ", err)
				os.Exit(0)
			}
		}()
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := serv.ListenAndServe(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				fmt.Println("web server error: ", err)
				os.Exit(0)
			}
		}
	}()

	<-ctx.Done()
	fmt.Println("got Interrupt. exiting")

	if err := serv.Shutdown(context.Background()); err != nil {
		fmt.Println("server shutdown error: ", err)
	}

	cancel()
	wg.Wait()

	fmt.Println("success")
}
