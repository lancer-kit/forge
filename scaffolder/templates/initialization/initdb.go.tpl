package initialization

import (
	"github.com/lancer-kit/armory/db"
	"github.com/sirupsen/logrus"

	"{{.project_name}}/config"
)

var (
	DB   initModule = "database connection"
)

func initDatabase(cfg *config.Cfg, entry *logrus.Entry) error {
	return db.Init(cfg.DB.ConnURL, entry)
}
