package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	store := new(Store)
	runShell(store)
}

func runShell(store *Store) {
	intr := Interpreter{store}
	fmt.Println("Type \"exit\" to leave")

	scanner := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("> ")

		txt, err := scanner.ReadString('\n')
		if err != nil {
			fmt.Println("Got the following error while retrieving input: ", err)
			continue
		}

		cmd := strings.TrimSpace(txt)

		switch cmd {
		case "":
		case "exit":
			fmt.Println("Exiting...")
			return
		default:
			if actual, err := intr.Exec(cmd); err == nil {
				fmt.Printf(" %v\n", actual)
			} else {
				fmt.Println("Command failed with the following error: ", err)
			}
		}
	}
}
