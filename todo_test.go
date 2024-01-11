package todo_test

import (
	"errors"
	"os"
	"testing"
	"todo"
)

// TestAdd tests the Add method of the type List
func TestAdd(t *testing.T) {
	l := todo.List{}
	taskName := "Buy Coffee"
	l.Add(taskName)
	if l[0].Task != taskName {
		t.Errorf("expected %v, got %v instead.", taskName, l[0].Task)
	}

}

// TestComplete
func TestComplete(t *testing.T) {
	l := todo.List{}
	taskName := "Buy Coffee"
	item1 := todo.Item{Task: taskName}
	l = append(l, item1)
	l = append(l, todo.Item{Task: "Read the Book"})
	l.Complete(1)
	if !l[0].Done {
		t.Errorf("new task should be completed")
	}
	if l[1].Done {
		t.Errorf("this task should not be complete yet")
	}
}

// TestDelete
func TestDelete(t *testing.T) {
	l := todo.List{}
	l = append(l, todo.Item{Task: "Task1"})
	l = append(l, todo.Item{Task: "Task2"})
	l = append(l, todo.Item{Task: "Task3"})

	if len(l) != 3 {
		t.Errorf("couldn't add to a List")
	}
	l.Delete(2)
	if len(l) == 3 {
		t.Errorf("couldn't delete a todo Item from a List")
	} else if len(l) < 2 {
		t.Errorf("more todo Items deleted from a List than expected")
	}
	if l[0].Task != "Task1" || l[1].Task != "Task3" {
		t.Errorf("wrong Item was deleted")
	}

}

// TestSave
func TestSave(t *testing.T) {
	l := todo.List{}
	l = append(l, todo.Item{Task: "Task1"})
	l = append(l, todo.Item{Task: "Task2"})
	l = append(l, todo.Item{Task: "Task3"})

	testFile := "testingTodo"
	l.Save(testFile)
	defer os.Remove(testFile)
	if _, err := os.Stat(testFile); err == nil {

	} else if errors.Is(err, os.ErrNotExist) {
		t.Errorf("List was not saved in a file")
	}

	
}

// TestGet
