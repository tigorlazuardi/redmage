package migrations

import (
	"os"
	"strconv"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		dao := daos.New(db)
		settings, err := dao.FindSettings()
		if err != nil {
			return err
		}
		logMaxDays, err := strconv.Atoi(os.Getenv("REDMAGE_INIT_LOG_MAX_DAYS"))
		if err != nil {
			logMaxDays = 7
		}
		settings.Meta.AppName = "Redmage"
		settings.Meta.AppUrl = os.Getenv("REDMAGE_INIT_META_APP_URL")
		settings.Logs.MaxDays = logMaxDays
		return nil
	}, nil)
}
