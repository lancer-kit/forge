package workers

import (
	"github.com/lancer-kit/uwe/v2"
	"github.com/sirupsen/logrus"

	"{{.project_name}}/config"
	{{if .api}}"{{.project_name}}/workers/api"{{end}}
	{{.project_name}}/workers/foobar"
)

func InitChief(logger *logrus.Entry, cfg *config.Cfg) uwe.Chief {
	chief := uwe.NewChief()
	chief.UseDefaultRecover()
	chief.SetEventHandler(func(event uwe.Event) {
		var level logrus.Level
		switch event.Level {
		case uwe.LvlFatal, uwe.LvlError:
			level = logrus.ErrorLevel
		case uwe.LvlInfo:
			level = logrus.InfoLevel
		default:
			level = logrus.WarnLevel
		}

		logger.WithFields(event.Fields).
			Log(level, event.Message)
	})

    {{if or .api .simple_worker}}logger = logger.WithField("app_layer", "workers"){{end}}

    {{if .api}}chief.AddWorker(config.WorkerAPIServer, api.GetServer(cfg, logger.WithField("worker", config.WorkerFooBar))){{end}}
	chief.AddWorker(config.WorkerFooBar, foobar.NewWorker(config.WorkerFooBar, logger.WithField("worker", config.WorkerFooBar)))
	return chief
}
