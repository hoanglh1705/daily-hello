package diregistry

import (
	"daily-hello-service/config"
	"go-libs/contexthelper"
	"go-libs/dihelper"
	auth_middleware "go-libs/http_middlewares/auth"

	"github.com/sarulabs/di"
)

func initAdaptersBuilder() {
	dihelper.AdaptersBuilder = func() []di.Def {
		arr := []di.Def{}
		arr = append(arr,
			di.Def{
				Name:  JWTMiddlewareDIName,
				Scope: di.App,
				Build: func(ctn di.Container) (any, error) {
					cfg := ctn.Get(ConfigDIName).(*config.Config)
					return auth_middleware.New(cfg.JwtConfig.Algorithm, cfg.JwtConfig.SecretKey, cfg.JwtConfig.Duration, string(contexthelper.ContextKeyType_UserInfo)), nil
				},
				Close: func(obj any) error {
					return nil
				},
			},
		)

		return arr
	}
}
