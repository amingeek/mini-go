package main

import (
	"strings"
	"testing"
)

func TestGetList(t *testing.T) {
	input := "hat coat scarf\n"
	reader := strings.NewReader(input)
	got := GetList(reader)
	want := []string{"coat", "scarf"}

	if len(got) != len(want) {
		t.Fatalf("Expected %d items, got %d", len(want), len(got))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("Expected %v, got %v", want[i], got[i])
		}
	}
}

func TestGetSeason(t *testing.T) {
	input := "winter cold dry\n"
	reader := strings.NewReader(input)
	got := GetSeason(reader)
	want := "winter"

	if got != want {
		t.Errorf("Expected %v, got %v", want, got)
	}
}
