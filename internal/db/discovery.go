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
	var certificates []Certificate
	page := pagination.Page
	pageSize := pagination.Limit
	offset := (page - 1) * pageSize

	pagination.TotalRows = d.db.Model(discovery).Association("Certificates").Count()
	if pagination.TotalRows == 0 {
		pagination.Rows = []Certificate{}
		pagination.TotalPages = 0
		return &pagination, nil
	}

	tempDiscovery := Discovery{Id: discovery.Id}
	err := d.db.Preload("Certificates", func(db *gorm.DB) *gorm.DB {
		return db.Order("id ASC").Offset(offset).Limit(pageSize)
	}).First(&tempDiscovery).Error

	if err != nil {
		return nil, err
	}

	certificates = tempDiscovery.Certificates
	pagination.Rows = certificates
	totalPages := int(math.Ceil(float64(pagination.TotalRows) / float64(pagination.Limit)))
	pagination.TotalPages = totalPages
	return &pagination, nil
}

func (d *DiscoveryRepository) DeleteDiscovery(discovery *Discovery) error {
	return d.db.Transaction(func(tx *gorm.DB) error {
		// 1. Clear the many-to-many associations in the join table (discovery_certificates)
		// This will NOT delete the actual Certificate records.
		err := tx.Model(discovery).Association("Certificates").Clear()
		if err != nil {
			return err
		}

		// 2. Delete the Discovery record itself.
		// Use discovery directly as it is already a *Discovery.
		return tx.Delete(discovery).Error
	})
}

type Pagination struct {
	Limit      int    `json:"limit,omitempty"`
	Page       int    `json:"page,omitempty"`
	Sort       string `json:"sort,omitempty"`
	TotalRows  int64  `json:"total_rows"`
	TotalPages int    `json:"total_pages"`
	Rows       any    `json:"rows"`
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
