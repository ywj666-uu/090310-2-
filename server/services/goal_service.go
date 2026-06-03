package services

import (
	"math"
	"time"

	"gorm.io/gorm"

	"carbon-tracker/models"
)

type GoalService struct {
	db *gorm.DB
}

func NewGoalService(db *gorm.DB) *GoalService {
	return &GoalService{db: db}
}

func (s *GoalService) CreateGoal(userID uint, req *models.GoalRequest) (*models.ReductionGoal, error) {
	s.db.Model(&models.ReductionGoal{}).
		Where("user_id = ? AND status = ?", userID, "active").
		Update("status", "cancelled")

	goal := &models.ReductionGoal{
		UserID:         userID,
		TargetEmission: req.TargetEmission,
		StartDate:      req.StartDate,
		EndDate:        req.EndDate,
		Status:         "active",
	}

	if err := s.db.Create(goal).Error; err != nil {
		return nil, err
	}
	return goal, nil
}

func (s *GoalService) GetActiveGoal(userID uint) (*models.GoalStatus, error) {
	var goal models.ReductionGoal
	if err := s.db.Where("user_id = ? AND status = ?", userID, "active").First(&goal).Error; err != nil {
		return nil, err
	}

	startDate, _ := time.Parse("2006-01-02", goal.StartDate)
	endDate, _ := time.Parse("2006-01-02", goal.EndDate)
	now := time.Now()

	var totalEmission float64
	var dayCount int64
	s.db.Model(&models.CarbonRecord{}).
		Where("user_id = ? AND record_date BETWEEN ? AND ?", userID, goal.StartDate, now.Format("2006-01-02")).
		Select("COALESCE(SUM(emission), 0)").
		Scan(&totalEmission)

	s.db.Model(&models.CarbonRecord{}).
		Where("user_id = ? AND record_date BETWEEN ? AND ?", userID, goal.StartDate, now.Format("2006-01-02")).
		Select("COUNT(DISTINCT record_date)").
		Scan(&dayCount)

	var currentAvg float64
	if dayCount > 0 {
		currentAvg = math.Round(totalEmission/float64(dayCount)*100) / 100
	}

	achieved := currentAvg <= goal.TargetEmission && currentAvg > 0
	daysRemaining := int(endDate.Sub(now).Hours() / 24)
	if daysRemaining < 0 {
		daysRemaining = 0
	}

	totalDays := endDate.Sub(startDate).Hours() / 24
	elapsedDays := now.Sub(startDate).Hours() / 24
	completionRate := math.Min(elapsedDays/totalDays*100, 100)

	if now.After(endDate) && achieved {
		s.db.Model(&goal).Update("status", "completed")
		goal.Status = "completed"
	}

	encouragement := ""
	if achieved {
		idx := int(now.Unix()) % len(encouragements)
		encouragement = encouragements[idx]
	}

	return &models.GoalStatus{
		Goal:           goal,
		CurrentAvg:     currentAvg,
		Achieved:       achieved,
		Encouragement:  encouragement,
		DaysRemaining:  daysRemaining,
		CompletionRate: math.Round(completionRate*100) / 100,
	}, nil
}

func (s *GoalService) GetGoalHistory(userID uint) ([]models.ReductionGoal, error) {
	var goals []models.ReductionGoal
	err := s.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&goals).Error
	return goals, err
}
