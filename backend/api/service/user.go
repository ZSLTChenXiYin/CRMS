package service

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	Email        string `gorm:"unique"`
	PasswordHash string `gorm:"not null"`

	ExpiredAt   time.Time
	LastLoginAt time.Time
}

func (u *User) TableName() string {
	return "users"
}

type UserService struct {
	db *gorm.DB
}

func NewUserService(database *gorm.DB) *UserService {
	return &UserService{
		db: database,
	}
}

func (s *UserService) CreateUser(email string, password_hash string) error {
	return s.db.Create(&User{
		Email:        email,
		PasswordHash: password_hash,
		ExpiredAt:    time.Now(),
		LastLoginAt:  time.Now(),
	}).Error
}

func (s *UserService) GetUserByID(id uint) (*User, error) {
	user := &User{}
	err := s.db.First(user, id).Error
	return user, err
}

func (s *UserService) GetUserByEmail(email string) (*User, error) {
	user := &User{}
	err := s.db.Where("email = ?", email).First(user).Error
	return user, err
}

func (s *UserService) DeleteUserById(id uint) error {
	return s.db.Delete(&User{}, id).Error
}

func (s *UserService) UpdateUserPasswordHashById(id uint, password_hash string) error {
	return s.db.Model(&User{}).Where("id = ?", id).Update("password_hash", password_hash).Error
}

func (s *UserService) UpdateUserExpiredById(id uint, expired_at time.Time) error {
	return s.db.Model(&User{}).Where("id = ?", id).Update("expired_at", expired_at).Error
}

func (s *UserService) UpdateUserLastLoginById(id uint) error {
	return s.db.Model(&User{}).Where("id = ?", id).Update("last_login_at", time.Now()).Error
}

type UserAssetMapping struct {
	gorm.Model

	UserID  uint `gorm:"index:idx_user_asset,unique;not null"`
	AssetID uint `gorm:"index:idx_user_asset,unique;not null"`

	Permission string `gorm:"type:ENUM('use', 'execute');not null"`
}

func (u *UserAssetMapping) TableName() string {
	return "user_asset_mapping"
}

type UserAssetService struct {
	db *gorm.DB
}

func NewUserAssetService(database *gorm.DB) *UserAssetService {
	return &UserAssetService{
		db: database,
	}
}

func (s *UserAssetService) CreateUserAssetMapping(user_id uint, asset_id uint, permission string) error {
	return s.db.Create(&UserAssetMapping{
		UserID:     user_id,
		AssetID:    asset_id,
		Permission: permission,
	}).Error
}

func (s *UserAssetService) GetUserAssetMappingsByUserId(user_id uint) ([]UserAssetMapping, error) {
	var mappings []UserAssetMapping
	err := s.db.Where("user_id = ?", user_id).Find(&mappings).Error
	return mappings, err
}

func (s *UserAssetService) GetUserAssetMappingsByAssetId(asset_id uint) ([]UserAssetMapping, error) {
	var mappings []UserAssetMapping
	err := s.db.Where("asset_id = ?", asset_id).Find(&mappings).Error
	return mappings, err
}

func (s *UserAssetService) DeleteUserAssetMappingByAssetId(asset_id uint) error {
	return s.db.Where("asset_id = ?", asset_id).Delete(&UserAssetMapping{}).Error
}

func (s *UserAssetService) DeleteUserAssetMappingByUserIdAndAssetId(user_id uint, asset_id uint) error {
	return s.db.Where("user_id = ? AND asset_id = ?", user_id, asset_id).Delete(&UserAssetMapping{}).Error
}

func (s *UserAssetService) UpdateUserAssetMappingPermissionByUserIDAndAssetID(user_id uint, asset_id uint, permission string) error {
	return s.db.Model(&UserAssetMapping{}).
		Where("user_id = ? AND asset_id = ?", user_id, asset_id).
		Update("permission", permission).Error
}
