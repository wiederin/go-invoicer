package pdf

import "context"

// Renderer converts HTML invoice markup to PDF bytes.
type Renderer interface {
	RenderHTML(ctx context.Context, html string) ([]byte, error)
}

// ErrNotImplemented is returned by stub renderers.
type ErrNotImplemented struct{}

func (ErrNotImplemented) Error() string {
	return "pdf: renderer not implemented; use platform renderer or inject a Chromium backend"
}

// StubRenderer is a placeholder until a local Chromium adapter exists.
type StubRenderer struct{}

func (StubRenderer) RenderHTML(_ context.Context, _ string) ([]byte, error) {
	return nil, ErrNotImplemented{}
}
