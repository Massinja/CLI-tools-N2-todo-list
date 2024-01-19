package todo

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Item struct {
	Task        string
	Done        bool
	CreatedAt   time.Time
	CompletedAt time.Time
}

type List []Item

// Add creates a new todo item and appends it to the list
func (l *List) Add(task string) {
	t := Item{
		Task:        task,
		Done:        false,
		CreatedAt:   time.Now(),
		CompletedAt: time.Time{},
	}

	*l = append(*l, t)
}

// Complete method marks a todo item as completed by
// Done = true, CompletedAt = current time
func (l *List) Complete(i int) error {
	ls := *l
	if i <= 0 || i > len(ls) {
		return fmt.Errorf("Item %d does not exist", i)
	}

	ls[i-1].Done = true
	ls[i-1].CompletedAt = time.Now()

	return nil
}

// Delete removes a todo item from the list
func (l *List) Delete(i int) error {
	ls := *l
	if i <= 0 || i > len(ls) {
		return fmt.Errorf("Item %d does not exist", i)
	}

	*l = append(ls[:i-1], ls[i:]...)
	return nil

}

// Save encodes the List as JSON and saves it
func (l *List) Save(filename string) error {

	lj, err := json.Marshal(l)
	if err != nil {
		return fmt.Errorf("couldn't convert to json: %v", err)
	}

	err = os.WriteFile(filename, lj, 0644)
	if err != nil {
		return fmt.Errorf("couldn't write to file %s: %v", filename, err)
	}

	return nil
}

// Get opens the file, decodes JSON and parses it into a List
func (l *List) Get(filename string) error {

	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		return fmt.Errorf("couldn't open file %s: %v", filename, err)
	}

	file, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	defer f.Close()

	if err := json.Unmarshal(file, l); err != nil {
		return fmt.Errorf("couldn't convert from json file %s: %v", file, err)
	}

	return nil
}

func (l *List) String() string {
	formatted := ""
	for _, t := range *l {
		if t.Done {
			formatted += fmt.Sprintf(" X: %s\n", t.Task)
		} else {
			formatted += fmt.Sprintf(" O: %s\n", t.Task)
		}
	}
	return formatted
}
