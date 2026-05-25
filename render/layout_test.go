package render

import "testing"

func TestLayoutViewShowAndOrd(t *testing.T) {
	l := NewLayoutView([]string{BlockTotals, BlockHeader, BlockLineItems})
	if !l.Show(BlockHeader) {
		t.Fatal("expected header")
	}
	if l.Show(BlockParties) {
		t.Fatal("parties not in layout")
	}
	if l.Ord(BlockTotals) != 0 {
		t.Fatalf("totals order: %d", l.Ord(BlockTotals))
	}
	if l.Ord(BlockHeader) != 1 {
		t.Fatalf("header order: %d", l.Ord(BlockHeader))
	}
}
