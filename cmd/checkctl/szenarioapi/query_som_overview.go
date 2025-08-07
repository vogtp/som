package szenarioapi

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/spf13/viper"
	"github.com/vogtp/go-icinga/pkg/check"
	"github.com/vogtp/go-icinga/pkg/icinga"
	"github.com/vogtp/som/pkg/visualiser/webstatus"
)

func querySomAPI(ctx context.Context, result *check.Result) error {
	url := viper.GetString(somURL)
	slog := slog.Default().With("som.url", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("create request for %s: %w", url, err)
	}
	req.Header.Set("Accept", "application/json")
	req = req.WithContext(ctx)

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("query %s: %w", url, err)
	}
	slog.Debug("SOM API resonse", "resp", resp)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("get http status %v (%s)", resp.StatusCode, resp.Status)
	}
	apiResp := webstatus.OverviewAPIResonse{}
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return fmt.Errorf("decoding json response from %s: %w", url, err)
	}
	szName := viper.GetString(szNameFlag)
	for _, sz := range apiResp.Szenarios {
		if !strings.EqualFold(sz.Name, szName) {
			continue
		}
		result.SetCode(imgcolor2Resultcode(sz.Img))
		slog.Debug("Szenaro response", "name", sz.Name, "data", sz)
		last := timeFormater("total", check.Data{Value: parseTime(string(sz.LastTime))})
		avg := timeFormater("total", check.Data{Value: parseTime(string(sz.AvgTime))})
		iList := fmt.Sprintf(`%s%s`, url, sz.IncidentList)
		result.WriteHeader("Duration %s", last)
		result.WriteHeader("Incident List: <a href='%s'>%s</a>", iList, "Link")
		result.SetHeader(`%s\n\n<br>Incident List: <a href="%s">%s</a>`, fmt.Sprintf("Duration %s", last), iList, iList)
		result.SetCounter("Response (current)", last)
		result.SetCounter("Response (Average)", avg)
		result.SetStatus("Availability (current)", sz.AvailabilityCur)
		result.SetStatus("Availability (average)", sz.AvailabilityAvg)
		result.SetStatus("Status", strings.TrimSpace(sz.Status))
		result.SetStatus("Incidents", sz.IncidentCount)
	}
	return nil
}

func imgcolor2Resultcode(l string) icinga.ResultCode {
	switch l {
	case "darkgray": //Unknown:
		return icinga.UNKNOWN
	case "green": //OK:
		return icinga.OK
	case "yellow": // Issues:
		return icinga.WARNING
	case "orange": // Warning:
		return icinga.WARNING
	case "red": //Down:
		return icinga.CRITICAL
	default:
		return icinga.UNKNOWN
	}
}

func parseTime(s string) time.Duration {
	if d, err := time.ParseDuration(s); err == nil {
		return d
	}
	idx := strings.Index(s, ">")
	if idx > 0 {
		s = s[idx+1 : len(s)-6]
		if d, err := time.ParseDuration(s); err == nil {
			return d
		}
	}
	return 0
}
