package webstatus

import (
	_ "embed" // embed needs it
	"html/template"
	"net/http"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
)

var (
	//go:embed README.md
	readme []byte
)

func (s *WebStatus) handleDocu(w http.ResponseWriter, r *http.Request) {
	p := parser.NewWithExtensions(
		parser.CommonExtensions |
			parser.Titleblock |
			parser.Mmark,
	)
	output := markdown.ToHTML(readme, p, nil)
	var data = struct {
		*commonData
		Docu template.HTML
	}{
		commonData: common("", r),
		Docu:       template.HTML(string(output)),
	}
	s.render(w, r, "docu.gohtml", data)
}
