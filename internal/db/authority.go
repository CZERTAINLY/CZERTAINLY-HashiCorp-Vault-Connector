package db

import (
	"gorm.io/gorm"
	"gorm.io/datatypes"
)

type AuthorityInstance struct {
	ID             int64   `db:"id"`
	UUID           string  `db:"uuid"`
	Name           *string `db:"name"`
	URL            *string `db:"url"`
	CredentialUUID *string `db:"credential_uuid"`
	CredentialData *string `db:"credential_data"`
	Attributes     datatypes.JSON `db:"attributes"`
}

func CreateAuthorityInstance(db *gorm.DB, authority *AuthorityInstance) error {
	return db.Create(authority).Error
}

func UpdateAuthorityInstance(db *gorm.DB, authority *AuthorityInstance) error {
	return db.Save(authority).Error
}

func DeleteAuthorityInstance(db *gorm.DB, authority *AuthorityInstance) error {
	return db.Delete(authority).Error
}

func FindAuthorityInstanceByName(db *gorm.DB, name string) (*AuthorityInstance, error) {
	var authority AuthorityInstance
	err := db.Where("name = ?", name).First(&authority).Error
	if err != nil {
		return nil, err
	}
	return &authority, nil
}
