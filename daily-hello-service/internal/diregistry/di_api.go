package diregistry

import (
	"go-libs/dihelper"

	"github.com/sarulabs/di"
)

func initApiBuilder() {
	dihelper.APIsBuilder = func() []di.Def {
		arr := []di.Def{}
		// arr = append(arr,
		// 	di.Def{
		// 		Name:  AdminCompanyAPIDIName,
		// 		Scope: di.App,
		// 		Build: func(ctn di.Container) (any, error) {
		// 			cuc := ctn.Get(CompanyUseCaseDIName).(company_usecase.CompanyUseCase)
		// 			return adminCompany.NewAdminCompanyController(cuc), nil
		// 		},
		// 		Close: func(obj any) error {
		// 			return nil
		// 		},
		// 	},
		// )

		return arr
	}
}
