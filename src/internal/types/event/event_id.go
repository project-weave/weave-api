package event

import (
	"database/sql/driver"
	"encoding/base64"
	"fmt"

	"github.com/google/uuid"
)

// Users will directly use EventIDs via URLs, we want the IDs to be short, however how we represent data on the backend doesnt matter
type EventUUID struct {
	uuid.UUID
}

func ToEventUUID(idStr string) (EventUUID, error) {
	// Attempt to parse the input string directly as a UUID
	parsedUUID, parseErr := uuid.Parse(idStr)
	if parseErr == nil {
		return EventUUID{parsedUUID}, nil
	}

	// If parsing fails, assume the input might be a Base64 encoded UUID
	parsedBase64ID, err := base64.RawURLEncoding.DecodeString(idStr)
	if err != nil {
		return EventUUID{}, fmt.Errorf("failed to decode Base64 UUID string: %w", err)
	}

	parsedUUID, err = uuid.FromBytes(parsedBase64ID)
	if err != nil {
		return EventUUID{}, fmt.Errorf("failed to parse UUID from Base64 bytes: %w", err)
	}

	out := EventUUID{parsedUUID}
	fmt.Printf("%s -> %s\n", idStr, out.UUID)
	return out, nil
}

func (eid *EventUUID) String() string {
	base64ID := base64.RawURLEncoding.EncodeToString(eid.UUID[:])
	return base64ID
}

// Convert Base64-encoded UUID string or UUID string to EventUUID
func (eid *EventUUID) UnmarshalText(data []byte) error {
	id, err := ToEventUUID(string(data))
	if err != nil {
		return err
	}
	*eid = id
	return nil
}

// Output Base64-encoded UUID
func (eid EventUUID) MarshalText() ([]byte, error) {
	return []byte(eid.String()), nil
}

// Output UUID string
func (eid EventUUID) Value() (driver.Value, error) {
	return eid.UUID.String(), nil
}

// Convert UUID string to EventUUID
func (eid *EventUUID) Scan(value interface{}) error {
	uuidStr, ok := value.(string)
	if !ok {
		return fmt.Errorf("cannot convert %v to string", value)
	}
	id, err := ToEventUUID(uuidStr)
	if err != nil {
		return err
	}

	*eid = id
	return nil
}
