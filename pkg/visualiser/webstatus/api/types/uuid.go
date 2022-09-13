package types

import (
	"fmt"

	"github.com/99designs/gqlgen/graphql"
	"github.com/google/uuid"
)

// MarshalUUID marshals uuids
func MarshalUUID(id uuid.UUID) graphql.Marshaler {
	return graphql.MarshalString(id.String())
}

// UnmarshalUUID unmarshals uuids
func UnmarshalUUID(v interface{}) (uuid.UUID, error) {
	switch v := v.(type) {
	case string:
		return uuid.Parse(v)
	case []byte:
		return uuid.ParseBytes(v)
	case uuid.UUID:
		return v, nil
	case nil:
		return uuid.Nil, nil
	default:
		return uuid.Nil, fmt.Errorf("%T is not a UUID", v)
	}
}
