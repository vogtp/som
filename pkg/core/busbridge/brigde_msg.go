package busbridge

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/vogtp/som/pkg/core/msg"
)

// BrigeMsg is used to encapsulate eventbus messages when sending over a bridge
type BrigeMsg struct {
	Type        string
	Src         uuid.UUID
	JSONPayload []byte
}

func getType(o any) string {
	return fmt.Sprintf("%T", o)
}

// ToBridgePayload converts to bridge payload
func ToBridgePayload(msg *msg.SzenarioEvtMsg) (string, error) {
	b, err := json.Marshal(msg)
	if err != nil {
		return "", fmt.Errorf("marshalling payload: %w", err)
	}
	bm := &BrigeMsg{
		Type:        getType(msg),
		JSONPayload: b,
	}
	m, err := json.Marshal(bm)
	if err != nil {
		return "", fmt.Errorf("marshalling brige message: %w", err)
	}
	return string(m), nil
}

// FromBridgePayload converts from bridge payload
func FromBridgePayload(payload string) (rmsg *msg.SzenarioEvtMsg, err error) {
	defer func() {
		// handle json errors
		if r := recover(); r != nil {
			err = fmt.Errorf("json panic: %v", r)
		}
	}()
	bm := &BrigeMsg{}
	err = json.Unmarshal([]byte(payload), bm)
	if err != nil {
		return nil, err
	}

	// if bm.Type != getType(msg) {
	// 	return nil, fmt.Errorf("wrong type: got %v expecting %v", bm.Type, getType(msg))
	// }
	msg, err := msg.SzenarioEvtMsgFromJSON(bm.JSONPayload)
	return msg, err
}
