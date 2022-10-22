package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"github.com/eidolon/wordwrap"

	"quick-rater/data"
)

func printPrompt(prompt data.Prompt) {
	wrapper := wordwrap.Wrapper(80, false)
	details := prompt.Get()

	wrappedTitle := wrapper("\033[1m" + details.ElementTitle + "\033[0m")
	formattedTitle := wordwrap.Indent(wrappedTitle, "What about:\t", false)

	formattedDescription := wrapper(details.ElementDetails)

	var answerType string
	if details.QuestionIsBinary {
		answerType = "[Y/n]"
	} else {
		answerType = "[1-5]"
	}

	fmt.Printf("\n%v\n%v\n\n%v %v\t", formattedTitle, formattedDescription, details.QuestionText, answerType)
}

func readResponse(prompt data.Prompt) (int, error) {
	scanner := bufio.NewScanner(os.Stdin)

	scanner.Scan()
	text := scanner.Text()

	if text == "" {
		return 0, fmt.Errorf("received empty string as input")
	}

	isBinary := prompt.Get().QuestionIsBinary

	// Binary question case
	if isBinary && (text == "Y" || text == "y") {
		return 5, nil
	} else if isBinary && (text == "N" || text == "n") {
		return 0, nil
	} else if isBinary {
		return 0, fmt.Errorf("invalid input for yes/no question, got %v", text)

	}

	// 5-star question case
	if rating, err := strconv.Atoi(text); err != nil {
		return 0, fmt.Errorf("invalid numeric input for 5-star question, got %v: %w", text, err)
	} else if rating < 1 || rating > 5 {
		return 0, fmt.Errorf("invalid numeric input outside range for 5-star question, got %v: %w", rating, err)
	} else {
		return rating, nil
	}
}

func main() {
	d := data.New()
	defer func() {
		d.Close()
		fmt.Println("Pringiedoobles!")
	}()

	for {
		prompt := d.Ask()
		var response int

		// Repeat prompt until a valid response is received
		var err error
		for answered := false; !answered; answered = err == nil {
			printPrompt(prompt)
			response, err = readResponse(prompt)
			if err != nil {
				fmt.Printf("\n\033[31m%v \033[0m\n", err)
			}
		}

		if err := d.Answer(prompt, response); err != nil {
			fmt.Println(fmt.Errorf("could not record response to prompt: %w", err))
		}
	}
}
