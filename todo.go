package todo

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

type item struct {
	Task        string
	Done        bool
	CreatedAt   time.Time
	CompletedAt time.Time
}

type List []item

func checkErr(e error) {
	if e != nil {
		fmt.Errorf("%v", e)
		return
	}
}

// Add creates a new todo item and appends it to the list
func (l *List) Add(task string) {
	t := item{
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
	checkErr(err)

	er := os.WriteFile(filename, lj, 0644)
	checkErr(er)

	return nil
}

// Get opens the file, decodes JSON and parses it into a List
func (l *List) Get(filename string) error {

	f, err := os.Open(filename)
	checkErr(err)

	file, err := io.ReadAll(f)
	checkErr(err)
	defer f.Close()

	er := json.Unmarshal(file, l)
	checkErr(er)

	return nil
}
