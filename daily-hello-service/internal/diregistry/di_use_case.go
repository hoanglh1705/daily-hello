package diregistry

import (
	"daily-hello-service/internal/repositories"
	"daily-hello-service/internal/services"
	"go-libs/dihelper"

	"github.com/sarulabs/di"
)

func initUseCasesBuilder() {
	dihelper.UsecasesBuilder = func() []di.Def {
		arr := []di.Def{}
		arr = append(arr,
			di.Def{
				Name:  BranchServiceDIName,
				Scope: di.App,
				Build: func(ctn di.Container) (any, error) {
					branchRepository := ctn.Get(BranchRepositoryDIName).(repositories.BranchRepository)
					return services.NewBranchService(branchRepository), nil
				},
				Close: func(obj any) error {
					return nil
				},
			},
			di.Def{
				Name:  BranchWifiServiceDIName,
				Scope: di.App,
				Build: func(ctn di.Container) (any, error) {
					branchWifiRepo := ctn.Get(BranchWifiRepositoryDIName).(repositories.BranchWifiRepository)
					branchRepo := ctn.Get(BranchRepositoryDIName).(repositories.BranchRepository)
					return services.NewBranchWifiService(branchWifiRepo, branchRepo), nil
				},
				Close: func(obj any) error {
					return nil
				},
			},
		)

		return arr
	}
}
