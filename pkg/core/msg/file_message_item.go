package msg

import (
	"crypto/md5"

	"github.com/google/uuid"
	"github.com/iancoleman/strcase"
	"github.com/vogtp/go-hcl"
	"github.com/vogtp/som/pkg/core/mime"
)

// FileMsgItem contains a file
type FileMsgItem struct {
	ID      uuid.UUID `json:"ID"  gorm:"primaryKey;type:uuid"`
	Name    string    `json:"Name"`
	Type    mime.Type
	Size    int `json:"Size"`
	Payload []byte
}

// NewFileMsgItem creates a new FileMsgItem
func NewFileMsgItem(name string, mtype mime.Type, payload []byte) *FileMsgItem {
	fmi := &FileMsgItem{
		Name:    strcase.ToLowerCamel(name),
		Type:    mtype,
		Size:    len(payload),
		Payload: payload,
	}
	fmi.CalculateID()
	return fmi
}

// CalculateID calculates the uuid from the file hash if needed
func (fmi *FileMsgItem) CalculateID() {
	fmi.ID = idFromMD5(fmi.Payload)
}

func idFromMD5(d []byte) uuid.UUID {
	h := make([]byte, md5.Size)
	for i, b := range md5.Sum(d) {
		h[i] = b
	}
	id, err := uuid.FromBytes(h)
	if err != nil {
		// should never been reacht md5 and uuid both have a size of 16
		hcl.Warn("Cannot create uuid from MD5", "error", err, "uuid_bytes", string(h))
	}
	if id == uuid.Nil {
		hcl.Warn("Generating non hash based id")
		id = uuid.New()
	}
	return id
}
