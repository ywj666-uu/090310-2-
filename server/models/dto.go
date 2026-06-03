package models

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6"`
	Email    string `json:"email"`
	Region   string `json:"region"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type RecordRequest struct {
	RecordDate string  `json:"record_date" binding:"required"`
	Category   string  `json:"category" binding:"required"`
	Item       string  `json:"item" binding:"required"`
	Amount     float64 `json:"amount" binding:"required,gt=0"`
	Note       string  `json:"note"`
}

type GoalRequest struct {
	TargetEmission float64 `json:"target_emission" binding:"required,gt=0"`
	StartDate      string  `json:"start_date" binding:"required"`
	EndDate        string  `json:"end_date" binding:"required"`
}

type DailySummary struct {
	Date          string  `json:"date"`
	TotalEmission float64 `json:"total_emission"`
}

type CategorySummary struct {
	Category      string  `json:"category"`
	TotalEmission float64 `json:"total_emission"`
}

type TrendData struct {
	UserData        []DailySummary `json:"user_data"`
	RegionalAverage float64        `json:"regional_average"`
}

type GoalStatus struct {
	Goal            ReductionGoal `json:"goal"`
	CurrentAvg      float64       `json:"current_avg"`
	Achieved        bool          `json:"achieved"`
	Encouragement   string        `json:"encouragement"`
	DaysRemaining   int           `json:"days_remaining"`
	CompletionRate  float64       `json:"completion_rate"`
}
