package webstatus

import (
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/vogtp/som/pkg/core"
	"github.com/vogtp/som/pkg/visualiser/webstatus/api"
)

func (s *WebStatus) routes() {
	w := core.Get().WebServer()

	w.Handle("/static/", http.StripPrefix(w.BasePath(), http.FileServer(http.FS(assetData))))
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

	w.HandleFunc(szenarioComponentPath, s.handleSzenarioComponent)

	w.HandleFunc("/api/", s.handleGraphiQL)
	w.Handle("/graphiql/", playground.Handler("SQM", w.BasePath()+"/graphql/"))

	srv := handler.NewDefaultServer(api.NewSchema(s.Ent()))
	w.Handle("/graphql/", srv)

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
		s.log.Info("Request finished", "duration", d, "url", r.URL.String())
	}
}
