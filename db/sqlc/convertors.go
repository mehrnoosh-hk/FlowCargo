package sqlc

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// TextToString converts a pointer to pgtype.Text to a string.
// It returns an empty string if the pointer is nil or the pgtype.Text is not valid.
func TextToString(t *pgtype.Text) string {
	if t == nil || !t.Valid {
		return ""
	}
	return t.String
}

// StringToText converts a pointer to a string to pgtype.Text.
// A nil pointer is converted to an invalid (NULL) pgtype.Text.
func StringToText(s *string) pgtype.Text {
	if s == nil {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: *s, Valid: true}
}

// BoolToBool converts a pointer to pgtype.Bool to a bool.
// It returns false if the pointer is nil or the pgtype.Bool is not valid.
func BoolToBool(b *pgtype.Bool) bool {
	if b == nil || !b.Valid {
		return false
	}
	return b.Bool
}

// BoolToPgtypeBool converts a pointer to a bool to pgtype.Bool.
// A nil pointer is converted to an invalid (NULL) pgtype.Bool.
func BoolToPgtypeBool(b *bool) pgtype.Bool {
	if b == nil {
		return pgtype.Bool{Valid: false}
	}
	return pgtype.Bool{Bool: *b, Valid: true}
}

// UUIDToUUID converts a pointer to pgtype.UUID to a uuid.UUID.
// It returns a zero UUID if the pointer is nil or the pgtype.UUID is not valid.
func UUIDToUUID(u *pgtype.UUID) uuid.UUID {
	if u == nil || !u.Valid {
		return uuid.UUID{}
	}
	return u.Bytes
}

// UUIDToPgtypeUUID converts a pointer to a uuid.UUID to pgtype.UUID.
// A nil pointer is converted to an invalid (NULL) pgtype.UUID.
func UUIDToPgtypeUUID(u *uuid.UUID) pgtype.UUID {
	if u == nil {
		return pgtype.UUID{Valid: false}
	}
	return pgtype.UUID{Bytes: *u, Valid: true}
}