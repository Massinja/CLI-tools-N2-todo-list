package todo_test

import (
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
	task1 := "Buy Coffee"
	item1 := todo.Item{Task: task1}
	task2 := "Read the Book"
	item2 := todo.Item{Task: task2}
	l = append(l, item1)
	l = append(l, item2)
	l.Complete(1)
	if !l[0].Done {
		t.Errorf("task %s was not marked complete", task1)
	}
	if l[1].Done {
		t.Errorf("task %s task should not be complete yet", task2)
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

// TestSaveGet
func TestSaveGet(t *testing.T) {
	listSave := todo.List{}
	listSave = append(listSave, todo.Item{Task: "Task1"})
	listSave = append(listSave, todo.Item{Task: "Task2"})
	listSave = append(listSave, todo.Item{Task: "Task3"})

	testFile := "testingTodo"
	if err := listSave.Save(testFile); err != nil {
		t.Errorf("Could not save Items in the list: %v", err)
	}
	defer os.Remove(testFile)

	listGet := todo.List{}
	if err := listGet.Get(testFile); err != nil {
		t.Errorf("Could not get list Items: %v", err)
	}
	if length := len(listGet); length < 3 {
		t.Errorf("List was truncated")
	} else if length > 3 {
		t.Errorf("List is too long")
	}
	if listGet[0].Task != "Task1" {
		t.Errorf("Tasks were not saved properly")
	}

}
