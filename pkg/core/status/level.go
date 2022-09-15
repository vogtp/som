//go:generate stringer -type=Level

package status

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

// Level is the serverity level of a state
type Level int

const (
	// Unknown indicates that no information is availabler
	Unknown Level = iota
	// OK indicates OK (or if an event at least the last event was OK)
	OK
	// Issues indicates that some states where not OK
	Issues
	// Warning indicates that most states where not OK
	Warning
	// Down indicates that all (or at least a lot) states where DOWN
	Down
)

// Img returns the name of the image
func (l Level) Img() string {
	switch l {
	case Unknown:
		return "darkgray"
	case OK:
		return "green"
	case Issues:
		return "yellow"
	case Warning:
		return "orange"
	case Down:
		return "red"
	default:
		return "blueQuestion"
	}
}

// FromString converts string to Level or panics
func (Level) FromString(lvl string) Level {
	lvl = strings.ToLower(lvl)
	for i := Unknown; i <= Down; i++ {
		if strings.ToLower(i.String()) == lvl {
			return i
		}
	}
	return Unknown
}

// UnmarshalGQL for graphql
func (e *Level) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Unknown.FromString(str)
	return nil
}

// MarshalGQL for graphql
func (e Level) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
