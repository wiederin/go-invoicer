package numbering

import (
	"testing"
	"time"
)

func TestFormat(t *testing.T) {
	at := time.Date(2026, 3, 15, 0, 0, 0, 0, time.UTC)
	got := Format("{PREFIX}{YYYY}-{SEQ:4}", "ACME-", 42, at)
	want := "ACME-2026-0042"
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}
