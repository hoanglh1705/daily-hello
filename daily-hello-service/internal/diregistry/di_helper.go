package diregistry

import (
	"daily-hello-service/config"
	"fmt"
	"go-libs/cachehelper"
	"go-libs/copyhelper"
	"go-libs/dihelper"
	"go-libs/redisclienthelper"
	"go-libs/sqlormhelper"

	"github.com/sarulabs/di"
)

func initHelpersBuilder() {
	dihelper.HelpersBuilder = func() []di.Def {
		arr := []di.Def{}
		arr = append(arr,
			di.Def{
				Name:  DataBaseDIName,
				Scope: di.App,
				Build: func(ctn di.Container) (any, error) {
					cfg := ctn.Get(ConfigDIName).(*config.Config)
					return sqlormhelper.NewGormPostgresqlDB(
						&sqlormhelper.GormConnectionOptions{
							Host:     cfg.Database.Host,
							Port:     cfg.Database.Port,
							Username: cfg.Database.Username,
							Password: cfg.Database.Password,
							Database: cfg.Database.Database,
							SSLMode:  cfg.Database.SSLMode,
							Schema:   cfg.Database.SearchPath,
							Timezone: cfg.Database.Timezone,
						}), nil
				},
				Close: func(obj any) error {
					return nil
				},
			},
			di.Def{
				Name:  ModelConverterDIName,
				Scope: di.App,
				Build: func(ctn di.Container) (any, error) {
					return copyhelper.NewModelConverter(), nil
				},
				Close: func(obj any) error {
					return nil
				},
			},
			di.Def{
				Name:  RedisClientHelperDIName,
				Scope: di.App,
				Build: func(ctn di.Container) (any, error) {
					cfg := ctn.Get(ConfigDIName).(*config.Config)
					return redisclienthelper.NewRedisClientHelper(&redisclienthelper.RedisConfigOptions{
						Addrs: []string{
							fmt.Sprintf("%v:%v", cfg.Cache.Host, cfg.Cache.Port),
						},
						Password: cfg.Cache.Password,
					}), nil
				},
				Close: func(obj any) error {
					return nil
				},
			},
			di.Def{
				Name:  CacheHelperDIName,
				Scope: di.App,
				Build: func(ctn di.Container) (any, error) {
					redisClient := ctn.Get(RedisClientHelperDIName).(*redisclienthelper.RedisClientHelper)
					return cachehelper.NewCacheHelper(&cachehelper.CacheConfigOptions{
						RedisClientHelper: redisClient,
					}), nil
				},
				Close: func(obj any) error {
					return nil
				},
			},
		)
		return arr
	}
}
