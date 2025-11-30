package template

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"path/filepath"
	"strings"
	"time"

	"github.com/shopspring/decimal"
	"github.com/wiederin/go-invoicer/currency"
)

type Source interface {
	Load(name string) (string, error)
	List() ([]string, error)
}

type FSSource struct {
	fs fs.FS
}

func NewFSSource(filesystem fs.FS) *FSSource {
	return &FSSource{fs: filesystem}
}

func (s *FSSource) Load(name string) (string, error) {
	data, err := fs.ReadFile(s.fs, name)
	if err != nil {
		return "", fmt.Errorf("failed to load template %s: %w", name, err)
	}
	return string(data), nil
}

func (s *FSSource) List() ([]string, error) {
	var templates []string
	err := fs.WalkDir(s.fs, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && (strings.HasSuffix(path, ".html") || strings.HasSuffix(path, ".tmpl")) {
			templates = append(templates, path)
		}
		return nil
	})
	return templates, err
}

type EmbedSource struct {
	fs embed.FS
}

func NewEmbedSource(filesystem embed.FS) *EmbedSource {
	return &EmbedSource{fs: filesystem}
}

func (s *EmbedSource) Load(name string) (string, error) {
	data, err := s.fs.ReadFile(name)
	if err != nil {
		return "", fmt.Errorf("failed to load embedded template %s: %w", name, err)
	}
	return string(data), nil
}

func (s *EmbedSource) List() ([]string, error) {
	var templates []string
	err := fs.WalkDir(s.fs, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && (strings.HasSuffix(path, ".html") || strings.HasSuffix(path, ".tmpl")) {
			templates = append(templates, path)
		}
		return nil
	})
	return templates, err
}

type Manager struct {
	sources   []Source
	templates map[string]*template.Template
	funcMap   template.FuncMap
}

func NewManager(sources ...Source) *Manager {
	m := &Manager{
		sources:   sources,
		templates: make(map[string]*template.Template),
	}
	m.funcMap = m.defaultFuncMap()
	return m
}

func (m *Manager) defaultFuncMap() template.FuncMap {
	return template.FuncMap{
		"formatMoney": func(amount decimal.Decimal, currencyCode string) string {
			return currency.FormatSimple(amount, currencyCode)
		},
		"formatDate": func(t time.Time, layout string) string {
			if layout == "" {
				layout = "2006-01-02"
			}
			return t.Format(layout)
		},
		"formatDateLong": func(t time.Time) string {
			return t.Format("January 2, 2006")
		},
		"upper": strings.ToUpper,
		"lower": strings.ToLower,
		"title": strings.Title,
		"join":  strings.Join,
		"add": func(a, b int) int {
			return a + b
		},
		"mul": func(a, b float64) float64 {
			return a * b
		},
		"div": func(a, b float64) float64 {
			if b == 0 {
				return 0
			}
			return a / b
		},
		"seq": func(n int) []int {
			result := make([]int, n)
			for i := range result {
				result[i] = i + 1
			}
			return result
		},
	}
}

func (m *Manager) AddFunc(name string, fn any) {
	m.funcMap[name] = fn
}

func (m *Manager) Load(name string) error {
	for _, source := range m.sources {
		content, err := source.Load(name)
		if err == nil {
			tmpl, err := template.New(filepath.Base(name)).Funcs(m.funcMap).Parse(content)
			if err != nil {
				return fmt.Errorf("failed to parse template %s: %w", name, err)
			}
			m.templates[name] = tmpl
			return nil
		}
	}
	return fmt.Errorf("template %s not found in any source", name)
}

func (m *Manager) RenderHTML(name string, data any) (string, error) {
	tmpl, ok := m.templates[name]
	if !ok {
		if err := m.Load(name); err != nil {
			return "", err
		}
		tmpl = m.templates[name]
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template %s: %w", name, err)
	}
	return buf.String(), nil
}

func (m *Manager) RenderToBytes(name string, data any) ([]byte, error) {
	html, err := m.RenderHTML(name, data)
	if err != nil {
		return nil, err
	}
	return []byte(html), nil
}

func (m *Manager) ListTemplates() []string {
	var all []string
	for _, source := range m.sources {
		if list, err := source.List(); err == nil {
			all = append(all, list...)
		}
	}
	return all
}
