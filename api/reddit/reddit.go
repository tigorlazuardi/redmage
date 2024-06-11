package reddit

import (
	"github.com/tigorlazuardi/redmage/config"
)

type Reddit struct {
	Client Client
	Config *config.Config
}
