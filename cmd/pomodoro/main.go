package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"

	"github.com/algrvvv/pomodoro/internal/app"
	"github.com/algrvvv/pomodoro/internal/config"
	"github.com/algrvvv/pomodoro/internal/database"
	"github.com/algrvvv/pomodoro/internal/notify"
	repos "github.com/algrvvv/pomodoro/internal/repositories/postgres"
)

func main() {
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

	wg.Add(1)
	notifier := notify.GetTerminalNotifier()
	sessionRepo := repos.NewPostgresRepo()

	if err := app.Start(ctx, notifier, sessionRepo, &wg); err != nil {
		fmt.Println("app error: ", err)
		os.Exit(0)
	}

	<-ctx.Done()
	fmt.Println("got Interrupt. exiting")

	cancel()
	wg.Wait()

	fmt.Println("success")
}
