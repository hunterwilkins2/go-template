package templates

import (
	"bytes"
	"errors"
	"html/template"
	"net/http"
	"path/filepath"
)

var (
	ErrTemplateNotFound = errors.New("template was not found")
)

// HTMLTemplate caches html templates
type HTMLTemplate struct {
	templates map[string]*template.Template
	hotReload bool
}

var functions = template.FuncMap{}

// New Creates a new HTMLTemplate cache
func New(dir string, hotReload bool) (*HTMLTemplate, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob(filepath.Join(dir, "pages/*.html"))
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.New(name).Funcs(functions).ParseFiles(filepath.Join(dir, "base.html"))
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob(filepath.Join(dir, "layout/*"))
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob(filepath.Join(dir, "components/*.html"))
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return &HTMLTemplate{templates: cache, hotReload: hotReload}, nil
}

// templateData executes template with common data
// any data passed to render will be put in the Data field
// access this with {{.Data.(struct name)}}
type templateData struct {
	Data      any
	HotReload bool
}

// Render renders HTML template and puts rendered HTML into response
func (htmlTs *HTMLTemplate) Render(w http.ResponseWriter, templateName string, data any) error {
	ts, ok := htmlTs.templates[templateName]
	if !ok {
		return ErrTemplateNotFound
	}

	var buff bytes.Buffer
	err := ts.ExecuteTemplate(&buff, "base.html", templateData{Data: data, HotReload: htmlTs.hotReload})
	if err != nil {
		return err
	}

	w.Write(buff.Bytes())
	return nil
}
