package main

import "testing"

func TestAddTag(t *testing.T) {
	entry, err := NewEntry(1, "Learn Go")
	if err != nil {
		t.Fatal("unexpected error creating entry:", err)
	}
	entry.AddTag("backend")

	if !entry.Tags["backend"] {
		t.Fatal("expected tag 'backend' to be added")
	}
}

func TestNewEntryRejectsBlankTitle(t *testing.T) {
	_, err	:= NewEntry(1, " ")

	if err == nil {
		t.Fatal("epected an error with a blank title")
	}
}
func TestNewEntryTrimsTitle(t *testing.T) {
	entry, err := NewEntry(1, "  Learn Go  ")
	if err != nil {
		t.Fatal(err)
	}

	if entry.Title != "Learn Go" {
		t.Fatalf("got %q, want %q", entry.Title, "Learn Go")
	}
}