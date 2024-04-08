package api

import (
	"database/sql"

	"github.com/tigorlazuardi/redmage/db/queries"
)

type API struct {
	Queries *queries.Queries
	DB      *sql.DB
}
