package app

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/algrvvv/pomodoro/internal/config"
	"github.com/algrvvv/pomodoro/internal/notify"
)

func saveTime(t *string) {
	fmt.Println("save working time: ", *t)
}

func Start(ctx context.Context, n notify.Notifier, wg *sync.WaitGroup) error {
	defer wg.Done()

	var timeDiff string
	defer saveTime(&timeDiff)
	defer n.Notify("Title", "Some message")

	start := time.Now()
	workDuration := time.Duration(config.Config.WorkMinutes) * time.Second
	c, cancel := context.WithTimeout(ctx, workDuration)
	defer cancel()

	fmt.Println("start work")
	go func(c context.Context) {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				dif := time.Since(start)

				hours := int(dif.Hours())
				minutes := int(dif.Minutes()) % 60
				seconds := int(dif.Seconds()) % 60
				timeDiff = fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)

				fmt.Printf("\rTime worked: %s", timeDiff)
			case <-c.Done():
				fmt.Println("\n\ngot interrupt. returned")
				return
			}
		}
	}(c)

	select {
	case <-time.After(workDuration):
		fmt.Println("\n\nend work")
	case <-ctx.Done():
		return errors.New("exit from app with interrupt")
	}

	return errors.New("exit from app start")
}
