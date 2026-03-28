package migration

import (
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

// EnablePostgreSQL: remove this and all tx.Set() functions bellow
// var defaultTableOpts = "ENGINE=InnoDB ROW_FORMAT=DYNAMIC"
var defaultTableOpts = ""

// Base represents base columns for all tables. Do not use gorm.Model because of uint ID
type Base struct {
	CreatedAt time.Time      `gorm:"column:created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index;column:deleted_at"`
}

type BaseWithID struct {
	ID        int64          `gorm:"column:id;primaryKey;autoIncrement"`
	CreatedAt time.Time      `gorm:"column:created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index;column:deleted_at"`
}

// BaseWithoutDeletedAt represents base columns for all tables without deleted_at. Do not use gorm.Model because of uint ID
type BaseWithoutDeletedAt struct {
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func migrateTable() []*gormigrate.Migration {
	migrations := []*gormigrate.Migration{}
	// Smart Attendance tables
	migrations = append(migrations, &gormigrate.Migration{
		ID: "202603271230",
		Migrate: func(tx *gorm.DB) error {
			type Branch struct {
				ID               uint      `gorm:"column:id;primaryKey;autoIncrement"`
				BranchCode       string    `gorm:"column:branch_code;type:varchar(100);uniqueIndex;not null"`
				ParentBranchCode string    `gorm:"column:parent_branch_code;type:varchar(100)"`
				Name             string    `gorm:"column:name;type:varchar(100);not null"`
				Address          string    `gorm:"column:address;type:text"`
				Lat              *float64  `gorm:"column:lat;type:double precision"`
				Lng              *float64  `gorm:"column:lng;type:double precision"`
				Radius           *int      `gorm:"column:radius;type:int"`
				Status           string    `gorm:"column:status;type:varchar(20);default:'active'"`
				CreatedAt        time.Time `gorm:"column:created_at"`
				UpdatedAt        time.Time `gorm:"column:updated_at"`
			}

			type User struct {
				ID        uint      `gorm:"column:id;primaryKey;autoIncrement"`
				Name      string    `gorm:"column:name;type:varchar(100);not null"`
				Code      string    `gorm:"column:code;type:varchar(100);not null"`
				Email     string    `gorm:"column:email;type:varchar(150);uniqueIndex;not null"`
				Phone     string    `gorm:"column:phone;type:varchar(50)"`
				Password  string    `gorm:"column:password;type:text;not null"`
				Role      string    `gorm:"column:role;type:varchar(20);not null"`
				BranchID  *uint     `gorm:"column:branch_id;index"`
				Status    string    `gorm:"column:status;type:varchar(20);default:'active'"`
				CreatedAt time.Time `gorm:"column:created_at"`
				UpdatedAt time.Time `gorm:"column:updated_at"`
			}

			type BranchWifi struct {
				ID        uint      `gorm:"column:id;primaryKey;autoIncrement"`
				Code      string    `gorm:"column:code;type:varchar(100);not null"`
				Name      string    `gorm:"column:name;type:varchar(100);not null"`
				BranchID  uint      `gorm:"column:branch_id;index;not null"`
				SSID      string    `gorm:"column:ssid;type:varchar(100)"`
				BSSID     string    `gorm:"column:bssid;type:varchar(100)"`
				CreatedAt time.Time `gorm:"column:created_at"`
			}

			type Attendance struct {
				ID           uint       `gorm:"column:id;primaryKey;autoIncrement"`
				UserID       uint       `gorm:"column:user_id;index;not null"`
				BranchID     uint       `gorm:"column:branch_id;index;not null"`
				CheckInTime  *time.Time `gorm:"column:check_in_time;index"`
				CheckOutTime *time.Time `gorm:"column:check_out_time"`
				CheckInLat   *float64   `gorm:"column:check_in_lat;type:double precision"`
				CheckInLng   *float64   `gorm:"column:check_in_lng;type:double precision"`
				CheckOutLat  *float64   `gorm:"column:check_out_lat;type:double precision"`
				CheckOutLng  *float64   `gorm:"column:check_out_lng;type:double precision"`
				WifiBSSID    string     `gorm:"column:wifi_bssid;type:varchar(100)"`
				DeviceID     string     `gorm:"column:device_id;type:varchar(100)"`
				Status       string     `gorm:"column:status;type:varchar(20)"`
				CreatedAt    time.Time  `gorm:"column:created_at"`
			}

			type Shift struct {
				ID        uint      `gorm:"column:id;primaryKey;autoIncrement"`
				BranchID  *uint     `gorm:"column:branch_id;index"`
				StartTime string    `gorm:"column:start_time;type:time"`
				EndTime   string    `gorm:"column:end_time;type:time"`
				CreatedAt time.Time `gorm:"column:created_at"`
			}

			type AttendanceSummary struct {
				ID            uint       `gorm:"column:id;primaryKey;autoIncrement"`
				UserID        *uint      `gorm:"column:user_id;index"`
				Date          *time.Time `gorm:"column:date;type:date"`
				TotalHours    *float64   `gorm:"column:total_hours"`
				OvertimeHours *float64   `gorm:"column:overtime_hours"`
				Status        string     `gorm:"column:status;type:varchar(20)"`
			}

			type Device struct {
				ID        uint      `gorm:"column:id;primaryKey;autoIncrement"`
				UserID    *uint     `gorm:"column:user_id;index"`
				DeviceID  string    `gorm:"column:device_id;type:varchar(100)"`
				IsTrusted *bool     `gorm:"column:is_trusted;default:true"`
				CreatedAt time.Time `gorm:"column:created_at"`
			}

			// AutoMigrate all tables
			if err := tx.Set("gorm:table_options", defaultTableOpts).AutoMigrate(
				&Branch{},
				&User{},
				&BranchWifi{},
				&Attendance{},
				&Shift{},
				&AttendanceSummary{},
				&Device{},
			); err != nil {
				return err
			}

			// Add FK constraints
			if err := tx.Exec(`ALTER TABLE users ADD CONSTRAINT fk_users_branch FOREIGN KEY (branch_id) REFERENCES branches(id)`).Error; err != nil {
				// Ignore if constraint already exists
			}

			if err := tx.Exec(`ALTER TABLE branch_wifis ADD CONSTRAINT fk_branch_wifis_branch FOREIGN KEY (branch_id) REFERENCES branches(id)`).Error; err != nil {
				// Ignore if constraint already exists
			}

			if err := tx.Exec(`ALTER TABLE attendances ADD CONSTRAINT fk_attendances_user FOREIGN KEY (user_id) REFERENCES users(id)`).Error; err != nil {
				// Ignore if constraint already exists
			}

			if err := tx.Exec(`ALTER TABLE attendances ADD CONSTRAINT fk_attendances_branch FOREIGN KEY (branch_id) REFERENCES branches(id)`).Error; err != nil {
				// Ignore if constraint already exists
			}

			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			if err := tx.Migrator().DropTable("devices"); err != nil {
				return err
			}
			if err := tx.Migrator().DropTable("attendance_summaries"); err != nil {
				return err
			}
			if err := tx.Migrator().DropTable("shifts"); err != nil {
				return err
			}
			if err := tx.Migrator().DropTable("attendances"); err != nil {
				return err
			}
			if err := tx.Migrator().DropTable("branch_wifis"); err != nil {
				return err
			}
			if err := tx.Migrator().DropTable("users"); err != nil {
				return err
			}
			if err := tx.Migrator().DropTable("branches"); err != nil {
				return err
			}
			return nil
		},
	})

	// Add branch_code and parent_branch_code to branches
	migrations = append(migrations, &gormigrate.Migration{
		ID: "202603271339",
		Migrate: func(tx *gorm.DB) error {
			if err := tx.Exec(`ALTER TABLE branches ADD COLUMN IF NOT EXISTS branch_code VARCHAR(100) NOT NULL DEFAULT ''`).Error; err != nil {
				return err
			}
			if err := tx.Exec(`ALTER TABLE branches ADD COLUMN IF NOT EXISTS parent_branch_code VARCHAR(100)`).Error; err != nil {
				return err
			}
			if err := tx.Exec(`CREATE UNIQUE INDEX IF NOT EXISTS idx_branches_branch_code ON branches(branch_code)`).Error; err != nil {
				return err
			}
			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			if err := tx.Exec(`DROP INDEX IF EXISTS idx_branches_branch_code`).Error; err != nil {
				return err
			}
			if err := tx.Exec(`ALTER TABLE branches DROP COLUMN IF EXISTS parent_branch_code`).Error; err != nil {
				return err
			}
			if err := tx.Exec(`ALTER TABLE branches DROP COLUMN IF EXISTS branch_code`).Error; err != nil {
				return err
			}
			return nil
		},
	})

	// Add refresh_tokens table
	migrations = append(migrations, &gormigrate.Migration{
		ID: "202603280900",
		Migrate: func(tx *gorm.DB) error {
			type RefreshToken struct {
				ID        uint      `gorm:"column:id;primaryKey;autoIncrement"`
				UserID    uint      `gorm:"column:user_id;not null;index"`
				Token     string    `gorm:"column:token;type:varchar(512);uniqueIndex;not null"`
				ExpiresAt time.Time `gorm:"column:expires_at;not null"`
				CreatedAt time.Time `gorm:"column:created_at"`
			}

			if err := tx.Set("gorm:table_options", defaultTableOpts).AutoMigrate(&RefreshToken{}); err != nil {
				return err
			}

			tx.Exec(`ALTER TABLE refresh_tokens ADD CONSTRAINT fk_refresh_tokens_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE`)

			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable("refresh_tokens")
		},
	})

	return migrations
}
