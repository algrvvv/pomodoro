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
	"github.com/algrvvv/pomodoro/internal/utils"
)

var (
	sessionCount = 0
	sessionType  = types.WorkSession
	notifier     notify.Notifier
)

func saveTime(t *string) {
	if sessionType == types.WorkSession {
		sessionCount++
	}

	if sessionCount == config.Config.SessionsGoal {
		// TODO: потом отслеживать в бд, чтобы не повторяться
		fmt.Print(notify.GoalMessage)
		notifier.Notify(
			"You have reached your daily goal!",
			fmt.Sprintf("You have already held %d working sessions today", sessionCount),
		)
	}
	// fmt.Printf("\n\n[%d] save working time: %v for %v\n\n", sessionCount, *t, sessionType)
}

func Start(ctx context.Context, n notify.Notifier, wg *sync.WaitGroup) error {
	defer wg.Done()
	scanner := bufio.NewScanner(os.Stdin)
	notifier = n

	for {
		if sessionType == types.WorkSession {
			utils.ClearTerminal()
			fmt.Print(notify.BackToWork)

			workDuration := time.Duration(config.Config.WorkMinutes) * time.Second
			if err := startTimer(ctx, workDuration); err != nil {
				return err
			}

			sessionType = types.BreakSession
		} else if sessionType == types.BreakSession {
			utils.ClearTerminal()
			fmt.Print(notify.BreakTime)

			var breakDuration time.Duration

			if sessionCount%config.Config.BreakAfterSessions == 0 && sessionCount != 0 {
				breakDuration = time.Duration(config.Config.LongBreakMinutes) * time.Second
			} else {
				breakDuration = time.Duration(config.Config.ShortBreakMinutes) * time.Second
			}

			if err := startTimer(ctx, breakDuration); err != nil {
				return err
			}

			sessionType = types.WorkSession
		}

		fmt.Println("\n\npress Enter to continue or write q to quit...")
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

func startTimer(ctx context.Context, duration time.Duration) error {
	var timeDiff string
	defer saveTime(&timeDiff)
	defer func(diff *string) {
		var title, message string
		if sessionType == types.WorkSession {
			// TODO: проверить чтобы при раннем отключении помидора все равно выводилось сообщение
			// со временем проведенным за работой
			title = "It's time to take a break!"
			message = fmt.Sprintf("Your working session #%d lasted %s", sessionCount+1, *diff)
		} else if sessionType == types.BreakSession {
			title = "It's time to get back to work"
			message = "Start a new work session right now!"
		}

		if err := notifier.Notify(title, message); err != nil {
			fmt.Println("failed to notify: ", err)
		}
	}(&timeDiff)

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

				fmt.Printf("\rTime passed: %s", timeDiff)
			case <-c.Done():
				fmt.Print("\n\n\n\n")
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
