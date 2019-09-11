package render

import (
	"fmt"
	"html/template"
	"io"
	"path/filepath"
	"reflect"
)

// New returns a render function with the given template path and function map.
func New(path string, funcMap template.FuncMap) func(io.Writer, interface{}) {
	ts := template.Must(template.New("").Funcs(funcMap).ParseGlob(filepath.Join(path, "*.tmpl")))

	var render func(io.Writer, interface{})
	render = func(w io.Writer, c interface{}) {
		templates := template.Must(ts.Clone())

		var name string
		if t, ok := c.(Templater); ok {
			name = t.Template()
		} else {
			name = templateName(c)
		}

		layout := "layout.tmpl"
		if l, ok := c.(Layouter); ok {
			layout = l.Layout()
		}

		t := templates.Lookup(name)
		if t == nil {
			panic(fmt.Errorf("template not found: %s", name))
		}

		t = template.Must(t.AddParseTree("content", t.Tree))
		if err := t.ExecuteTemplate(w, layout, c); err != nil {
			panic(err)
		}
	}
	return render
}

func templateName(c interface{}) string {
	t := reflect.TypeOf(c)
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return fmt.Sprintf("%s.tmpl", t)
}

// Templater is an interface of view objects which have a custom template name other than <type>.tmpl.
type Templater interface {
	Template() string
}

// Layouter is an interface of view objects which have a custom layout other than layout.tmpl.
type Layouter interface {
	Layout() string
}

