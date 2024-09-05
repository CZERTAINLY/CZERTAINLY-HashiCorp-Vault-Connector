package db

import (
	"math"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Discovery struct {
	Id           uint `gorm:"primarykey"`
	UUID         string
	Name         string
	Status       string
	Meta         datatypes.JSON
	Certificates []Certificate `gorm:"many2many:discovery_certificates;"`
}

type Certificate struct {
	Id            uint `gorm:"primarykey"`
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

func (d *DiscoveryRepository) CreateDiscovery(discovery *Discovery) error {
	result := d.db.Create(&discovery)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (d *DiscoveryRepository) AssociateCertificatesToDiscovery(discovery *Discovery, certificates ...*Certificate) error {
	for _, certificate := range certificates {
		d.db.FirstOrCreate(&certificate, Certificate{SerialNumber: certificate.SerialNumber})
	}
	assoc := d.db.Model(&discovery).Association("Certificates")
	err := assoc.Append(&certificates)
	if err != nil {
		return err
	}
	return nil
}

func (d *DiscoveryRepository) CreateDiscoveryAndAssociateCertificates(discovery *Discovery, certificates ...*Certificate) error {
	d.db.Create(&discovery)
	for _, certificate := range certificates {
		d.db.FirstOrCreate(&certificate, Certificate{SerialNumber: certificate.SerialNumber})
		assoc := d.db.Model(&discovery).Association("Certificates")
		err := assoc.Append(&certificate)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *DiscoveryRepository) FindDiscoveryByUUID(uuid string) (*Discovery, error) {
	var discovery Discovery
	if err := d.db.Preload("Certificates").First(&discovery, "uuid = ?", uuid).Error; err != nil {
		return nil, err
	}
	return &discovery, nil
}

func (d *DiscoveryRepository) UpdateDiscovery(discovery *Discovery) error {
	result := d.db.Save(&discovery)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (d *DiscoveryRepository) List(pagination Pagination, discovery *Discovery) (*Pagination, error) {
	var certificates []*Certificate
	page := pagination.Page
	pageSize := pagination.Limit
	var count int
	offset := (page - 1) * pageSize
	d.db.Table("certificates").Select("certificates.*").
		Joins("JOIN discovery_certificates ON discovery_certificates.certificate_id = certificates.id").
		Where("discovery_certificates.discovery_id = ?", discovery.Id).
		Offset(offset).Limit(pageSize).
		Order("certificates.id ASC").
		Find(&certificates)

	pagination.Rows = certificates
	pagination.TotalRows = int64(len(discovery.Certificates))
	totalPages := int(math.Ceil(float64(count) / float64(pagination.Limit)))
	pagination.TotalPages = totalPages
	return &pagination, nil
}

func (d *DiscoveryRepository) DeleteDiscovery(discovery *Discovery) error {
	result := d.db.Delete(&discovery)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

type Pagination struct {
	Limit      int         `json:"limit,omitempty"`
	Page       int         `json:"page,omitempty"`
	Sort       string      `json:"sort,omitempty"`
	TotalRows  int64       `json:"total_rows"`
	TotalPages int         `json:"total_pages"`
	Rows       interface{} `json:"rows"`
}

func (p *Pagination) GetOffset() int {
	return (p.GetPage() - 1) * p.GetLimit()
}

func (p *Pagination) GetLimit() int {
	if p.Limit == 0 {
		p.Limit = 10
	}
	return p.Limit
}

func (p *Pagination) GetPage() int {
	if p.Page == 0 {
		p.Page = 1
	}
	return p.Page
}
