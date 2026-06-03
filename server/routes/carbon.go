package routes

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"carbon-tracker/models"
	"carbon-tracker/services"
)

type CarbonHandler struct {
	service *services.CarbonService
}

func NewCarbonHandler(service *services.CarbonService) *CarbonHandler {
	return &CarbonHandler{service: service}
}

func (h *CarbonHandler) GetFactors(c *gin.Context) {
	factors, err := h.service.GetEmissionFactors()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取碳排因子失败"})
		return
	}
	c.JSON(http.StatusOK, factors)
}

func (h *CarbonHandler) CreateRecord(c *gin.Context) {
	userID := c.GetUint("userID")
	var req models.RecordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数无效: " + err.Error()})
		return
	}

	record, err := h.service.CreateRecord(userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response := gin.H{"record": record}
	if goalStatus := h.service.CheckGoalAfterRecord(userID, req.RecordDate); goalStatus != nil {
		response["goal_achieved"] = true
		response["encouragement"] = goalStatus.Encouragement
	}

	c.JSON(http.StatusCreated, response)
}

func (h *CarbonHandler) GetDailyRecords(c *gin.Context) {
	userID := c.GetUint("userID")
	date := c.Query("date")
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}

	records, err := h.service.GetDailyRecords(userID, date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取记录失败"})
		return
	}

	total, _ := h.service.GetDailySummary(userID, date)

	c.JSON(http.StatusOK, gin.H{
		"records":        records,
		"total_emission": total,
		"date":           date,
	})
}

func (h *CarbonHandler) GetWeeklyTrend(c *gin.Context) {
	userID := c.GetUint("userID")
	region := c.Query("region")
	if region == "" {
		region = "全国"
	}
	trend, err := h.service.GetWeeklyTrend(userID, region)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取趋势数据失败"})
		return
	}
	c.JSON(http.StatusOK, trend)
}

func (h *CarbonHandler) GetMonthlyTrend(c *gin.Context) {
	userID := c.GetUint("userID")
	region := c.Query("region")
	if region == "" {
		region = "全国"
	}
	trend, err := h.service.GetMonthlyTrend(userID, region)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取趋势数据失败"})
		return
	}
	c.JSON(http.StatusOK, trend)
}

func (h *CarbonHandler) GetCategorySummary(c *gin.Context) {
	userID := c.GetUint("userID")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	if startDate == "" || endDate == "" {
		now := time.Now()
		endDate = now.Format("2006-01-02")
		startDate = now.AddDate(0, 0, -29).Format("2006-01-02")
	}

	summaries, err := h.service.GetCategorySummary(userID, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取分类统计失败"})
		return
	}
	c.JSON(http.StatusOK, summaries)
}

func (h *CarbonHandler) DeleteRecord(c *gin.Context) {
	userID := c.GetUint("userID")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的记录ID"})
		return
	}

	if err := h.service.DeleteRecord(userID, uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}
