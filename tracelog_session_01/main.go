package main

import (
	"errors"
	"fmt"
	"strings"
)

type Entry struct {
	ID    int
	Title string
	Tags  map[string]bool
}

func (e *Entry) AddTag(tag string) {
	if e.Tags == nil {
		e.Tags = make(map[string]bool)
	}
	e.Tags[tag] = true
}

func NewEntry(id int, title string) (*Entry, error) {
	if strings.TrimSpace(title) == "" {
		return nil, errors.New("title cannot be blank")
	}

	return &Entry{
		ID:    id,
		Title: strings.TrimSpace(title),
	}, nil
}

func (e *Entry) Summary() string {
	return fmt.Sprintf("Entry %d: %s", e.ID, e.Title)
}

func main() {
	entry, err := NewEntry(1, "My First Entry")
	if err != nil {
		fmt.Println("Error creating entry:", err)
		return
	}

	fmt.Println(entry.Summary())
}
