package templates

import (
	"bytes"
	"errors"
	"html/template"
	"net/http"
	"path/filepath"
	"time"

	"github.com/justinas/nosurf"
)

var (
	ErrTemplateNotFound = errors.New("template was not found")
)

// HTMLTemplate caches html templates
type HTMLTemplate struct {
	templates  map[string]*template.Template
	components *template.Template
	hotReload  bool
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

	components, err := template.ParseGlob(filepath.Join(dir, "components/*.html"))
	if err != nil {
		return nil, err
	}

	return &HTMLTemplate{templates: cache, components: components, hotReload: hotReload}, nil
}

// templateData executes template with common data
// any data passed to render will be put in the Data field
// access this with {{.Data.(struct name)}}
type templateData struct {
	Data        any
	CurrentYear int
	CSRFToken   string
	HotReload   bool
}

// Render renders HTML template and puts rendered HTML into response
func (htmlTs *HTMLTemplate) Render(w http.ResponseWriter, r *http.Request, templateName string, data any) error {
	ts, ok := htmlTs.templates[templateName]
	if !ok {
		return ErrTemplateNotFound
	}

	tData := templateData{
		Data:        data,
		CurrentYear: time.Now().Year(),
		CSRFToken:   nosurf.Token(r),
		HotReload:   htmlTs.hotReload,
	}

	var buff bytes.Buffer
	err := ts.ExecuteTemplate(&buff, "base.html", tData)
	if err != nil {
		return err
	}

	w.Write(buff.Bytes())
	return nil
}

// RenderComponent renders a partial HTML component
func (htmlTs *HTMLTemplate) RenderComponent(w http.ResponseWriter, r *http.Request, templateName string, data any) error {
	tData := templateData{
		Data:        data,
		CurrentYear: time.Now().Year(),
		CSRFToken:   nosurf.Token(r),
		HotReload:   htmlTs.hotReload,
	}

	var buff bytes.Buffer
	err := htmlTs.components.ExecuteTemplate(&buff, templateName, tData)
	if err != nil {
		return err
	}

	w.Write(buff.Bytes())
	return nil
}
