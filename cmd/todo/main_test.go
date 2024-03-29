package main_test

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
	"time"
	"todo"
)

var (
	binName  = "todo"
	fileName = "todo.json"
)

func TestMain(m *testing.M) {
	if os.Getenv("TODO_FILENAME") != "" {
		fileName = os.Getenv("TODO_FILENAME")
	}
	// Remove todo list to start clean
	os.Remove(fileName)
	fmt.Println("Building tool...")

	if runtime.GOOS == "windows" {
		binName += ".exe"
	}

	build := exec.Command("go", "build", "-o", binName)

	if err := build.Run(); err != nil {
		fmt.Errorf("Cannot build tool %s: %s", binName, err)
		return
	}
	fmt.Println("Running tests...")
	result := m.Run()

	fmt.Println("Cleaning up...")
	os.Remove(fileName)
	os.Remove(binName)

	os.Exit(result)
}

func TestTodoCLI(t *testing.T) {
	task := "task number 1"
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	cmdPath := filepath.Join(dir, binName)
	t.Run("AddNewTaskFromArguments", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-add", task)
		if err := cmd.Run(); err != nil {
			t.Fatal(err)
		}
	})
	task2 := "task number 2"
	t.Run("AddNewTaskFromSTDIN", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-add")
		cmdStdIn, err := cmd.StdinPipe()
		if err != nil {
			t.Fatal(err)
		}
		io.WriteString(cmdStdIn, task2)
		cmdStdIn.Close()
		if err := cmd.Run(); err != nil {
			t.Fatal(err)
		}
	})
	t.Run("ListAllTasks", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-list", "a")
		out, err := cmd.Output()
		if err != nil {
			t.Fatal(err)
		}
		expected := fmt.Sprintf(" O: %s\n O: %s\n", task, task2)
		if expected != string(out) {
			t.Errorf("Expected %q, got %q instead\n", expected, string(out))
		}
	})
	t.Run("CompleteTask", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-complete", "1")
		if err := cmd.Run(); err != nil {
			t.Fatal(err)
		}
	})
	t.Run("CheckCompleteTask", func(t *testing.T) {
		l := todo.List{}
		if err := l.Get(fileName); err != nil {
			t.Errorf("Could not get the list: %v", err)
			return
		}
		if !l[0].Done {
			t.Errorf("Task was not marked as complete")
		}
		timecomp := l[0].CompletedAt.Format(time.DateTime)
		cmd := exec.Command(cmdPath, "-list", "c")
		out, err := cmd.Output()
		if err != nil {
			t.Fatal(err)
		}

		expected := fmt.Sprintf(" 1) %v: %v\n", timecomp, task)
		if expected != string(out) {
			t.Errorf("Expected %q, got %q instead\n", expected, string(out))
		}
	})
	t.Run("ListUncompletedTasks", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-list", "u")
		out, err := cmd.Output()
		if err != nil {
			t.Fatal(err)
		}
		expected := " 1: " + task2 + "\n"
		if expected != string(out) {
			t.Errorf("Expected %q, got %q instead\n", expected, string(out))
		}
	})

	t.Run("DeleteTask", func(t *testing.T) {
		l := &todo.List{}
		if err := l.Get(fileName); err != nil {
			t.Fatal(err)
		}
		length := len(*l)

		cmd := exec.Command(cmdPath, "-del", "1")
		if err := cmd.Run(); err != nil {
			t.Fatal(err)
		}
		if err := l.Get(fileName); err != nil {
			t.Fatal(err)
		}

		if len(*l) == length {
			t.Errorf("Item was not deleted")
		}
		tasks := *l

		if len(*l) < length && tasks[0].Task != task2 {
			t.Errorf("Wrong item was deleted")
		}
	})

	t.Run("AddMultipleTasksFromSTDIN", func(t *testing.T) {
		l := &todo.List{}
		if err := l.Get(fileName); err != nil {
			t.Fatal(err)
		}
		length := len(*l)

		cmd := exec.Command(cmdPath, "-add")
		cmdStdIn, err := cmd.StdinPipe()
		if err != nil {
			t.Fatal(err)
		}
		task3 := "task 3\n"
		task4 := "task 4\n"
		io.WriteString(cmdStdIn, task3)
		io.WriteString(cmdStdIn, task4)
		cmdStdIn.Close()
		if err := cmd.Run(); err != nil {
			t.Fatal(err)
		}
		if err := l.Get(fileName); err != nil {
			t.Fatal(err)
		}
		if len(*l) == length {
			t.Errorf("Items were not added")
		}
		if len(*l) == length+1 {
			t.Errorf("Only one item was added instead of two")
		}

	})
}
