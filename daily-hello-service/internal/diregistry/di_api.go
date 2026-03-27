package diregistry

import (
	"daily-hello-service/internal/handlers"
	"daily-hello-service/internal/services"
	"go-libs/dihelper"

	"github.com/sarulabs/di"
)

func initApiBuilder() {
	dihelper.APIsBuilder = func() []di.Def {
		arr := []di.Def{}
		arr = append(arr,
			di.Def{
				Name:  BranchAPIDIName,
				Scope: di.App,
				Build: func(ctn di.Container) (any, error) {
					cuc := ctn.Get(BranchServiceDIName).(*services.BranchService)
					return handlers.NewBranchHandler(cuc), nil
				},
				Close: func(obj any) error {
					return nil
				},
			},
			di.Def{
				Name:  BranchWifiAPIDIName,
				Scope: di.App,
				Build: func(ctn di.Container) (any, error) {
					svc := ctn.Get(BranchWifiServiceDIName).(*services.BranchWifiService)
					return handlers.NewBranchWifiHandler(svc), nil
				},
				Close: func(obj any) error {
					return nil
				},
			},
		)

		return arr
	}
}
