package ent

import (
	"github.com/vogtp/som/pkg/core/mime"
	"github.com/vogtp/som/pkg/core/msg"
)

// MsgItem converts to msg.FileMsgItem
func (f File) MsgItem() msg.FileMsgItem {
	return msg.FileMsgItem{
		ID:   f.UUID,
		Name: f.Name,
		Type: mime.Type{
			MimeType: f.Type,
			Ext:      f.Ext,
		},
		Size: f.Size,
		//Payload: f.Payload,
	}
}
