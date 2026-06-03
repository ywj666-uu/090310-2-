package services

import (
	"fmt"
	"math"
	"time"

	"gorm.io/gorm"

	"carbon-tracker/models"
)

type CarbonService struct {
	db *gorm.DB
}

func NewCarbonService(db *gorm.DB) *CarbonService {
	return &CarbonService{db: db}
}

func (s *CarbonService) CheckGoalAfterRecord(userID uint, recordDate string) *models.GoalStatus {
	var goal models.ReductionGoal
	if err := s.db.Where("user_id = ? AND status = ?", userID, "active").First(&goal).Error; err != nil {
		return nil
	}

	var todayTotal float64
	s.db.Model(&models.CarbonRecord{}).
		Where("user_id = ? AND record_date = ?", userID, recordDate).
		Select("COALESCE(SUM(emission), 0)").
		Scan(&todayTotal)

	if todayTotal > 0 && todayTotal <= goal.TargetEmission {
		idx := int(time.Now().UnixNano()) % len(encouragements)
		return &models.GoalStatus{
			Goal:          goal,
			CurrentAvg:    todayTotal,
			Achieved:      true,
			Encouragement: encouragements[idx],
		}
	}
	return nil
}

var encouragements = []string{
	"太棒了！你正在为地球做出积极的贡献！🌍",
	"继续保持！每一点减排都是对未来的投资！💚",
	"你是环保先锋！你的努力正在改变世界！🌱",
	"目标达成！你证明了绿色生活可以很简单！🎉",
	"了不起！你的碳足迹比区域平均值还低！⭐",
	"坚持就是胜利！你的减排习惯正在成为日常！🏆",
}

func (s *CarbonService) GetEmissionFactors() ([]models.EmissionFactor, error) {
	var factors []models.EmissionFactor
	err := s.db.Find(&factors).Error
	return factors, err
}

func (s *CarbonService) CreateRecord(userID uint, req *models.RecordRequest) (*models.CarbonRecord, error) {
	var factor models.EmissionFactor
	if err := s.db.Where("category = ? AND item = ?", req.Category, req.Item).First(&factor).Error; err != nil {
		return nil, fmt.Errorf("未找到对应的碳排因子: %s - %s", req.Category, req.Item)
	}

	emission := req.Amount * factor.Factor

	record := &models.CarbonRecord{
		UserID:     userID,
		RecordDate: req.RecordDate,
		Category:   req.Category,
		Item:       req.Item,
		Amount:     req.Amount,
		Emission:   math.Round(emission*10000) / 10000,
		Note:       req.Note,
	}

	if err := s.db.Create(record).Error; err != nil {
		return nil, err
	}
	return record, nil
}

func (s *CarbonService) GetDailyRecords(userID uint, date string) ([]models.CarbonRecord, error) {
	var records []models.CarbonRecord
	err := s.db.Where("user_id = ? AND record_date = ?", userID, date).Find(&records).Error
	return records, err
}

func (s *CarbonService) GetDailySummary(userID uint, date string) (float64, error) {
	var total float64
	err := s.db.Model(&models.CarbonRecord{}).
		Where("user_id = ? AND record_date = ?", userID, date).
		Select("COALESCE(SUM(emission), 0)").
		Scan(&total).Error
	return total, err
}

func (s *CarbonService) GetWeeklyTrend(userID uint, region string) (*models.TrendData, error) {
	now := time.Now()
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	monday := now.AddDate(0, 0, -(weekday - 1))
	sunday := monday.AddDate(0, 0, 6)
	return s.getTrend(userID, monday, sunday, region)
}

func (s *CarbonService) GetMonthlyTrend(userID uint, region string) (*models.TrendData, error) {
	now := time.Now()
	firstDay := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	lastDay := firstDay.AddDate(0, 1, -1)
	return s.getTrend(userID, firstDay, lastDay, region)
}

func (s *CarbonService) getTrend(userID uint, startDate, endDate time.Time, region string) (*models.TrendData, error) {
	var summaries []models.DailySummary
	err := s.db.Model(&models.CarbonRecord{}).
		Where("user_id = ? AND record_date BETWEEN ? AND ?", userID, startDate.Format("2006-01-02"), endDate.Format("2006-01-02")).
		Select("record_date as date, SUM(emission) as total_emission").
		Group("record_date").
		Order("record_date").
		Scan(&summaries).Error
	if err != nil {
		return nil, err
	}

	filled := s.fillMissingDates(summaries, startDate, endDate)

	regionalAvg := s.calcRegionalAverage(region, startDate, endDate)

	return &models.TrendData{
		UserData:        filled,
		RegionalAverage: regionalAvg,
	}, nil
}

func (s *CarbonService) fillMissingDates(summaries []models.DailySummary, startDate, endDate time.Time) []models.DailySummary {
	existing := make(map[string]float64)
	for _, d := range summaries {
		existing[d.Date] = d.TotalEmission
	}

	var filled []models.DailySummary
	for cur := startDate; !cur.After(endDate); cur = cur.AddDate(0, 0, 1) {
		dateStr := cur.Format("2006-01-02")
		emission := existing[dateStr]
		filled = append(filled, models.DailySummary{Date: dateStr, TotalEmission: emission})
	}
	return filled
}

func (s *CarbonService) calcRegionalAverage(region string, startDate, endDate time.Time) float64 {
	var result struct {
		Avg float64
	}
	s.db.Model(&models.CarbonRecord{}).
		Joins("JOIN users ON users.id = carbon_records.user_id").
		Where("users.region = ? AND carbon_records.record_date BETWEEN ? AND ?", region, startDate.Format("2006-01-02"), endDate.Format("2006-01-02")).
		Select("COALESCE(SUM(carbon_records.emission) / NULLIF(COUNT(DISTINCT carbon_records.user_id), 0) / NULLIF(DATEDIFF(?, ?)+1, 0), 0) as avg", endDate.Format("2006-01-02"), startDate.Format("2006-01-02")).
		Scan(&result)

	return math.Round(result.Avg*100) / 100
}

func (s *CarbonService) GetCategorySummary(userID uint, startDate, endDate string) ([]models.CategorySummary, error) {
	var summaries []models.CategorySummary
	err := s.db.Model(&models.CarbonRecord{}).
		Where("user_id = ? AND record_date BETWEEN ? AND ?", userID, startDate, endDate).
		Select("category, SUM(emission) as total_emission").
		Group("category").
		Scan(&summaries).Error
	return summaries, err
}

func (s *CarbonService) DeleteRecord(userID uint, recordID uint) error {
	return s.db.Where("id = ? AND user_id = ?", recordID, userID).Delete(&models.CarbonRecord{}).Error
}
