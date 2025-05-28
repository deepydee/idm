package role

import "idm/inner/role"

type Fixture struct {
	roles *role.Repository
}

func NewFixture(roles *role.Repository) *Fixture {
	return &Fixture{roles: roles}
}

func (f *Fixture) FindAll() ([]*role.Role, error) {
	return f.roles.FindAll()
}

func (f *Fixture) FindById(id int64) (*role.Role, error) {
	return f.roles.FindById(id)
}

func (f *Fixture) FindByIds(ids []int64) ([]*role.Role, error) {
	return f.roles.FindByIds(ids)
}

func (f *Fixture) Create(name string) (*role.Role, error) {
	roleEntity := &role.Role{Name: name}
	err := f.roles.Create(roleEntity)

	return roleEntity, err
}

func (f *Fixture) Remove(id int64) error {
	return f.roles.Remove(id)
}

func (f *Fixture) RemoveByIds(ids []int64) error {
	return f.roles.RemoveByIds(ids)
}
