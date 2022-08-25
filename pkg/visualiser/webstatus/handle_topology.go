package webstatus

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"html/template"
	"image/png"
	"net/http"

	"github.com/google/uuid"
	"github.com/nfnt/resize"
	"github.com/vogtp/som/pkg/core"
	"github.com/vogtp/som/pkg/core/cfg"
	"github.com/vogtp/som/pkg/core/mime"
	"github.com/vogtp/som/pkg/core/status"
)

type statusData struct {
	status.Status
	Image   map[string]template.HTML
	OK      status.Level
	Baseurl string
}

func (s *WebStatus) handleTopology(w http.ResponseWriter, r *http.Request) {

	var data = struct {
		*commonData
		TimeFormat string
		Status     *statusData
	}{
		commonData: common("SOM Topology", r),
		TimeFormat: cfg.TimeFormatString,
		Status:     prepaireStatus(s.data.Status),
	}
	err := templates.ExecuteTemplate(w, "topology.gohtml", data)
	if err != nil {
		s.hcl.Errorf("Topology Template error %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func prepaireStatus(stat status.Status) *statusData {
	data := &statusData{
		Status:  stat,
		Image:   make(map[string]template.HTML),
		OK:      status.OK,
		Baseurl: core.Get().WebServer().BasePath(),
	}
	for _, sz := range data.Status.Szenarios() {
		for _, reg := range sz.Regions() {
			for _, u := range reg.Users() {
				if u.Level() <= status.OK {
					continue
				}
				evt := u.LastEvent()
				if evt == nil || evt.Err() == nil {
					continue
				}
				for _, f := range evt.Files {
					if f.Type == mime.Png {
						img, err := getImage(f.Payload)
						if err != nil {
							data.Image[evt.ID.String()] = template.HTML(err.Error())
						}
						data.Image[evt.ID.String()] = template.HTML(img)
						break
					}
				}
			}
		}
	}
	return data
}

func getImage(img []byte) (string, error) {
	image, err := png.Decode(bytes.NewReader(img))
	if err != nil {
		return "", fmt.Errorf("cannot decode image: %w", err)
	}

	newImage := resize.Resize(600, 0, image, resize.Lanczos3)

	buf := new(bytes.Buffer)
	enc := png.Encoder{CompressionLevel: png.BestCompression}
	err = enc.Encode(buf, newImage)
	if err != nil {
		return "", fmt.Errorf("cannot encode image to png: %w", err)
	}

	imgb64 := base64.StdEncoding.EncodeToString(buf.Bytes())
	id := uuid.NewString()
	return fmt.Sprintf("<img id='%s' onclick='resize(\"%s\")' width='100'' src='data:image/png;base64, %s' /><br><small>Click to resize...</small>", id, id, imgb64), nil
}
