package services

import (
	"context"
	"time"

	"daily-hello-service/internal/models"
	appErrors "daily-hello-service/internal/pkg/errors"
	"daily-hello-service/internal/repositories"

	"gorm.io/gorm"
)

type AttendanceService struct {
	repo            *repositories.AttendanceRepository
	branchRepo      repositories.BranchRepository
	locationService *LocationService
	timezone        *time.Location
}

func NewAttendanceService(
	repo *repositories.AttendanceRepository,
	branchRepo repositories.BranchRepository,
	locationService *LocationService,
	timezone *time.Location,
) *AttendanceService {
	return &AttendanceService{
		repo:            repo,
		branchRepo:      branchRepo,
		locationService: locationService,
		timezone:        timezone,
	}
}

func (s *AttendanceService) CheckIn(ctx context.Context, userID uint, req models.AttendanceRequest) (*models.Attendance, error) {
	// 1. Fetch branch
	branch, err := s.branchRepo.FindByID(ctx, req.BranchID)
	if err != nil {
		return nil, appErrors.ErrBranchNotFound
	}

	// 2. Validate location (GPS or WiFi)
	validGPS := s.locationService.IsValidGPS(branch, req.Lat, req.Lng)
	validWifi := s.locationService.IsValidWifi(branch, req.WifiBSSID)
	if !validGPS && !validWifi {
		return nil, appErrors.ErrInvalidLocation
	}

	// 3. Check if already checked in today
	todayStart, tomorrowStart := s.todayRange()
	_, err = s.repo.FindTodayCheckIn(ctx, userID, todayStart, tomorrowStart)
	if err == nil {
		return nil, appErrors.ErrAlreadyCheckedIn
	}
	if err != gorm.ErrRecordNotFound {
		return nil, appErrors.ErrInternal
	}

	// 4. Create attendance record with check-in
	now := time.Now().In(s.timezone)
	checkInLat := req.Lat
	checkInLng := req.Lng
	checkInType := "gps"
	checkInStatus := models.StatusWaitingApprove
	if validWifi {
		checkInType = "wifi"
		checkInStatus = models.StatusApproved
	}

	att := &models.Attendance{
		UserID:           userID,
		BranchID:         req.BranchID,
		CheckInTime:      &now,
		CheckInLat:       &checkInLat,
		CheckInLng:       &checkInLng,
		CheckInWifiBSSID: req.WifiBSSID,
		CheckInDeviceID:  req.DeviceID,
		CheckInType:      checkInType,
		CheckInStatus:    checkInStatus,
	}

	if err := s.repo.Create(ctx, att); err != nil {
		return nil, appErrors.ErrInternal
	}

	return att, nil
}

func (s *AttendanceService) CheckOut(ctx context.Context, userID uint, req models.AttendanceRequest) (*models.Attendance, error) {
	// 1. Verify checked in today
	todayStart, tomorrowStart := s.todayRange()
	att, err := s.repo.FindTodayCheckIn(ctx, userID, todayStart, tomorrowStart)
	if err != nil {
		return nil, appErrors.ErrNotCheckedIn
	}

	// 2. Fetch branch
	branch, err := s.branchRepo.FindByID(ctx, req.BranchID)
	if err != nil {
		return nil, appErrors.ErrBranchNotFound
	}

	// 3. Validate location
	validGPS := s.locationService.IsValidGPS(branch, req.Lat, req.Lng)
	validWifi := s.locationService.IsValidWifi(branch, req.WifiBSSID)
	if !validGPS && !validWifi {
		return nil, appErrors.ErrInvalidLocation
	}

	// 4. Update check-out on existing record
	checkOutType := "gps"
	checkOutStatus := models.StatusWaitingApprove
	if validWifi {
		checkOutType = "wifi"
		checkOutStatus = models.StatusApproved
	}

	now := time.Now().In(s.timezone)
	checkOutLat := req.Lat
	checkOutLng := req.Lng

	updates := map[string]interface{}{
		"check_out_time":       now,
		"check_out_lat":        checkOutLat,
		"check_out_lng":        checkOutLng,
		"check_out_wifi_bssid": req.WifiBSSID,
		"check_out_device_id":  req.DeviceID,
		"check_out_type":       checkOutType,
		"check_out_status":     checkOutStatus,
	}

	if err := s.repo.UpdateCheckOut(ctx, att.ID, updates); err != nil {
		return nil, appErrors.ErrInternal
	}

	att.CheckOutTime = &now
	att.CheckOutLat = &checkOutLat
	att.CheckOutLng = &checkOutLng
	att.CheckOutType = checkOutType
	att.CheckOutWifiBSSID = req.WifiBSSID
	att.CheckOutDeviceID = req.DeviceID
	att.CheckOutStatus = checkOutStatus

	return att, nil
}

func (s *AttendanceService) GetHistory(ctx context.Context, filter models.AttendanceFilter, pq models.PaginationQuery) (*models.PaginatedResponse, error) {
	items, total, err := s.repo.List(ctx, filter, pq)
	if err != nil {
		return nil, appErrors.ErrInternal
	}

	return &models.PaginatedResponse{
		Items: items,
		Meta: models.PaginationMeta{
			Page:  pq.GetPage(),
			Limit: pq.GetLimit(),
			Total: total,
		},
	}, nil
}

func (s *AttendanceService) GetToday(ctx context.Context, userID uint) (*models.Attendance, error) {
	todayStart, tomorrowStart := s.todayRange()
	att, err := s.repo.FindTodayByUserID(ctx, userID, todayStart, tomorrowStart)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, appErrors.ErrInternal
	}
	return att, nil
}

func (s *AttendanceService) ApproveCheckIn(ctx context.Context, id uint) (*models.Attendance, error) {
	return s.updatePendingStatus(ctx, id, "check_in_status", models.StatusApproved)
}

func (s *AttendanceService) RejectCheckIn(ctx context.Context, id uint) (*models.Attendance, error) {
	return s.updatePendingStatus(ctx, id, "check_in_status", models.StatusRejected)
}

func (s *AttendanceService) ApproveCheckOut(ctx context.Context, id uint) (*models.Attendance, error) {
	return s.updatePendingStatus(ctx, id, "check_out_status", models.StatusApproved)
}

func (s *AttendanceService) RejectCheckOut(ctx context.Context, id uint) (*models.Attendance, error) {
	return s.updatePendingStatus(ctx, id, "check_out_status", models.StatusRejected)
}

func (s *AttendanceService) updatePendingStatus(ctx context.Context, id uint, field string, status models.AttendanceStatus) (*models.Attendance, error) {
	att, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, appErrors.ErrNotFound
		}
		return nil, appErrors.ErrInternal
	}

	switch field {
	case "check_in_status":
		if att.CheckInStatus != models.StatusWaitingApprove {
			return nil, appErrors.ErrInvalidInput
		}
	case "check_out_status":
		if att.CheckOutTime == nil || att.CheckOutStatus != models.StatusWaitingApprove {
			return nil, appErrors.ErrInvalidInput
		}
	default:
		return nil, appErrors.ErrInvalidInput
	}

	if err := s.repo.UpdateStatus(ctx, id, map[string]interface{}{field: status}); err != nil {
		return nil, appErrors.ErrInternal
	}

	if field == "check_in_status" {
		att.CheckInStatus = status
	} else {
		att.CheckOutStatus = status
	}

	return att, nil
}

func (s *AttendanceService) todayRange() (time.Time, time.Time) {
	now := time.Now().In(s.timezone)
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, s.timezone)
	return today, today.Add(24 * time.Hour)
}
