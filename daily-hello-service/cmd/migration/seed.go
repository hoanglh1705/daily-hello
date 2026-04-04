package migration

import (
	config "daily-hello-service/config"
	"fmt"
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func migrateData(_ *config.Config) []*gormigrate.Migration {
	migrations := []*gormigrate.Migration{}
	migrations = append(migrations, &gormigrate.Migration{
		ID: "202604011400",
		Migrate: func(tx *gorm.DB) error {
			type Branch struct {
				ID               uint      `gorm:"column:id;primaryKey;autoIncrement"`
				BranchCode       string    `gorm:"column:branch_code"`
				ParentBranchCode string    `gorm:"column:parent_branch_code"`
				Name             string    `gorm:"column:name"`
				Address          string    `gorm:"column:address"`
				Status           string    `gorm:"column:status"`
				CreatedAt        time.Time `gorm:"column:created_at"`
				UpdatedAt        time.Time `gorm:"column:updated_at"`
			}

			const (
				totalBranches = 10000
				parentCount   = 2500
				batchSize     = 1000
				codePrefix    = "SEED-BR-"
			)

			var existed int64
			if err := tx.Table("branches").
				Where("branch_code LIKE ?", codePrefix+"%").
				Count(&existed).Error; err != nil {
				return err
			}
			if existed > 0 {
				return nil
			}

			now := time.Now()
			branches := make([]Branch, 0, totalBranches)
			parentCodes := make([]string, 0, parentCount)

			for i := 1; i <= parentCount; i++ {
				code := fmt.Sprintf("%s%05d", codePrefix, i)
				parentCodes = append(parentCodes, code)
				branches = append(branches, Branch{
					BranchCode:       code,
					ParentBranchCode: "",
					Name:             fmt.Sprintf("Seed Branch %05d", i),
					Address:          fmt.Sprintf("Address for branch %05d", i),
					Status:           "active",
					CreatedAt:        now,
					UpdatedAt:        now,
				})
			}

			for i := parentCount + 1; i <= totalBranches; i++ {
				code := fmt.Sprintf("%s%05d", codePrefix, i)
				parentCode := parentCodes[(i-parentCount-1)%len(parentCodes)]
				branches = append(branches, Branch{
					BranchCode:       code,
					ParentBranchCode: parentCode,
					Name:             fmt.Sprintf("Seed Child Branch %05d", i),
					Address:          fmt.Sprintf("Address for child branch %05d", i),
					Status:           "active",
					CreatedAt:        now,
					UpdatedAt:        now,
				})
			}

			return tx.Table("branches").CreateInBatches(branches, batchSize).Error
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Exec(`DELETE FROM branches WHERE branch_code LIKE 'SEED-BR-%'`).Error
		},
	})

	migrations = append(migrations, &gormigrate.Migration{
		ID: "202604011410",
		Migrate: func(tx *gorm.DB) error {
			type BranchRef struct {
				ID         uint   `gorm:"column:id"`
				BranchCode string `gorm:"column:branch_code"`
			}

			type User struct {
				ID        uint      `gorm:"column:id;primaryKey;autoIncrement"`
				Name      string    `gorm:"column:name"`
				Code      string    `gorm:"column:code"`
				Email     string    `gorm:"column:email"`
				Phone     string    `gorm:"column:phone"`
				Password  string    `gorm:"column:password"`
				Role      string    `gorm:"column:role"`
				BranchID  *uint     `gorm:"column:branch_id"`
				Status    string    `gorm:"column:status"`
				CreatedAt time.Time `gorm:"column:created_at"`
				UpdatedAt time.Time `gorm:"column:updated_at"`
			}

			branchCodes := []string{
				"SEED-BR-00001",
				"SEED-BR-02501",
				"SEED-BR-05001",
				"SEED-BR-07501",
			}

			var existed int64
			if err := tx.Table("users").
				Where("email LIKE ?", "seed.user.%@dailyhello.local").
				Count(&existed).Error; err != nil {
				return err
			}
			if existed > 0 {
				return nil
			}

			var branches []BranchRef
			if err := tx.Table("branches").
				Select("id, branch_code").
				Where("branch_code IN ?", branchCodes).
				Order("branch_code ASC").
				Find(&branches).Error; err != nil {
				return err
			}
			if len(branches) != len(branchCodes) {
				return fmt.Errorf("missing seed branches for user seed")
			}

			hashedPassword, err := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
			if err != nil {
				return err
			}

			now := time.Now()
			users := make([]User, 0, 100)
			userIndex := 1

			for _, branch := range branches {
				for i := 1; i <= 25; i++ {
					role := "employee"
					if i <= 2 {
						role = "manager"
					}

					branchID := branch.ID
					users = append(users, User{
						Name:      fmt.Sprintf("Seed User %03d", userIndex),
						Code:      fmt.Sprintf("SEED-USER-%03d", userIndex),
						Email:     fmt.Sprintf("seed.user.%03d@dailyhello.local", userIndex),
						Phone:     fmt.Sprintf("090%07d", userIndex),
						Password:  string(hashedPassword),
						Role:      role,
						BranchID:  &branchID,
						Status:    "active",
						CreatedAt: now,
						UpdatedAt: now,
					})
					userIndex++
				}
			}

			return tx.Table("users").CreateInBatches(users, 100).Error
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Exec(`DELETE FROM users WHERE email LIKE 'seed.user.%@dailyhello.local'`).Error
		},
	})

	migrations = append(migrations, &gormigrate.Migration{
		ID: "202604011420",
		Migrate: func(tx *gorm.DB) error {
			type SeedUser struct {
				ID       uint  `gorm:"column:id"`
				BranchID *uint `gorm:"column:branch_id"`
			}

			type Attendance struct {
				ID                uint       `gorm:"column:id;primaryKey;autoIncrement"`
				UserID            uint       `gorm:"column:user_id"`
				BranchID          uint       `gorm:"column:branch_id"`
				CheckInTime       *time.Time `gorm:"column:check_in_time"`
				CheckOutTime      *time.Time `gorm:"column:check_out_time"`
				CheckInLat        *float64   `gorm:"column:check_in_lat"`
				CheckInLng        *float64   `gorm:"column:check_in_lng"`
				CheckOutLat       *float64   `gorm:"column:check_out_lat"`
				CheckOutLng       *float64   `gorm:"column:check_out_lng"`
				CheckInType       string     `gorm:"column:check_in_type"`
				CheckOutType      string     `gorm:"column:check_out_type"`
				CheckInWifiBSSID  string     `gorm:"column:check_in_wifi_bssid"`
				CheckOutWifiBSSID string     `gorm:"column:check_out_wifi_bssid"`
				CheckInDeviceID   string     `gorm:"column:check_in_device_id"`
				CheckOutDeviceID  string     `gorm:"column:check_out_device_id"`
				CheckInStatus     string     `gorm:"column:check_in_status"`
				CheckOutStatus    string     `gorm:"column:check_out_status"`
				CheckInImage      string     `gorm:"column:check_in_image"`
				CheckOutImage     string     `gorm:"column:check_out_image"`
				CreatedAt         time.Time  `gorm:"column:created_at"`
			}

			var existed int64
			if err := tx.Table("attendances").
				Where("check_in_device_id LIKE ?", "SEED-ATT-%").
				Count(&existed).Error; err != nil {
				return err
			}
			if existed > 0 {
				return nil
			}

			var users []SeedUser
			if err := tx.Table("users").
				Select("id, branch_id").
				Where("email LIKE ?", "seed.user.%@dailyhello.local").
				Order("id ASC").
				Find(&users).Error; err != nil {
				return err
			}
			if len(users) != 100 {
				return fmt.Errorf("expected 100 seeded users, got %d", len(users))
			}

			const (
				totalDays = 40
				batchSize = 500
			)

			baseLat := 10.776889
			baseLng := 106.700806
			startDay := time.Now().AddDate(0, 0, -(totalDays - 1))
			records := make([]Attendance, 0, len(users)*totalDays)

			for dayOffset := 0; dayOffset < totalDays; dayOffset++ {
				day := startDay.AddDate(0, 0, dayOffset)

				for userIndex, user := range users {
					if user.BranchID == nil {
						return fmt.Errorf("seed user %d missing branch_id", user.ID)
					}

					isLate := (dayOffset+userIndex)%3 == 0
					checkInHour := 7
					checkInMinute := 50 + (userIndex % 8)
					if isLate {
						checkInHour = 8
						checkInMinute = 5 + (userIndex % 20)
					}

					checkIn := time.Date(day.Year(), day.Month(), day.Day(), checkInHour, checkInMinute, 0, 0, day.Location())
					checkOut := time.Date(day.Year(), day.Month(), day.Day(), 17, 30+(userIndex%25), 0, 0, day.Location())
					lat := baseLat + float64(userIndex%10)*0.0001
					lng := baseLng + float64(dayOffset%10)*0.0001
					checkInStatus := "approved"
					checkOutStatus := "approved"
					if isLate {
						checkInStatus = "waiting_approve"
					}

					records = append(records, Attendance{
						UserID:            user.ID,
						BranchID:          *user.BranchID,
						CheckInTime:       &checkIn,
						CheckOutTime:      &checkOut,
						CheckInLat:        &lat,
						CheckInLng:        &lng,
						CheckOutLat:       &lat,
						CheckOutLng:       &lng,
						CheckInType:       "wifi",
						CheckOutType:      "wifi",
						CheckInWifiBSSID:  fmt.Sprintf("SEED-BSSID-%03d", userIndex+1),
						CheckOutWifiBSSID: fmt.Sprintf("SEED-BSSID-%03d", userIndex+1),
						CheckInDeviceID:   fmt.Sprintf("SEED-ATT-%03d-%02d", userIndex+1, dayOffset+1),
						CheckOutDeviceID:  fmt.Sprintf("SEED-ATT-%03d-%02d", userIndex+1, dayOffset+1),
						CheckInStatus:     checkInStatus,
						CheckOutStatus:    checkOutStatus,
						CheckInImage:      "",
						CheckOutImage:     "",
						CreatedAt:         checkIn,
					})
				}
			}

			return tx.Table("attendances").CreateInBatches(records, batchSize).Error
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Exec(`DELETE FROM attendances WHERE check_in_device_id LIKE 'SEED-ATT-%'`).Error
		},
	})

	migrations = append(migrations, &gormigrate.Migration{
		ID: "202604021433",
		Migrate: func(tx *gorm.DB) error {
			type User struct {
				ID        uint      `gorm:"column:id;primaryKey;autoIncrement"`
				Name      string    `gorm:"column:name"`
				Code      string    `gorm:"column:code"`
				Email     string    `gorm:"column:email"`
				Phone     string    `gorm:"column:phone"`
				Password  string    `gorm:"column:password"`
				Role      string    `gorm:"column:role"`
				BranchID  *uint     `gorm:"column:branch_id"`
				Status    string    `gorm:"column:status"`
				CreatedAt time.Time `gorm:"column:created_at"`
				UpdatedAt time.Time `gorm:"column:updated_at"`
			}

			var existed int64
			if err := tx.Table("users").
				Where("email LIKE ?", "seed.admin.%@dailyhello.local").
				Count(&existed).Error; err != nil {
				return err
			}
			if existed > 0 {
				return nil
			}

			hashedPassword, err := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
			if err != nil {
				return err
			}

			now := time.Now()
			admins := make([]User, 0, 10)
			for i := 1; i <= 10; i++ {
				admins = append(admins, User{
					Name:      fmt.Sprintf("Seed Admin %02d", i),
					Code:      fmt.Sprintf("SEED-ADMIN-%02d", i),
					Email:     fmt.Sprintf("seed.admin.%02d@dailyhello.local", i),
					Phone:     fmt.Sprintf("091%07d", i),
					Password:  string(hashedPassword),
					Role:      "admin",
					BranchID:  nil,
					Status:    "active",
					CreatedAt: now,
					UpdatedAt: now,
				})
			}

			return tx.Table("users").CreateInBatches(admins, 10).Error
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Exec(`DELETE FROM users WHERE email LIKE 'seed.admin.%@dailyhello.local'`).Error
		},
	})

	return migrations
}
