package repl

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/peterh/liner"
)

type Cmd func(args []string) error

func Run(prompt string, commands map[string]Cmd) {

	line := liner.NewLiner()
	defer line.Close()

	var keys []string

	for k := range commands {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	line.SetCompleter(func(line string) []string {
		var completions []string
		for _, k := range keys {
			if strings.HasPrefix(k, line) {
				completions = append(completions, k)
			}
		}
		return completions
	})

REPL:
	for {
		var err error
		var command string
		command, err = line.Prompt(prompt)
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("error reading line: ", err)
			continue
		}

		args := strings.Fields(command)

		if len(args) == 0 {
			continue
		}

		line.AppendHistory(command)

		switch args[0] {

		case "h", "help":
			fmt.Println("h[elp] -- this help")
			fmt.Println("q[uit] -- quit")
			fmt.Println("known commands:\n", strings.Join(keys, "\n"))

		case "q", "quit":
			break REPL

		default:

			if f, ok := commands[args[0]]; ok {
				err := f(args[1:])
				if err != nil {
					fmt.Println("error: ", err)
				}
			} else {
				fmt.Println("unknown command, try `help`")
			}
		}
	}

	fmt.Println("bye")
}
