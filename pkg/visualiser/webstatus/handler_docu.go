package webstatus

import (
	"html/template"
	"net/http"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
	"github.com/vogtp/som"
)

func (s *WebStatus) handleDocu(w http.ResponseWriter, r *http.Request) {
	p := parser.NewWithExtensions(
		parser.CommonExtensions |
			parser.Titleblock |
			parser.Mmark,
	)
	output := markdown.ToHTML(som.README, p, nil)
	var data = struct {
		*commonData
		Docu template.HTML
	}{
		commonData: common("", r),
		Docu:       template.HTML(string(output)),
	}
	s.render(w, r, "docu.gohtml", data)
}
