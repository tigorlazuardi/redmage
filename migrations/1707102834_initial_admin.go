package migrations

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models"
)

func init() {
	_ = godotenv.Load()
	m.Register(func(db dbx.Builder) error {
		email := os.Getenv("REDMAGE_INIT_EMAIL")
		password := os.Getenv("REDMAGE_INIT_PASSWORD")

		if email == "" || password == "" {
			return nil
		}

		admin := &models.Admin{Email: email}
		if err := admin.SetPassword(password); err != nil {
			return fmt.Errorf("failed to set initial admin password: %w", err)
		}
		// add up queries...
		dao := daos.New(db)
		return dao.SaveAdmin(admin)
	}, func(db dbx.Builder) error {
		_ = godotenv.Load()
		email := os.Getenv("REDMAGE_INIT_EMAIL")

		dao := daos.New(db)

		admin, err := dao.FindAdminByEmail(email)
		if err == nil {
			return dao.DeleteAdmin(admin)
		}

		// already deleted
		return nil
	})
}
