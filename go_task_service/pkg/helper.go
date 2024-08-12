package pkg

import (
	"database/sql"
)

func NullStringToString(s sql.NullString) string {
	if s.Valid {
		return s.String
	}

	return ""
}

func NullFloatToFloat(s sql.NullFloat64) float64 {
	if s.Valid {
		return s.Float64
	}

	return 0
}

func NullTimeToString(s sql.NullTime) string {
	if s.Valid {
		return s.Time.String()
	}

	return ""
}
