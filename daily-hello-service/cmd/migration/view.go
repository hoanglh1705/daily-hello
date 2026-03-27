package migration

import (
	"github.com/go-gormigrate/gormigrate/v2"
)

func migrateView() []*gormigrate.Migration {
	migrations := []*gormigrate.Migration{}
	// migrations = append(migrations, &gormigrate.Migration{
	// 	ID: "202410101011",
	// 	Migrate: func(tx *gorm.DB) error {
	// 		changes := []string{
	// 			`CREATE OR REPLACE VIEW public.view_list_all_contract
	// 				AS SELECT contract.id,
	// 					contract.area_name,
	// 					contract.branch_no,
	// 					contract.branch_name,
	// 					contract.cif_no AS cif,
	// 					contract.contract_number,
	// 					'' AS credit_no,
	// 					contract.full_name,
	// 					contract.issue_date,
	// 					contract.principal_amount,
	// 					contract.interest_amount,
	// 					contract.fee_amount,
	// 					contract.late_payment_days,
	// 					dg.code AS debt_group,
	// 					'CONTRACT' AS contract_type,
	// 					contract.reg_number,
	// 					contract.collateral,
	// 					contract.branch_no_original,
	// 					contract.total_owed_amount
	// 				   FROM contract
	// 				   left join debt_group dg on contract.debt_group_id  = dg.id
	// 				  WHERE 1 = 1 AND contract.deleted_at IS NULL
	// 				  and contract.late_payment_days > 0
	// 				UNION ALL
	// 				 SELECT xlrr_contract.id,
	// 					xlrr_contract.area_name,
	// 					xlrr_contract.branch_no,
	// 					xlrr_contract.branch_name,
	// 					xlrr_contract.cif,
	// 					xlrr_contract.contract_number,
	// 					xlrr_contract.credit_no,
	// 					xlrr_contract.full_name,
	// 					xlrr_contract.issue_date,
	// 					xlrr_contract.principal_amount,
	// 					xlrr_contract.interest_amount,
	// 					xlrr_contract.fee_amount,
	// 					xlrr_contract.late_payment_days,
	// 					xlrr_contract.debt_group,
	// 					'XLRR_CONTRACT' AS contract_type,
	// 					xlrr_contract.reg_number,
	// 					xlrr_contract.collateral,
	// 					xlrr_contract.branch_no_original,
	// 					xlrr_contract.principal_amount + xlrr_contract.interest_amount as total_owed_amount
	// 				   FROM xlrr_contract;`,
	// 		}
	// 		return migratehelper.ExecMultiple(tx, strings.Join(changes, " "))
	// 	},
	// 	Rollback: func(tx *gorm.DB) error {
	// 		changes := []string{
	// 			`DROP VIEW view_list_all_contract;`,
	// 		}
	// 		return migratehelper.ExecMultiple(tx, strings.Join(changes, " "))
	// 	},
	// })

	return migrations
}
