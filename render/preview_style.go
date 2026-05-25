package render

import (
	"fmt"
	"strings"
)

// InjectLayoutDesignerStyles adds preview-only CSS for template studio block highlighting.
func InjectLayoutDesignerStyles(html string, activeBlock string, visibleBlocks []string) string {
	if html == "" {
		return html
	}
	visible := make(map[string]bool, len(visibleBlocks))
	for _, id := range visibleBlocks {
		visible[id] = true
	}
	all := []string{BlockHeader, BlockParties, BlockLineItems, BlockTotals, BlockNotes, BlockPayment}
	var rules []string
	rules = append(rules,
		"[data-block-id]{transition:outline .15s ease,opacity .15s ease}",
		"[data-block-id]:hover{outline:2px dashed rgba(56,168,255,.45);outline-offset:2px}",
	)
	for _, id := range all {
		if !visible[id] {
			rules = append(rules, fmt.Sprintf(`[data-block-id="%s"]{opacity:.35;outline:2px dashed rgba(148,163,184,.5)}`, id))
		}
	}
	if activeBlock != "" {
		rules = append(rules, fmt.Sprintf(`[data-block-id="%s"]{outline:3px solid #38a8ff!important;outline-offset:3px;opacity:1!important}`, activeBlock))
	}
	style := fmt.Sprintf(`<style id="gi-layout-designer">%s</style>`, strings.Join(rules, ""))
	if i := strings.Index(html, "</head>"); i >= 0 {
		return html[:i] + style + html[i:]
	}
	return style + html
}
