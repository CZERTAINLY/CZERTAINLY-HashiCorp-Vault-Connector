package db

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Discovery struct {
	UUID         string
	Name         string
	Status       string
	Meta         datatypes.JSON
	Certificates []Certificate `gorm:"many2many:discovery_certificates;"`
}

type Certificate struct {
	SerialNumber  string
	UUID          string
	Base64Content string
	Meta          datatypes.JSON
	Discoveries   []Discovery `gorm:"many2many:discovery_certificates;"`
}

type DiscoveryRepository struct {
	db *gorm.DB
}

func NewDiscoveryRepository(db *gorm.DB) (*DiscoveryRepository, error) {
	return &DiscoveryRepository{db: db}, nil
}


func (d *DiscoveryRepository) CreateDiscoveryAndAssociateCertificates(discovery *Discovery, certificates ...*Certificate) {
	d.db.Create(&discovery)
	for _, certificate := range certificates {
		d.db.FirstOrCreate(&certificate, Certificate{SerialNumber: certificate.SerialNumber})
		d.db.Model(&discovery).Association("Certificates").Append(&certificate)
	}
}

func (d *DiscoveryRepository) FindDiscoveryByUUID(uuid string) (*Discovery, error) {
	var discovery Discovery
	if err := d.db.Preload("Certificates").First(&discovery, "uuid = ?", uuid).Error; err != nil {
		return nil, err
	}
	return &discovery, nil
}
