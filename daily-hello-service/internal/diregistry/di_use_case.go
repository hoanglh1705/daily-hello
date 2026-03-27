package diregistry

import (
	"go-libs/dihelper"

	"github.com/sarulabs/di"
)

func initUseCasesBuilder() {
	dihelper.UsecasesBuilder = func() []di.Def {
		arr := []di.Def{}
		// arr = append(arr,
		// 	di.Def{
		// 		Name:  HelloWorldUsecaseDIName,
		// 		Scope: di.App,
		// 		Build: func(ctn di.Container) (any, error) {
		// 			return usecase.NewHelloWorldUsecase(), nil
		// 		},
		// 		Close: func(obj any) error {
		// 			return nil
		// 		},
		// 	},
		// )

		return arr
	}
}
