package db

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type AuthorityInstance struct {
	ID             int64          `db:"id"`
	UUID           string         `db:"uuid"`
	Name           *string        `db:"name"`
	URL            *string        `db:"url"`
	CredentialUUID *string        `db:"credential_uuid"`
	CredentialData *string        `db:"credential_data"`
	Attributes     datatypes.JSON `db:"attributes"`
}

type AuthorityRepository struct {
	db *gorm.DB
}
func NewAuthorityRepository(db *gorm.DB) (*AuthorityRepository, error) {
	return &AuthorityRepository{db: db}, nil
}

func (d *AuthorityRepository) CreateAuthorityInstance(authority *AuthorityInstance) error {
	return d.db.Create(authority).Error
}

func (d *AuthorityRepository) UpdateAuthorityInstance(authority *AuthorityInstance) error {
	return d.db.Save(authority).Error
}

func (d *AuthorityRepository) DeleteAuthorityInstance(authority *AuthorityInstance) error {
	return d.db.Delete(authority).Error
}

func (d *AuthorityRepository) FindAuthorityInstanceByName(name string) (*AuthorityInstance, error) {
	var authority AuthorityInstance
	err := d.db.Where("name = ?", name).First(&authority).Error
	if err != nil {
		return nil, err
	}
	return &authority, nil
}
