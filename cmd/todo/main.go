package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
	"todo"
)

var todoFileName = "todo.json"

// getTask decides where to get the description for a new task from: arguments or STDIN
// Each line in STDIN is a new task
// Empty line will end adding tasks
func getTask(r io.Reader, args ...string) ([]string, error) {
	tasks := []string{}
	if len(args) > 0 {
		tasks := append(tasks, strings.Join(args, ""))
		return tasks, nil
	}

	s := bufio.NewScanner(r)
	s.Split(bufio.ScanLines)
	for s.Scan() {
		if err := s.Err(); err != nil {
			return tasks, err
		}
		if s.Text() == "" {
			return tasks, nil
		}
		tasks = append(tasks, s.Text())
	}

	return tasks, nil
}

func main() {
	// Check if the user defined the ENV VAR for a custom file name
	if os.Getenv("TODO_FILENAME") != "" {
		todoFileName = os.Getenv("TODO_FILENAME")
	}
	// Parsing command line flags
	add := flag.Bool("add", false, "Add task to the ToDo list")
	list := flag.String("list", "", "List tasks. Options:\na - all, u - uncompleted, c - completed")
	complete := flag.Int("complete", 0, "Item to be completed")
	del := flag.Int("del", 0, "Item to be deleted")

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(),
			"%s tool. Developed by Massinja\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "Usage information:\n")
		flag.PrintDefaults()
	}
	flag.Parse()
	// Define an Items List
	l := &todo.List{}
	// If file with todo items already exists, get it
	if _, err := os.Stat(todoFileName); err == nil {
		if err := l.Get(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}

	switch {
	case *add:
		t, err := getTask(os.Stdin, flag.Args()...)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		for _, v := range t {
			l.Add(v)
		}

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

	case *list == "a":
		fmt.Print(l)

	case *list == "u":
		i := 0
		for _, task := range *l {
			if task.Done == false {
				i = i + 1
				fmt.Printf(" %v: %v\n", i, task.Task)
			}
		}

	case *list == "c":
		i := 0
		for _, task := range *l {
			if task.Done == true {
				i = i + 1
				fmt.Printf(" %v) %v: %v\n", i, task.CompletedAt.Format(time.DateTime), task.Task)
			}
		}

	case *del > 0:
		if err := l.Delete(*del); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		if err := l.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

	default:
		flag.Usage()
		os.Exit(0)

	}

}
