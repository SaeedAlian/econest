package json_types

import (
	"database/sql"
	"encoding/json"
)

type JSONNullString struct {
	sql.NullString
}

func (ns JSONNullString) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ns.String)
}

type JSONNullTime struct {
	sql.NullTime
}

func (nt JSONNullTime) MarshalJSON() ([]byte, error) {
	if !nt.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(nt.Time)
}

type JSONNullInt32 struct {
	sql.NullInt32
}

func (nt JSONNullInt32) MarshalJSON() ([]byte, error) {
	if !nt.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(nt.Int32)
}

type JSONNullInt64 struct {
	sql.NullInt64
}

func (nt JSONNullInt64) MarshalJSON() ([]byte, error) {
	if !nt.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(nt.Int64)
}
