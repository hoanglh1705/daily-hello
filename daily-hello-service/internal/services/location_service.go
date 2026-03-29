package services

import (
	"math"

	"daily-hello-service/internal/models"
	"daily-hello-service/internal/repositories"
)

type LocationService struct {
	branchRepo     repositories.BranchRepository
	branchWifiRepo repositories.BranchWifiRepository
}

func NewLocationService(branchRepo repositories.BranchRepository, branchWifiRepo repositories.BranchWifiRepository) *LocationService {
	return &LocationService{branchRepo: branchRepo, branchWifiRepo: branchWifiRepo}
}

// IsValidGPS checks if the given lat/lng is within the branch's configured radius
func (s *LocationService) IsValidGPS(branch *models.Branch, lat, lng float64) bool {
	if branch.Lat == nil || branch.Lng == nil || branch.Radius == nil {
		return false
	}
	distance := haversine(*branch.Lat, *branch.Lng, lat, lng)
	return distance <= float64(*branch.Radius)
}

// IsValidWifi checks if the given BSSID matches any of the branch's configured WiFi networks
func (s *LocationService) IsValidWifi(branch *models.Branch, bssid string) bool {
	if bssid == "" {
		return false
	}
	for _, wifi := range branch.WifiList {
		if wifi.BSSID == bssid {
			return true
		}
	}
	return false
}

// haversine calculates the distance in meters between two GPS coordinates
func haversine(lat1, lng1, lat2, lng2 float64) float64 {
	const earthRadius = 6371000 // meters

	dLat := toRadians(lat2 - lat1)
	dLng := toRadians(lng2 - lng1)

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(toRadians(lat1))*math.Cos(toRadians(lat2))*
			math.Sin(dLng/2)*math.Sin(dLng/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadius * c
}

func toRadians(deg float64) float64 {
	return deg * math.Pi / 180
}
