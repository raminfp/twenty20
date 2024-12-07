package main

import (
	"errors"
	"fmt"
	"os/exec"
	"time"
	"twenty20/sound"
)

// Notify displays a desktop notification.
// title: The title of the notification.
// message: The body of the notification.
// iconPath: Path to the notification icon (can be empty).
func Notify(title, message, iconPath string) error {
	// Check for required command availability
	cmd := "notify-send"
	if _, err := exec.LookPath(cmd); err != nil {
		return errors.New("notify-send command not found. Ensure it is installed and available in PATH")
	}

	// Prepare the arguments
	args := []string{title, message}
	if iconPath != "" {
		args = append(args, "-i", iconPath)
	}

	// Execute the command
	err := exec.Command(cmd, args...).Run()
	if err != nil {
		return errors.New("failed to send notification: " + err.Error())
	}

	return nil
}

func main() {

	fmt.Println("Playing beep sound...")
	for {

		err := Notify("20-20-20 Rule", "Time to rest! Look at something far away for 20 seconds.", "")
		if err != nil {
			fmt.Printf("Notification error: %v\n", err)
		} else {
			fmt.Println("Notification sent!")
		}
		// Beep at 440 Hz for 200 milliseconds
		err = sound.Sound(440.0, 200)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		time.Sleep(15 * time.Minute) // Every 20 minutes
	}

}
