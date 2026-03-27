package diregistry

import (
	"daily-hello-service/config"
	"go-libs/dihelper"

	"github.com/sarulabs/di"
)

func initConfigBuilder() {
	dihelper.ConfigsBuilder = func() []di.Def {
		arr := []di.Def{}
		arr = append(arr, di.Def{
			Name:  ConfigDIName,
			Scope: di.App,
			Build: func(ctn di.Container) (any, error) {
				cfg, err := config.Load()
				return cfg, err
			},
			Close: func(obj any) error {
				return nil
			},
		})

		return arr
	}
}
