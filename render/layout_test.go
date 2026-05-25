package render

import "testing"

func TestLayoutViewEmptyShowsAllExceptPayment(t *testing.T) {
	l := NewLayoutView(nil)
	for _, id := range []string{BlockHeader, BlockParties, BlockLineItems, BlockTotals, BlockNotes} {
		if !l.Show(id) {
			t.Fatalf("expected %s visible with default layout", id)
		}
	}
	if l.Show(BlockPayment) {
		t.Fatal("payment hidden by default")
	}
}

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
