package service

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type OperationLog struct {
	ID uint `gorm:"primarykey" json:"-"`

	UserID                uint           `gorm:"index;not null" json:"user_id"`
	ResourceType          string         `gorm:"type:ENUM('user', 'asset');not null" json:"resource_type"`
	ResourceID            uint           `gorm:"index;not null" json:"resource_id"`
	Action                string         `gorm:"type:ENUM('create', 'update', 'delete');not null" json:"action"`
	ActionDetails         string         `gorm:"type:TEXT;not null" json:"action_details"`
	AdditionalInformation datatypes.JSON `json:"additional_information"`

	CreatedAt time.Time `json:"created_at"`
}

func (o *OperationLog) TableName() string {
	return "operation_log"
}

type OperationLogService struct {
	db *gorm.DB
}

func NewOperationLogService(database *gorm.DB) *OperationLogService {
	return &OperationLogService{
		db: database,
	}
}

func (s *OperationLogService) CreateOperationLog(operation_log *OperationLog) error {
	return s.db.Create(operation_log).Error
}

func (s *OperationLogService) GetOperationLogsByUserId(user_id uint) ([]OperationLog, error) {
	var operation_logs []OperationLog
	err := s.db.Where("user_id = ?", user_id).Find(&operation_logs).Error
	return operation_logs, err
}

func (s *OperationLogService) GetOperationLogsByResourceId(resource_type string, resource_id uint) ([]OperationLog, error) {
	var operation_logs []OperationLog
	err := s.db.Where("resource_type = ? AND resource_id = ?", resource_type, resource_id).Find(&operation_logs).Error
	return operation_logs, err
}

func (s *OperationLogService) GetOperationLogsByUserIdAndResourceId(user_id uint, resource_type string, resource_id uint) ([]OperationLog, error) {
	var operation_logs []OperationLog
	err := s.db.Where("user_id = ? AND resource_type = ? AND resource_id = ?", user_id, resource_type, resource_id).Find(&operation_logs).Error
	return operation_logs, err
}

func (s *OperationLogService) GetOperationLogsByUserIdAndResourceType(user_id uint, resource_type string) ([]OperationLog, error) {
	var operation_logs []OperationLog
	err := s.db.Where("user_id = ? AND resource_type = ?", user_id, resource_type).Find(&operation_logs).Error
	return operation_logs, err
}
