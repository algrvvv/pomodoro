package app

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/algrvvv/pomodoro/internal/config"
	"github.com/algrvvv/pomodoro/internal/notify"
	"github.com/algrvvv/pomodoro/internal/types"
)

var (
	sessionCount = 0
	sessionType  = types.WorkSession
)

func saveTime(t *string) {
	if sessionType == types.WorkSession {
		sessionCount++
	}

	if sessionCount == config.Config.SessionsGoal {
		fmt.Println("\n\nGOAL\n\n")
	}
	fmt.Printf("[%d] save working time: %v for %v", sessionCount, *t, sessionType)
}

func Start(ctx context.Context, n notify.Notifier, wg *sync.WaitGroup) error {
	defer wg.Done()
	scanner := bufio.NewScanner(os.Stdin)

	for {
		if sessionType == types.WorkSession {
			workDuration := time.Duration(config.Config.WorkMinutes) * time.Second
			fmt.Println("start work")
			if err := startTimer(ctx, workDuration, n); err != nil {
				return err
			}
			fmt.Println("\n\nend work")

			sessionType = types.BreakSession
		} else if sessionType == types.BreakSession {
			var breakDuration time.Duration

			if sessionCount%config.Config.BreakAfterSessions == 0 && sessionCount != 0 {
				breakDuration = time.Duration(config.Config.LongBreakMinutes) * time.Second
			} else {
				breakDuration = time.Duration(config.Config.ShortBreakMinutes) * time.Second
			}

			fmt.Println("start break")
			if err := startTimer(ctx, breakDuration, n); err != nil {
				return err
			}
			fmt.Println("\n\nend break")

			sessionType = types.WorkSession
		}

		fmt.Println("\n\npress Enter to continue...")
		if scanner.Scan() {
			line := scanner.Text()
			if line == "q" || line == "quit" {
				return errors.New("quit by user")
			}
		} else {
			if scanner.Err() != nil {
				return scanner.Err()
			}
			return errors.New("failed to get data from scanner")
		}
	}
}

func startTimer(ctx context.Context, duration time.Duration, n notify.Notifier) error {
	var timeDiff string
	defer saveTime(&timeDiff)
	defer n.Notify("Title", "Some message")

	start := time.Now()
	c, cancel := context.WithTimeout(ctx, duration)
	defer cancel()

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
	case <-ctx.Done():
		return errors.New("exit from app with interrupt")
	case <-time.After(duration):
	}

	return nil
}
