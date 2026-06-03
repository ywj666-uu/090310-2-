package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"carbon-tracker/models"
	"carbon-tracker/services"
)

type GoalHandler struct {
	service *services.GoalService
}

func NewGoalHandler(service *services.GoalService) *GoalHandler {
	return &GoalHandler{service: service}
}

func (h *GoalHandler) CreateGoal(c *gin.Context) {
	userID := c.GetUint("userID")
	var req models.GoalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数无效: " + err.Error()})
		return
	}

	goal, err := h.service.CreateGoal(userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "创建目标失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, goal)
}

func (h *GoalHandler) GetActiveGoal(c *gin.Context) {
	userID := c.GetUint("userID")
	status, err := h.service.GetActiveGoal(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "暂无活跃目标"})
		return
	}
	c.JSON(http.StatusOK, status)
}

func (h *GoalHandler) GetGoalHistory(c *gin.Context) {
	userID := c.GetUint("userID")
	goals, err := h.service.GetGoalHistory(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取目标历史失败"})
		return
	}
	c.JSON(http.StatusOK, goals)
}
