// Package numbering formats auto-generated invoice numbers from org rules.
package numbering

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var seqPadRE = regexp.MustCompile(`\{SEQ:(\d+)\}`)

// DefaultPattern is used when the org has no custom pattern.
const DefaultPattern = "{PREFIX}{YYYY}-{SEQ:4}"

// Format builds an invoice number from pattern, prefix, sequence, and issue date.
// Supported tokens: {PREFIX}, {YYYY}, {YY}, {MM}, {SEQ}, {SEQ:N} (zero-padded width N).
func Format(pattern, prefix string, seq int, at time.Time) string {
	if strings.TrimSpace(prefix) == "" {
		prefix = "INV-"
	}
	pat := strings.TrimSpace(pattern)
	if pat == "" {
		pat = DefaultPattern
	}
	if seq < 1 {
		seq = 1
	}
	out := pat
	out = strings.ReplaceAll(out, "{PREFIX}", prefix)
	out = strings.ReplaceAll(out, "{YYYY}", fmt.Sprintf("%04d", at.Year()))
	out = strings.ReplaceAll(out, "{YY}", fmt.Sprintf("%02d", at.Year()%100))
	out = strings.ReplaceAll(out, "{MM}", fmt.Sprintf("%02d", int(at.Month())))
	out = strings.ReplaceAll(out, "{SEQ}", strconv.Itoa(seq))
	out = seqPadRE.ReplaceAllStringFunc(out, func(m string) string {
		sub := seqPadRE.FindStringSubmatch(m)
		if len(sub) < 2 {
			return m
		}
		w, _ := strconv.Atoi(sub[1])
		if w <= 0 {
			w = 4
		}
		return fmt.Sprintf("%0*d", w, seq)
	})
	return out
}
