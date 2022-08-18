package mime

// Type represents a mime type
type Type struct {
	MimeType string
	Ext      string
}

var (
	// Png is a png image
	Png = Type{MimeType: "image/png", Ext: "png"}
	// Jpg is a jpeg image
	//	Jpg = Type{MimeType: "image/jpeg", Ext: "jpeg"}
	// HTML is a html page
	HTML = Type{MimeType: "text/html", Ext: "html"}
)
