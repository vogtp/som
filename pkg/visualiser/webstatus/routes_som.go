package webstatus

import (
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/vogtp/som/pkg/core"
	"github.com/vogtp/som/pkg/visualiser/webstatus/database/graphql"
)

func (s *WebStatus) routes() {
	w := core.Get().WebServer()

	w.Handle("/static/", http.StripPrefix("/som/", http.FileServer(http.FS(assetData))))
	w.HandleFunc(FilesPath, s.handleFiles)

	// handlers below this line will be reported in the log
	w.AddMiddleware(s.reportRequest)
	w.HandleFunc("/", s.handleIndex)
	w.HandleFunc("/docu", s.handleDocu)
	w.HandleFunc("/topology/", s.handleTopology)
	w.HandleFunc("/mesh_info/", s.handleMeshInfo)
	w.HandleFunc(alertListPath, s.handleAlertList)
	w.HandleFunc(AlertDetailPath, s.handleAlertDetail)
	w.HandleFunc(incidentListPath, s.handleIncidentList)
	w.HandleFunc(IncidentDetailPath, s.handleIncidentDetail)

	w.Handle("/gql", playground.Handler("SQM", w.BasePath()+"/api"))

	srv := handler.NewDefaultServer(graphql.NewSchema(s.Ent()))
	w.Handle("/api", srv)

	// legacy
	// w.HandleFunc("/chart/", s.handleChart)
}

func (s *WebStatus) reportRequest(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		f(w, r)
		d := time.Since(start)
		d = d.Truncate(time.Microsecond)
		if d > time.Millisecond {
			d = d.Truncate(time.Millisecond)
		}
		if d > 100*time.Millisecond {
			d = d.Truncate(10 * time.Millisecond)
		}
		s.hcl.Infof("Request (%v): %v", d, r.URL.String())
	}
}
