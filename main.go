package main

import (
	"fmt"

	"quick-rater/data"
	"quick-rater/input"
)

func main() {
	input.Init()
	defer input.Close()

	d := data.New()
	defer d.Close()

	// Create queue to iterate over with initial random prompt value
	q := []data.Prompt{d.Ask()}

	// Store the last asked question for ability to "go back" once
	lastQuestion := q[0]

	for {
		// Pop-left from queue
		prompt := q[0]

		// Repeat prompt until a valid response is received, using the fact that the prompt is still first in queue
		var err error
		key, err := input.Get(prompt)
		if err != nil {
			fmt.Printf("\n\033[31m%v \033[0m\n", err)
			continue
		}

		fmt.Println()

		// Exit on escape key
		if key.IsEscape() {
			return
		}

		// Go back on backspace
		if key.IsBackspace() {
			// Allow only one jump back to the last question
			if q[0] != lastQuestion {
				q = append([]data.Prompt{lastQuestion}, q...)
			}

			continue
		}

		// Get numeric response data
		numeric, ok := key.GetNumeric()
		// Shouldn't be possible to not have numeric value given mutually exclusive conditions
		if !ok {
			fmt.Printf("\n\033[31m%v \033[0m\n", "cannot get numeric value from input")
			continue
		}

		// Record response
		if err := d.Answer(prompt, numeric); err != nil {
			fmt.Println(fmt.Errorf("could not record response to prompt: %w", err))
		}

		// Update queue and lastQuestion for next prompt
		lastQuestion = q[0]
		q = append(q[1:], d.Ask())
	}
}
