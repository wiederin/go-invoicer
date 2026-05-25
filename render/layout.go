package render

// Default layout block IDs (invoice section order).
const (
	BlockHeader     = "header"
	BlockParties    = "parties"
	BlockLineItems  = "line_items"
	BlockTotals     = "totals"
	BlockNotes      = "notes"
)

// DefaultLayoutBlocks is the standard section order for block-based templates.
func DefaultLayoutBlocks() []string {
	return []string{BlockHeader, BlockParties, BlockLineItems, BlockTotals, BlockNotes}
}

// LayoutView controls section visibility and CSS flex order in templates.
type LayoutView struct {
	order []string
}

// NewLayoutView returns a layout helper; empty blocks uses defaults.
func NewLayoutView(blocks []string) LayoutView {
	if len(blocks) == 0 {
		blocks = DefaultLayoutBlocks()
	}
	out := make([]string, 0, len(blocks))
	seen := make(map[string]bool, len(blocks))
	for _, id := range blocks {
		if id == "" || seen[id] {
			continue
		}
		seen[id] = true
		out = append(out, id)
	}
	return LayoutView{order: out}
}

// Show reports whether a block id is included in the layout.
func (l LayoutView) Show(name string) bool {
	for _, id := range l.order {
		if id == name {
			return true
		}
	}
	return false
}

// Ord returns flex order for a block (unknown blocks sort last).
func (l LayoutView) Ord(name string) int {
	for i, id := range l.order {
		if id == name {
			return i
		}
	}
	return len(l.order) + 1
}

// WithLayout attaches block layout to a view.
func WithLayout(view InvoiceView, blocks []string) InvoiceView {
	view.Layout = NewLayoutView(blocks)
	return view
}
