package notify

import "os/exec"

type TerminalNotifier struct {
	iconPath string
}

func GetTerminalNotifier() TerminalNotifier {
	return TerminalNotifier{iconPath: "./assets/icon.png"}
}

func (t TerminalNotifier) Notify(title, message string) error {
	args := []string{
		"-title", title, "-message", message, "-timeout", "10",
		"-sound", "Glass",
		// "-appIcon", t.iconPath,
	}

	cmd := exec.Command("terminal-notifier", args...)
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
