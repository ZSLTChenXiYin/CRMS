package service

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Asset struct {
	gorm.Model

	Type    string         `gorm:"type:ENUM('server');not null"`
	Name    string         `gorm:"not null"`
	Data    datatypes.JSON `gorm:"not null"`
	OwnerID uint           `gorm:"index;not null"`
}

func (a *Asset) TableName() string {
	return "assets"
}

type AssetService struct {
	db *gorm.DB
}

func NewAssetService(database *gorm.DB) *AssetService {
	return &AssetService{db: database}
}

func (s *AssetService) CreateAsset(asset *Asset) error {
	return s.db.Create(asset).Error
}

func (s *AssetService) GetAssetByOwnerID(owner_id uint) ([]Asset, error) {
	var assets []Asset
	err := s.db.Where("owner_id = ?", owner_id).Find(&assets).Error
	return assets, err
}

func (s *AssetService) GetAssetByID(id uint) (*Asset, error) {
	asset := &Asset{}
	err := s.db.First(asset, id).Error
	return asset, err
}

func (s *AssetService) DeleteAssetByID(id uint) error {
	return s.db.Delete(&Asset{}, id).Error
}

func (s *AssetService) UpdateAssetNameByID(id uint, name string) error {
	return s.db.Model(&Asset{}).Where("id = ?", id).Update("name", name).Error
}

func (s *AssetService) UpdateAssetDataByID(id uint, data datatypes.JSON) error {
	return s.db.Model(&Asset{}).Where("id = ?", id).Update("data", data).Error
}

type Server struct {
	InstanceID   string `json:"instance_id"`
	InstanceName string `json:"instance_name"`

	CPUCoreNum   uint `json:"cpu_core_num"`
	GBMemorySize uint `json:"gb_memory_size"`
	GBDiskSize   uint `json:"gb_disk_size"`
	Mbps         uint `json:"mbps"`

	IPv4Address string `json:"ipv4_address"`
	IPv6Address string `json:"ipv6_address"`

	SSHUsername string `json:"ssh_username"`
	SSHPassword string `json:"ssh_password"`

	PanelAddress  string `json:"panel_address"`
	PanelUsername string `json:"panel_username"`
	PanelPassword string `json:"panel_password"`

	Expire int64 `json:"expire"`
}
