package busbridge

import (
	"context"
	"fmt"
	"reflect"

	"log/slog"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
	"github.com/vogtp/som/pkg/core/log"
	"github.com/vogtp/som/pkg/core/msg"
)

// Slack bridges the eventbus over slack
type Slack struct {
	log     *slog.Logger
	ctx     context.Context
	client  *socketmode.Client
	auth    *slack.AuthTestResponse
	channel string
}

// NewSlack creates a new slack event bridge
func NewSlack(ctx context.Context, logger *slog.Logger, oAuth string, appLevelToken string) (*Slack, error) {
	logger = logger.With(log.Component, "slack")
	client := socketmode.New(
		slack.New(
			oAuth,
			slack.OptionAppLevelToken(appLevelToken),
			slack.OptionDebug(false),
		),
		socketmode.OptionDebug(false),
		//socketmode.OptionLog(log),
	)
	auth, err := client.AuthTestContext(ctx)
	if err != nil {
		logger.Error("Cannot login", log.Error, err)
		return nil, err
	}
	logger.Info("connected to slack", log.User, auth.User, "user_id", auth.UserID, "bot_id", auth.BotID)
	s := &Slack{
		log:     logger,
		client:  client,
		auth:    auth,
		channel: "#dev_event",
		ctx:     ctx,
	}
	go s.start()
	// _, err = client.CreateConversation(s.channel[1:], false)
	// if err != nil && err.Error() != "name_taken" {
	// 	panic(err)
	// }

	// users, cursor, err := client.GetUsersInConversation(&slack.GetUsersInConversationParameters{
	// 	ChannelID: s.channel,
	// })
	// if err != nil {
	// 	s.log.Errorf("channel %s probably not found: %v", s.channel, err)
	// 	os.Exit(1)
	// }
	// if len(cursor) > 0 {
	// 	s.log.Warnf("users cursor: %v should get more", cursor)
	// }
	// for _, u := range users {
	// 	s.log.Infof("user: %v", u)
	// }
	s.log.Info("Started bus handler")
	return s, nil
}

// start starts the slack connection
func (s Slack) start() {
	go s.handleEvents()
	//s.say("listening") // TODO should say something more meaningfull here

	if err := s.client.Run(); err != nil {
		s.log.Error("could not connect to slack", log.Error, err)
	}
}

// Send messge to the bridge
func (s Slack) Send(evt *msg.SzenarioEvtMsg) error {
	s.log.Warn("Sending message", "type", reflect.TypeOf(evt))
	b, err := ToBridgePayload(evt)
	if err != nil {
		s.log.Warn("bridge payload", log.Error, err)
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
	s.log.Info("Handling slack events", "channel", s.channel)
	for {
		select {
		case <-s.ctx.Done():
			s.log.Warn("Shutting down socketmode listener", "err", s.ctx.Err())
			return
		case event := <-s.client.Events:
			s.log.Debug("Slack Event", "event_type", event.Type, "type", reflect.TypeOf(event.Type), "evnt", event)

			switch event.Type {
			case socketmode.EventTypeEventsAPI:
				eventsAPIEvent, ok := event.Data.(slackevents.EventsAPIEvent)
				if !ok {
					s.log.Warn("Could not type cast the event to the EventsAPIEvent: %v\n", "event", event)
					continue
				}
				// We need to send an Acknowledge to the slack server
				s.client.Ack(*event.Request)

				err := s.handleEventMessage(s.ctx, eventsAPIEvent)
				if err != nil {
					s.log.Warn("handle message error", log.Error, err)
				}
			case socketmode.EventTypeSlashCommand:
				cmdEvt, ok := event.Data.(slack.SlashCommand)
				if !ok {
					s.log.Warn("Could not cast the command event")
					continue
				}
				s.client.Ack(*event.Request)
				s.log.Info("socket mode command", "command", cmdEvt)
				//go s.handleCommand(ctx, &cmdEvt)
			default:
				s.log.Info("Unsupported eventtype", "type", event.Type, "data", event.Data)
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
			s.log.Debug("Mention", "text", ev.Text)
			// if err := s.client.AddReactionContext(ctx, "eyes", slack.NewRefToMessage(ev.Channel, ev.TimeStamp)); err != nil {
			// 	s.log.Infof("Adding readtion: %v", err)
			// }
			// go handleBackendQuery(ctx, ev.Text, client, ev.Channel)
			// if err := s.client.AddReactionContext(ctx, "white_check_mark", slack.NewRefToMessage(ev.Channel, ev.TimeStamp)); err != nil {
			// 	s.log.Infof("Adding readtion: %v", err)
			// }
		case *slackevents.MessageEvent:
			if ev.BotID == s.auth.BotID {
				s.log.Error("There is a BOT with the same ID probably missing messages")
				//FIXME send a alert message
			}
			busMsg, err := FromBridgePayload(ev.Text)
			if err != nil {
				return fmt.Errorf("could not process msg: %w %s", err, ev.Text)
			}
			s.log.Info("GOT MESSAGE", "name", busMsg.Name, "message", busMsg)
			// go handleGeneralMessage(ctx, ev, client)

		default:
			return fmt.Errorf("unsupported event type: %v", event.Type)
		}
	default:
		return fmt.Errorf("unsupported event type: %v", event.Type)
	}
	return nil
}
