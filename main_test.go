package render

import (
	"bytes"
	"html/template"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		assert := assert.New(t)

		render := New("./testdata/simple", nil)
		assert.NotNil(render)

		var buf bytes.Buffer
		assert.NotPanics(func() {
			render(&buf, &foo{Value: "12345"})
		})

		assert.Equal("<layout>\n<foo>\n12345\n</foo>\n</layout>\n", buf.String())
	})

	t.Run("layouter", func(t *testing.T) {
		assert := assert.New(t)


		render := New("./testdata/layouter", nil)
		assert.NotNil(render)

		var buf bytes.Buffer
		assert.NotPanics(func() {
			render(&buf, &bar{Value: "12345"})
		})

		assert.Equal("<non-default-layout>\n<bar>\n12345\n</bar>\n</non-default-layout>\n", buf.String())
	})

	t.Run("templater", func(t *testing.T) {
		assert := assert.New(t)

		render := New("./testdata/templater", nil)
		assert.NotNil(render)

		var buf bytes.Buffer
		assert.NotPanics(func() {
			render(&buf, &baz{Value: "12345"})
		})

		assert.Equal("<layout>\n<non-default-template>\n12345\n</non-default-template>\n</layout>\n", buf.String())
	})

	t.Run("template not found", func(t *testing.T) {
		assert := assert.New(t)

		render := New("./testdata/template-not-found", nil)
		assert.NotNil(render)

		var buf bytes.Buffer
		assert.Panics(func() {
			render(&buf, &foo{Value: "12345"})
		})

		assert.Equal("", buf.String())
	})

	t.Run("layout not found", func(t *testing.T) {
		assert := assert.New(t)

		render := New("./testdata/layout-not-found", nil)
		assert.NotNil(render)

		var buf bytes.Buffer
		assert.Panics(func() {
			render(&buf, &foo{Value: "12345"})
		})

		assert.Equal("", buf.String())
	})
}

type foo struct{
	Value template.HTML
}

type bar struct{
	Value template.HTML
}

func (b *bar) Layout() string {
	return "non-default-layout.tmpl"
}

type baz struct{
	Value template.HTML
}

func (b *baz) Template() string {
	return "non-default-template.tmpl"
}