package notify

import "fmt"

type MacosNotifier struct{}

func GetMacosNotifier() MacosNotifier {
	return MacosNotifier{}
}

func (m MacosNotifier) Notify(title, message string) error {
	fmt.Println("MACOS NOTIFICATION: ", title, message)
	return nil
}
