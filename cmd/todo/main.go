package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"todo"
)

const todoFileName = "todo.json"

func main() {

	// Parsing command line flags
	task := flag.String("task", "", "Task to be included in the todo list")
	list := flag.Bool("list", false, "List all tasks")
	complete := flag.Int("complete", 0, "Item to be completed")

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(),
			"%s tool. Developed by Massinja\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "Usage information:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	//Define an Items List
	l := &todo.List{}

	// Use the Get command to read todo otems from a file
	if err := l.Get(todoFileName); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	switch {
	case *task != "":
		firstString := *task

		// temporarily solution.
		// In case user doesn't use "" to join multiple words to describe a task
		allStrings := flag.Args()
		allStrings = append([]string{firstString}, allStrings...)
		l.Add(strings.Join(allStrings, " "))
		if err := l.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

	case *complete > 0:
		if err := l.Complete(*complete); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		if err := l.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

	case *list:
		for _, item := range *l {
			if !item.Done {
				fmt.Println(item.Task)
			}
		}
	default:
		flag.Usage()
		os.Exit(0)

	}

}
