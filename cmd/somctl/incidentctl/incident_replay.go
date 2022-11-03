//nolint:all
package incidentctl

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/vogtp/go-hcl"
	"github.com/vogtp/som/cmd/somctl/term"
	"github.com/vogtp/som/pkg/core"
	"github.com/vogtp/som/pkg/core/cfg"
	"github.com/vogtp/som/pkg/core/mime"
	"github.com/vogtp/som/pkg/core/msg"
	"github.com/vogtp/som/pkg/visualiser/webstatus/db"
	"github.com/vogtp/som/pkg/visualiser/webstatus/db/ent"
	"github.com/vogtp/som/pkg/visualiser/webstatus/db/ent/incident"
)

const (
	flagReplayDelay = "replay.delay"
)

func init() {
	pflag.Duration(flagReplayDelay, time.Second, "Delay between events in replay")
}

var incidentReplay = &cobra.Command{
	Use:   "replay <Incident UUID>",
	Short: "replay all events from an incident",
	Long:  ``,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if !hcl.IsGoRun() {
			fmt.Println("Incident replay can be used to replay incidents to debug")
			ans := term.Read("DO YOU REALLY WANT TO CONTINE? (yes|no)")
			if ans != "yes" {
				fmt.Printf("%s? OK! Stopping!\n", ans)
				return nil
			}
		}
		incidenId, err := uuid.Parse(args[0])
		if err != nil {
			return fmt.Errorf("%s is not a UUID: %w", args[0], err)
		}
		szBus := core.Get().Bus().Szenario
		client, err := db.New()
		if err != nil {
			return fmt.Errorf("Cannot open DB: %w", err)
		}
		incidents, err := client.Incident.Query().Where(incident.IncidentID(incidenId)).All(cmd.Context())
		if err != nil {
			return fmt.Errorf("Cannot query incident %s: %w", incidenId, err)
		}
		if len(incidents) < 1 {
			return fmt.Errorf("No such incident found")
		}
		replayDelay := viper.GetDuration(flagReplayDelay)
		estDuration := time.Duration(len(incidents)) * replayDelay
		fmt.Printf("\nReplaying %d events\nName: %s\nStart: %s\nEstimated time: %v (delay: %v)\nYou will need to copy alertmgr.json from stater to get sensible results!\n", len(incidents), incidents[0].Name, incidents[0].Start.Format(cfg.TimeFormatString), estDuration, replayDelay)
		ans := term.Read("Replay? (y|n)")
		if ans != "y" {
			return fmt.Errorf("You answered %s, stopping", ans)
		}

		ticker := time.NewTicker(replayDelay)
		defer fmt.Println("")
		bar := progressbar.Default(int64(len(incidents)), "Replay")
		for _, inci := range incidents {
			bar.Add(1)
			m, err := getMsg(cmd.Context(), inci)
			if err != nil {
				return fmt.Errorf("cannot build msg from model: %v", err)
			}
			szBus.Send(m)
			select {
			case <-ticker.C:
			case <-cmd.Context().Done():
				fmt.Println("Canceled")
				return cmd.Context().Err()
			}
		}
		fmt.Printf("Finished...\n")
		return nil
	},
}

func getMsg(ctx context.Context, m *ent.Incident) (*msg.SzenarioEvtMsg, error) {
	if m == nil {
		return nil, fmt.Errorf("Got nil message: %v", m)
	}
	e := &msg.SzenarioEvtMsg{
		ID:         m.UUID,
		IncidentID: m.IncidentID.String(),
		Name:       m.Name,
		Time:       m.Time,
		Username:   m.Username,
		Region:     m.Region,
		ProbeOS:    m.ProbeOS,
		ProbeHost:  m.ProbeHost,
		Counters:   make(map[string]float64),
		Stati:      make(map[string]string),
	}
	failures, err := m.QueryFailures().All(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting failures (errors): %v", err)
	}
	e.Errors = make([]string, len(failures))
	for i, f := range failures {
		e.Errors[i] = f.Error
	}
	cntr, err := m.QueryCounters().All(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting counters: %v", err)
	}
	for _, c := range cntr {
		e.Counters[c.Name] = c.Value
	}
	stati, err := m.QueryStati().All(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting Stati: %v", err)
	}
	for _, s := range stati {
		e.Stati[s.Name] = s.Value
	}
	files, err := m.QueryFiles().All(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting Files: %v", err)
	}
	e.Files = make([]msg.FileMsgItem, len(files))
	for i, f := range files {
		e.Files[i] = msg.FileMsgItem{
			ID:   f.UUID,
			Name: f.Name,
			Type: mime.Type{
				MimeType: f.Type,
				Ext:      f.Ext,
			},
			Size:    f.Size,
			Payload: f.Payload,
		}
	}
	return e, nil
}
