package carta

import "github.com/jmoiron/sqlx"

// Mapx maps sqlx.Rows onto a struct or slice of structs.
// It is a convenience wrapper around the Map function.
func Mapx(rows *sqlx.Rows, dst interface{}) error {
	return Map(rows.Rows, dst)
}
