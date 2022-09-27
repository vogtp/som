package webstatus

import (
	"fmt"
	"html/template"
	"math"
	"net/http"
	"strconv"
)

const pageSize = 10

type pageInfo struct {
	ID    template.HTML
	State string
	URL   string
}

func (s *WebStatus) getPages(r *http.Request, total int) ([]pageInfo, int) {

	var pages []pageInfo
	page := 1
	r.ParseForm()
	if str := r.Form.Get("page"); len(str) > 0 {
		if p, err := strconv.Atoi(str); err == nil {
			page = p
		} else {
			s.hcl.Warnf("Cannot parse page %q", p)
		}
	}
	offset := (page - 1) * pageSize
	s.logTime("Paging offset %v len %v total %v", offset, pageSize, total)

	pgCnt := int(math.Ceil(float64(total / pageSize)))
	url := r.URL

	r.Form.Set("page", fmt.Sprintf("%d", page-1))
	r.URL.RawQuery = r.Form.Encode()
	p := pageInfo{
		ID:  template.HTML("&laquo;"),
		URL: url.String(),
	}
	if page < 2 {
		p.State = "disabled"
	}
	pages = append(pages, p)

	for i := 1; i <= pgCnt+1; i++ {
		id := fmt.Sprintf("%v", i)
		r.Form.Set("page", id)
		r.URL.RawQuery = r.Form.Encode()
		p := pageInfo{
			ID:  template.HTML(id),
			URL: url.String(),
		}
		if i == page {
			p.State = "active"
		}
		pages = append(pages, p)
	}
	if len(pages) > 18 {
		dots := pageInfo{ID: "...", State: "disabled"}
		start := 9
		end := len(pages) - 9
		mid := []pageInfo{dots}
		if !(page < start || page > end) {
			start -= 3
			end += 3
			mid = append(mid, pages[page-3:page+3]...)
			mid = append(mid, dots)
		}
		backP := pages[end:]
		pages = append(pages[:start], mid...)
		pages = append(pages, backP...)
	}
	r.Form.Set("page", fmt.Sprintf("%d", page+1))
	r.URL.RawQuery = r.Form.Encode()
	p = pageInfo{
		ID:  template.HTML("&raquo;"),
		URL: url.String(),
	}
	if page >= pgCnt+1 {
		p.State = "disabled"
	}
	pages = append(pages, p)

	r.Form.Del("page")
	r.URL.RawQuery = r.Form.Encode()

	if len(pages) < 4 {
		// 3 pages mean 1 with information -> no paging
		pages = pages[:0]
	}
	return pages, offset
}
