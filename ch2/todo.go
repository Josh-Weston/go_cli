package todo

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"
)

type item struct {
	Task        string
	Done        bool
	CreatedAt   time.Time
	CompletedAt time.Time
}

type List []item

func (l *List) Add(task string) {
	t := item{
		Task:        task,
		Done:        false,
		CreatedAt:   time.Now(),
		CompletedAt: time.Time{},
	}
	*l = append(*l, t)
}

func (l *List) Complete(i int) error {
	ls := *l
	if i <= 0 || i > len(ls) {
		return fmt.Errorf("Item %d does not exist", i)
	}

	// adjusting for 0 based indexing
	ls[i-1].Done = true
	ls[i-1].CompletedAt = time.Now()
	return nil
}

func (l *List) Delete(i int) error {
	ls := *l
	if i <= 0 || i > len(ls) {
		return fmt.Errorf("Item %d does not exist", i)
	}

	// adjusting for 0 based indexing
	*l = append(ls[:i-1], ls[i:]...)
	return nil
}

func (l *List) Save(fileName string) error {
	js, err := json.Marshal(l)
	if err != nil {
		return err
	}
	return os.WriteFile(fileName, js, 0644)
}

func (l *List) Get(fileName string) error {
	file, err := os.ReadFile(fileName)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}
	if len(file) == 0 {
		return nil
	}
	return json.Unmarshal(file, l)
}

func (l *List) String() string {
	var sb strings.Builder
	for k, t := range *l {
		prefix := "  "
		if t.Done {
			prefix = "X "
		}
		// Adjust to be indexed starting at 1
		sb.WriteString(fmt.Sprintf("%s%d: %s\n", prefix, k+1, t.Task))
	}
	return sb.String()
}

func (l *List) StringVerbose() string {
	var sb strings.Builder
	for k, t := range *l {
		prefix := "  "
		completedAt := ""
		if t.Done {
			prefix = "X "
			completedAt = fmt.Sprintf(", CompletedAt: %s", t.CompletedAt)
		}
		sb.WriteString(fmt.Sprintf("%s%d: %s (Created At: %s%s)\n", prefix, k+1, t.Task, t.CreatedAt, completedAt))
	}
	return sb.String()
}

func (l *List) ShowInComplete() string {
	var sb strings.Builder
	for k, t := range *l {
		prefix := "  "
		if t.Done {
			continue
		}
		// Adjust to be indexed starting at 1
		sb.WriteString(fmt.Sprintf("%s%d: %s\n", prefix, k+1, t.Task))
	}
	return sb.String()
}
