package main_test

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
	"todo"
)

var (
	binName  = "todo"
	fileName = "todo.json"
)

func TestMain(m *testing.T) {
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
		fmt.Fprintf(os.Stderr, "Cannot build tool %s: %s", binName, err)
		os.Exit(1)
	}

	fmt.Println("Running tests...")
	result := m.Run(binName, func(m *testing.T) {})
	fmt.Println("result:", result)
	fmt.Println("Cleaning up...")
	os.Remove(fileName)

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
	t.Run("ListTasks", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-list")
		out, err := cmd.Output()
		if err != nil {
			t.Fatal(err)
		}
		expected := " 1: " + task + "\n" + " 2: " + task2 + "\n"
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
		err := l.Get(fileName)
		todo.CheckErr(err)
		if !l[0].Done {
			t.Errorf("task was not marked as complete")
		}
	})
	t.Run("RemoveTodoList", func(t *testing.T) {
		os.Remove(fileName)
		os.Remove(binName)
	})

}
