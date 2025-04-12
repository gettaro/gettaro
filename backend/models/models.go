package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Auth0ID     string `gorm:"uniqueIndex"`
	Email       string `gorm:"uniqueIndex"`
	Name        string
	Role        string // admin, manager, engineer
	TeamID      uint
	Team        Team
	WorkLogs    []WorkLog
	Performance []PerformanceMetric
}

type Team struct {
	gorm.Model
	Name        string    `gorm:"size:255;not null" json:"name"`
	Description string    `gorm:"size:1000" json:"description"`
	Users       []User    `gorm:"many2many:team_users;" json:"users"`
	Projects    []Project `gorm:"foreignKey:TeamID" json:"projects"`
	Metrics     []TeamMetric
}

type Project struct {
	gorm.Model
	Name        string
	Description string
	TeamID      uint
	Team        Team
	Status      string // active, completed, on-hold
	StartDate   time.Time
	EndDate     time.Time
	Tasks       []Task
	Metrics     []ProjectMetric
}

type Task struct {
	gorm.Model
	Title       string
	Description string
	ProjectID   uint
	Project     Project
	AssigneeID  uint
	Assignee    User
	Status      string // todo, in-progress, review, done
	Priority    string // low, medium, high
	Estimate    int    // in hours
	WorkLogs    []WorkLog
}

type WorkLog struct {
	gorm.Model
	UserID    uint
	User      User
	TaskID    uint
	Task      Task
	StartTime time.Time
	EndTime   time.Time
	Duration  int    // in minutes
	Type      string // coding, review, meeting, etc.
}

type PerformanceMetric struct {
	gorm.Model
	UserID         uint
	User           User
	Period         string // weekly, monthly, quarterly
	StartDate      time.Time
	EndDate        time.Time
	TasksCompleted int
	CodeCommits    int
	PRsMerged      int
	ReviewCount    int
	Velocity       float64
}

type TeamMetric struct {
	gorm.Model
	TeamID         uint
	Team           Team
	Period         string // weekly, monthly, quarterly
	StartDate      time.Time
	EndDate        time.Time
	TotalTasks     int
	CompletedTasks int
	Velocity       float64
	LeadTime       float64
	CycleTime      float64
}

type ProjectMetric struct {
	gorm.Model
	ProjectID      uint
	Project        Project
	Period         string // weekly, monthly, quarterly
	StartDate      time.Time
	EndDate        time.Time
	TotalTasks     int
	CompletedTasks int
	Velocity       float64
	LeadTime       float64
	CycleTime      float64
}

type GenAIUsage struct {
	gorm.Model
	UserID    uint
	User      User
	TaskID    uint
	Task      Task
	Timestamp time.Time
	Tool      string // GitHub Copilot, etc.
	Action    string // suggestion, accept, reject
	Duration  int    // in seconds
}
