package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"

	"github.com/algrvvv/pomodoro/internal/app"
	"github.com/algrvvv/pomodoro/internal/config"
	"github.com/algrvvv/pomodoro/internal/notify"
)

func main() {
	if err := config.Parse("config.yml"); err != nil {
		fmt.Println("failed to load config")
		os.Exit(1)
	}

	wg := sync.WaitGroup{}
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)

	wg.Add(1)
	notifier := notify.GetMacosNotifier()

	if err := app.Start(ctx, notifier, &wg); err != nil {
		fmt.Println("app error: ", err)
		os.Exit(0)
	}

	<-ctx.Done()
	fmt.Println("got Interrupt. exiting")

	cancel()
	wg.Wait()

	fmt.Println("success")
}
