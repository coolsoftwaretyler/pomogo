package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gen2brain/beeep"
)

// Main function that runs when pomogo is executed
func main() {
	// Get the amount of time
	timeInput := getTheTime()
	timeInput = strings.TrimSuffix(timeInput, "\n")
	pomodoroTime, _ := strconv.Atoi(timeInput)
	// Get the task
	task := getTheTask()
	task = strings.TrimSuffix(task, "\n")
	// Get the start time and calculate the end time
	startTime := time.Now()
	endTime := startTime.Add(time.Minute * time.Duration(pomodoroTime))
	// Log the start time, the time, the task
	startMessage := fmt.Sprintf("%v, %v, %v\n", startTime, timeInput, task)
	setIntention(startMessage)
	// Set the timer
	for range time.Tick(1 * time.Second) {
		timeRemaining := getTimeRemaining(endTime)

		if timeRemaining.t <= 0 {
			break
		}
		// Update the status
		fmt.Printf("Minutes: %d Seconds: %d\n", timeRemaining.m, timeRemaining.s)
	}
	// When timer comes back, make a sound
	err := beeep.Notify("pomogo", "Pomodoro over!", "assets/information.png")
	if err != nil {
		panic(err)
	}
	// When the timer completes, ask how it went
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("How did that go?")
	result, _ := reader.ReadString('\n')
	result = strings.TrimSuffix(result, "\n")
	setIntention(result)
}

func getTheTime() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("How much time?")
	text, _ := reader.ReadString('\n')
	return text
}

func getTheTask() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("What is your primary task this pomodoro?")
	text, _ := reader.ReadString('\n')
	return text
}

type countdown struct {
	t int
	m int
	s int
}

func getTimeRemaining(t time.Time) countdown {
	currentTime := time.Now()
	difference := t.Sub(currentTime)

	total := int(difference.Seconds())
	minutes := int(total/60) % 60
	seconds := int(total % 60)

	return countdown{
		t: total,
		m: minutes,
		s: seconds,
	}
}

func setIntention(intention string) {
	f, err := os.OpenFile("pomodoro.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
	f.WriteString(intention)
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
}
