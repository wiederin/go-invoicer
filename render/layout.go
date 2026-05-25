package render

// Default layout block IDs (invoice section order).
const (
	BlockHeader     = "header"
	BlockParties    = "parties"
	BlockLineItems  = "line_items"
	BlockTotals     = "totals"
	BlockNotes      = "notes"
	BlockPayment    = "payment" // Swiss QR-bill and payment details
)

var defaultBlockOrder = map[string]int{
	BlockHeader: 0, BlockParties: 1, BlockLineItems: 2, BlockTotals: 3, BlockNotes: 4, BlockPayment: 5,
}

// DefaultLayoutBlocks is the standard section order for block-based templates.
func DefaultLayoutBlocks() []string {
	return []string{BlockHeader, BlockParties, BlockLineItems, BlockTotals, BlockNotes}
}

// SwissDefaultLayoutBlocks includes payment / QR section.
func SwissDefaultLayoutBlocks() []string {
	return append(DefaultLayoutBlocks(), BlockPayment)
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
	if len(l.order) == 0 {
		if name == BlockPayment {
			return false
		}
		return true
	}
	for _, id := range l.order {
		if id == name {
			return true
		}
	}
	return false
}

// Ord returns flex order for a block (unknown blocks sort last).
func (l LayoutView) Ord(name string) int {
	if len(l.order) == 0 {
		if o, ok := defaultBlockOrder[name]; ok {
			return o
		}
		return 99
	}
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
