package input

import (
	"fmt"
	"os"
	"strconv"

	"github.com/eidolon/wordwrap"
	term "github.com/nsf/termbox-go"

	"quick-rater/data"
)

func reset() {
	fmt.Fprint(os.Stdout, "\r \r")
}

func Init() {
	err := term.Init()

	if err != nil {
		panic(err)
	}
}

func Close() {
	term.Close()
}

func validate(isBinary bool, key Key) error {
	if key.IsBackspace() || key.IsEscape() {
		return nil
	}

	if _, ok := key.GetBool(); isBinary && !ok {
		return fmt.Errorf("invalid input for yes/no question")
	} else if _, ok := key.GetRating(); !isBinary && !ok {
		return fmt.Errorf("invalid numeric input for 5-star question")
	}

	return nil
}

func generatePromptText(prompt data.Prompt) (intro, question string) {
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

	intro = fmt.Sprintf("\n%v\n%v\n\n", formattedTitle, formattedDescription)
	question = fmt.Sprintf("%v %v\t", details.QuestionText, answerType)

	return intro, question
}

func Get(prompt data.Prompt) (Key, error) {
	intro, question := generatePromptText(prompt)

	fmt.Print(intro, question)

	var res Key

	for {
		switch ev := term.PollEvent(); ev.Type {
		case term.EventKey:
			switch ev.Key {
			case term.KeyEnter:
				if err := validate(prompt.Get().QuestionIsBinary, res); err != nil {
					return nil, fmt.Errorf("could not read input: %w", err)
				}

				return res, nil
			case term.KeyEsc:
				return escapeKey{}, nil
			case term.KeyBackspace, term.KeyBackspace2:
				return backspaceKey{}, nil
			case 0:
				if ev.Ch == 89 {
					reset()
					fmt.Print(question + "Y")
					res = boolKey{true}
				} else if ev.Ch == 121 {
					reset()
					fmt.Print(question + "y")
					res = boolKey{true}
				} else if ev.Ch == 78 {
					reset()
					fmt.Print(question + "N")
					res = boolKey{false}
				} else if ev.Ch == 110 {
					reset()
					fmt.Print(question + "n")
					res = boolKey{false}
				} else if ev.Ch >= 49 && ev.Ch <= 53 {
					reset()
					rating := int(ev.Ch) - 48
					fmt.Print(question + strconv.Itoa(rating))
					res = ratingKey{rating}
				}
			}
		case term.EventError:
			panic(ev.Err)
		}
	}
}
