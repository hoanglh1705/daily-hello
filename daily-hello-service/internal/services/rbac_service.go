package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"daily-hello-service/internal/models"
	"daily-hello-service/internal/repositories"

	"go-libs/redisclienthelper"
)

const managerBranchScopeTTL = 5 * time.Minute

type RBACService struct {
	branchRepo     repositories.BranchRepository
	userRepo       *repositories.UserRepository
	deviceRepo     *repositories.DeviceRepository
	attendanceRepo *repositories.AttendanceRepository
	redis          *redisclienthelper.RedisClientHelper
}

func NewRBACService(
	branchRepo repositories.BranchRepository,
	userRepo *repositories.UserRepository,
	deviceRepo *repositories.DeviceRepository,
	attendanceRepo *repositories.AttendanceRepository,
	redis *redisclienthelper.RedisClientHelper,
) *RBACService {
	return &RBACService{
		branchRepo:     branchRepo,
		userRepo:       userRepo,
		deviceRepo:     deviceRepo,
		attendanceRepo: attendanceRepo,
		redis:          redis,
	}
}

func (s *RBACService) GetAllowedBranchIDs(ctx context.Context, role string, currentBranchID *uint) ([]uint, error) {
	switch models.Role(role) {
	case models.RoleAdmin:
		return nil, nil
	case models.RoleManager:
		if currentBranchID == nil || *currentBranchID == 0 {
			return nil, fmt.Errorf("forbidden")
		}

		cacheKey := fmt.Sprintf("rbac:manager:branches:%d", *currentBranchID)
		if s.redis != nil && s.redis.Client != nil {
			if cached, err := s.redis.Client.Get(ctx, cacheKey).Result(); err == nil && cached != "" {
				var branchIDs []uint
				if unmarshalErr := json.Unmarshal([]byte(cached), &branchIDs); unmarshalErr == nil {
					return branchIDs, nil
				}
			}
		}

		branch, err := s.branchRepo.FindByID(ctx, *currentBranchID)
		if err != nil {
			return nil, err
		}

		branchIDs := []uint{branch.ID}
		if branch.Status == "active" && branch.BranchCode != "" {
			children, err := s.branchRepo.FindByParentBranchCode(ctx, branch.BranchCode)
			if err != nil {
				return nil, err
			}
			for _, child := range children {
				branchIDs = append(branchIDs, child.ID)
			}
		}

		if s.redis != nil && s.redis.Client != nil {
			if payload, err := json.Marshal(branchIDs); err == nil {
				_ = s.redis.Client.Set(ctx, cacheKey, payload, managerBranchScopeTTL).Err()
			}
		}

		return branchIDs, nil
	default:
		if currentBranchID == nil || *currentBranchID == 0 {
			return nil, fmt.Errorf("forbidden")
		}
		return []uint{*currentBranchID}, nil
	}
}

func (s *RBACService) EnsureBranchAccess(ctx context.Context, role string, currentBranchID *uint, targetBranchID uint) error {
	if models.Role(role) == models.RoleAdmin {
		return nil
	}

	allowedBranchIDs, err := s.GetAllowedBranchIDs(ctx, role, currentBranchID)
	if err != nil {
		return err
	}

	for _, branchID := range allowedBranchIDs {
		if branchID == targetBranchID {
			return nil
		}
	}

	return fmt.Errorf("forbidden")
}

func (s *RBACService) EnsureUserAccess(ctx context.Context, role string, currentBranchID *uint, targetUserID uint) error {
	if models.Role(role) == models.RoleAdmin {
		return nil
	}

	user, err := s.userRepo.FindByID(ctx, targetUserID)
	if err != nil {
		return err
	}
	if user.BranchID == nil {
		return fmt.Errorf("forbidden")
	}

	return s.EnsureBranchAccess(ctx, role, currentBranchID, *user.BranchID)
}

func (s *RBACService) EnsureDeviceAccess(ctx context.Context, role string, currentBranchID *uint, deviceID uint) error {
	if models.Role(role) == models.RoleAdmin {
		return nil
	}

	device, err := s.deviceRepo.FindByID(ctx, deviceID)
	if err != nil {
		return err
	}
	if device.User == nil || device.User.BranchID == nil {
		return fmt.Errorf("forbidden")
	}

	return s.EnsureBranchAccess(ctx, role, currentBranchID, *device.User.BranchID)
}

func (s *RBACService) EnsureAttendanceAccess(ctx context.Context, role string, currentBranchID *uint, attendanceID uint) error {
	if models.Role(role) == models.RoleAdmin {
		return nil
	}

	attendance, err := s.attendanceRepo.FindByID(ctx, attendanceID)
	if err != nil {
		return err
	}

	return s.EnsureBranchAccess(ctx, role, currentBranchID, attendance.BranchID)
}
