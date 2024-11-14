package notify

import (
	"fmt"
	"os/exec"
)

type MacosNotifier struct{}

func GetMacosNotifier() MacosNotifier {
	return MacosNotifier{}
}

func (m MacosNotifier) Notify(title, message string) error {
	cmd := exec.Command(
		"osascript",
		"-e",
		fmt.Sprintf(
			`display notification "%s" with title "%s"`,
			message,
			title,
		),
	)
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
