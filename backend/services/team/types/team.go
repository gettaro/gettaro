package types

type Team struct {
	ID          string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name        string
	Description string
}

type TeamMember struct {
	ID     string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID string
	TeamID string
}

type DirectReport struct {
	ID        string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	ManagerID string
	ReportID  string
	Depth     int

	// Unique constraint
	UniqueManagerReport string `gorm:"uniqueIndex:idx_manager_report"`
}
