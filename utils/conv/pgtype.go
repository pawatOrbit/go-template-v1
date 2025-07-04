package conv

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func PgTypeTextToString(t pgtype.Text) string {
	if !t.Valid {
		return ""
	}
	return t.String
}

func PgTypeTextToStringPointer(t pgtype.Text) *string {
	if !t.Valid {
		return nil
	}
	return &t.String
}

func PgTypeBoolToBool(b pgtype.Bool) bool {
	if !b.Valid {
		return false
	}
	return b.Bool
}

func PgTypeBoolToBoolPointer(b pgtype.Bool) *bool {
	if !b.Valid {
		return nil
	}
	return &b.Bool
}

func PgTypeInt4ToInt32(i pgtype.Int4) int32 {
	if !i.Valid {
		return 0
	}
	return i.Int32
}

func PgTypeInt4ToIntPointer(i pgtype.Int4) *int32 {
	if !i.Valid {
		return nil
	}
	return &i.Int32
}

func PgTypeInt8ToInt64(i pgtype.Int8) int64 {
	if !i.Valid {
		return 0
	}
	return i.Int64
}

func PgTypeInt8ToInt64Pointer(i pgtype.Int8) *int64 {
	if !i.Valid {
		return nil
	}
	return &i.Int64
}

func PgTypeTimestamptzToTimePointer(pt pgtype.Timestamptz) *time.Time {
	if !pt.Valid {
		return nil
	}

	t := pt.Time
	return &t
}

func PgtypeUUIDToString(pgUUID pgtype.UUID) string {
	var u uuid.UUID
	err := u.UnmarshalBinary(pgUUID.Bytes[:])
	if err != nil {
		return ""
	}
	return u.String()
}

func PgtypeUUIDToStringPointer(pgUUID pgtype.UUID) *string {
	var u uuid.UUID
	err := u.UnmarshalBinary(pgUUID.Bytes[:])
	if err != nil {
		return nil
	}
	s := u.String()
	return &s
}

func PgtypeUUIDToUUID(pgUUID pgtype.UUID) uuid.UUID {
	var u uuid.UUID
	err := u.UnmarshalBinary(pgUUID.Bytes[:])
	if err != nil {
		return uuid.Nil
	}
	return u
}

func PgtypeUUIDToUUIDNull(pgUUID pgtype.UUID) uuid.NullUUID {
	var u uuid.NullUUID
	err := u.UnmarshalBinary(pgUUID.Bytes[:])
	if err != nil {
		return uuid.NullUUID{}
	}

	return u
}
