package busbridge

import (
	"context"
	"fmt"
	"reflect"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
	"github.com/vogtp/go-hcl"
	"github.com/vogtp/som/pkg/core/msg"
)

// Slack bridges the eventbus over slack
type Slack struct {
	hcl     hcl.Logger
	ctx     context.Context
	client  *socketmode.Client
	auth    *slack.AuthTestResponse
	channel string
}

// NewSlack creates a new slack event bridge
func NewSlack(ctx context.Context, hcl hcl.Logger, oAuth string, appLevelToken string) (*Slack, error) {
	hcl = hcl.Named("slack")
	client := socketmode.New(
		slack.New(
			oAuth,
			slack.OptionAppLevelToken(appLevelToken),
			slack.OptionDebug(false),
		),
		socketmode.OptionDebug(false),
		//socketmode.OptionLog(hcl),
	)
	auth, err := client.AuthTestContext(ctx)
	if err != nil {
		hcl.Errorf("Cannot login: %v", err)
		return nil, err
	}
	hcl.Infof("connected to slack as: %s (%s) - botID: %s", auth.User, auth.UserID, auth.BotID)
	s := &Slack{
		hcl:     hcl,
		client:  client,
		auth:    auth,
		channel: "#dev_event",
		ctx:     ctx,
	}
	go s.start()
	_, err = client.CreateConversation(s.channel[1:], false)
	if err != nil && err.Error() != "name_taken" {
		panic(err)
	}

	// users, cursor, err := client.GetUsersInConversation(&slack.GetUsersInConversationParameters{
	// 	ChannelID: s.channel,
	// })
	// if err != nil {
	// 	s.hcl.Errorf("channel %s probably not found: %v", s.channel, err)
	// 	os.Exit(1)
	// }
	// if len(cursor) > 0 {
	// 	s.hcl.Warnf("users cursor: %v should get more", cursor)
	// }
	// for _, u := range users {
	// 	s.hcl.Infof("user: %v", u)
	// }
	s.hcl.Info("Started bus handler")
	return s, nil
}

// start starts the slack connection
func (s Slack) start() {
	go s.handleEvents()
	//s.say("listening") // TODO should say something more meaningfull here

	if err := s.client.Run(); err != nil {
		s.hcl.Errorf("could not connect to slack: %v", err)
	}
}

// Send messge to the bridge
func (s Slack) Send(evt *msg.SzenarioEvtMsg) error {
	s.hcl.Warnf("Sending message: %v", reflect.TypeOf(evt))
	b, err := ToBridgePayload(evt)
	if err != nil {
		s.hcl.Warnf("bridge payload: %v", err)
	}
	chanID, ts, err := s.client.PostMessage(
		s.channel,
		slack.MsgOptionText(b, true),
	)
	if err != nil {
		return fmt.Errorf("send %v %v: %w", ts, chanID, err)
	}
	return nil
}

//nolint:unused
func (s Slack) say(txt string) error {
	chanID, ts, err := s.client.PostMessage(
		s.channel,
		slack.MsgOptionText(txt, false),
	)
	if err != nil {
		return fmt.Errorf("say %v %v: %w", ts, chanID, err)
	}
	return nil
}

func (s Slack) handleEvents() {
	s.hcl.Infof("Handling slack events in %s", s.channel)
	for {
		select {
		case <-s.ctx.Done():
			s.hcl.Warn("Shutting down socketmode listener: %v", s.ctx.Err())
			return
		case event := <-s.client.Events:
			s.hcl.Tracef("Slack Event %v(%T): %v", event.Type, event.Type, event)

			switch event.Type {
			case socketmode.EventTypeEventsAPI:
				eventsAPIEvent, ok := event.Data.(slackevents.EventsAPIEvent)
				if !ok {
					s.hcl.Warn("Could not type cast the event to the EventsAPIEvent: %v\n", event)
					continue
				}
				// We need to send an Acknowledge to the slack server
				s.client.Ack(*event.Request)

				err := s.handleEventMessage(s.ctx, eventsAPIEvent)
				if err != nil {
					s.hcl.Warnf("handle message: %v", err)
				}
			case socketmode.EventTypeSlashCommand:
				cmdEvt, ok := event.Data.(slack.SlashCommand)
				if !ok {
					s.hcl.Warn("Could not cast the command event")
					continue
				}
				s.client.Ack(*event.Request)
				s.hcl.Infof("Command: %v", cmdEvt)
				//go s.handleCommand(ctx, &cmdEvt)
			default:
				s.hcl.Infof("Unsupported eventtype %v: %T", event.Type, event.Data)
			}

		}
	}
}

func (s Slack) handleEventMessage(ctx context.Context, event slackevents.EventsAPIEvent) error {
	switch event.Type {
	case slackevents.CallbackEvent:
		innerEvent := event.InnerEvent
		switch ev := innerEvent.Data.(type) {
		case *slackevents.AppMentionEvent:
			s.hcl.Debugf("Mention: %v", ev.Text)
			// if err := s.client.AddReactionContext(ctx, "eyes", slack.NewRefToMessage(ev.Channel, ev.TimeStamp)); err != nil {
			// 	s.hcl.Infof("Adding readtion: %v", err)
			// }
			// go handleBackendQuery(ctx, ev.Text, client, ev.Channel)
			// if err := s.client.AddReactionContext(ctx, "white_check_mark", slack.NewRefToMessage(ev.Channel, ev.TimeStamp)); err != nil {
			// 	s.hcl.Infof("Adding readtion: %v", err)
			// }
		case *slackevents.MessageEvent:
			if ev.BotID == s.auth.BotID {
				s.hcl.Errorf("There is a BOT with the same ID probably missing messages")
				//FIXME send a alert message
			}
			busMsg, err := FromBridgePayload(ev.Text)
			if err != nil {
				return fmt.Errorf("could not process msg: %w %s", err, ev.Text)
			}
			s.hcl.Infof("GOT MESSAGE: %v %+v ", busMsg.Name, busMsg)
			// go handleGeneralMessage(ctx, ev, client)

		default:
			return fmt.Errorf("unsupported event type: %v", event.Type)
		}
	default:
		return fmt.Errorf("unsupported event type: %v", event.Type)
	}
	return nil
}
