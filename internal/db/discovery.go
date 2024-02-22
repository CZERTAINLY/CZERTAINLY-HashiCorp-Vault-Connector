package db

import (
	"gorm.io/gorm"
    "gorm.io/datatypes"
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

func CreateDiscoveryAndAssociateCertificates(db *gorm.DB, discovery *Discovery, certificates ...*Certificate) {
	db.Create(&discovery)
	for _, certificate := range certificates {
		db.FirstOrCreate(&certificate, Certificate{SerialNumber: certificate.SerialNumber})
		db.Model(&discovery).Association("Certificates").Append(&certificate)
	}
}

func FindDiscoveryByUUID(db *gorm.DB, uuid string) (*Discovery, error) {
	var discovery Discovery
	if err := db.Preload("Certificates").First(&discovery, "uuid = ?", uuid).Error; err != nil {
		return nil, err
	}
	return &discovery, nil
}



