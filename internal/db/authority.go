package db

import (
	"gorm.io/gorm"
)

type AuthorityInstance struct {
	ID             int64  `db:"id"`
	UUID           string `db:"uuid"`
	Name           string `db:"name"`
	URL            string `db:"url"`
	CredentialType string `db:"credential_type"`
	RoleId         string `db:"role_id"`
	RoleSecret     string `db:"role_secret"`
	VaultRole      string `db:"vault_role"`
	MountPath      string `db:"login_mount_path"`
	Attributes     string `db:"attributes"`
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

func (d *AuthorityRepository) DeleteAuthorityInstanceByUUID(uuid string) error {
	var authority AuthorityInstance
	err := d.db.Where("uuid = ?", uuid).First(&authority).Error
	if err != nil {
		return err
	}
	return d.db.Delete(&authority).Error
}

func (d *AuthorityRepository) FindAuthorityInstanceByUUID(uuid string) (*AuthorityInstance, error) {
	var authority AuthorityInstance
	err := d.db.Where("uuid = ?", uuid).First(&authority).Error
	if err != nil {
		return nil, err
	}
	return &authority, nil
}

func (d *AuthorityRepository) ListAuthorityInstances() ([]*AuthorityInstance, error) {
	var authorities []*AuthorityInstance
	err := d.db.Find(&authorities).Error
	if err != nil {
		return nil, err
	}
	return authorities, nil
}
