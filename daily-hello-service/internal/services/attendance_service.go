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
}

func NewAttendanceService(
	repo *repositories.AttendanceRepository,
	branchRepo repositories.BranchRepository,
	locationService *LocationService,
) *AttendanceService {
	return &AttendanceService{
		repo:            repo,
		branchRepo:      branchRepo,
		locationService: locationService,
	}
}

func (s *AttendanceService) CheckIn(ctx context.Context, userID uint, req models.CheckInRequest) (*models.Attendance, error) {
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
	_, err = s.repo.FindTodayCheckIn(ctx, userID)
	if err == nil {
		return nil, appErrors.ErrAlreadyCheckedIn
	}
	if err != gorm.ErrRecordNotFound {
		return nil, appErrors.ErrInternal
	}

	// 4. Create attendance record with check-in
	now := time.Now()
	checkInLat := req.Lat
	checkInLng := req.Lng
	att := &models.Attendance{
		UserID:      userID,
		BranchID:    req.BranchID,
		CheckInTime: &now,
		CheckInLat:  &checkInLat,
		CheckInLng:  &checkInLng,
		WifiBSSID:   req.WifiBSSID,
		DeviceID:    req.DeviceID,
		Status:      models.StatusOnTime, // TODO: compute based on shift schedule
	}

	if err := s.repo.Create(ctx, att); err != nil {
		return nil, appErrors.ErrInternal
	}

	return att, nil
}

func (s *AttendanceService) CheckOut(ctx context.Context, userID uint, req models.CheckOutRequest) (*models.Attendance, error) {
	// 1. Verify checked in today
	att, err := s.repo.FindTodayCheckIn(ctx, userID)
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
	now := time.Now()
	checkOutLat := req.Lat
	checkOutLng := req.Lng
	if err := s.repo.UpdateCheckOut(ctx, att.ID, now, &checkOutLat, &checkOutLng); err != nil {
		return nil, appErrors.ErrInternal
	}

	att.CheckOutTime = &now
	att.CheckOutLat = &checkOutLat
	att.CheckOutLng = &checkOutLng

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
