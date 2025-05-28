package employee

import (
	"idm/inner/employee"
)

type Fixture struct {
	employees *employee.Repository
}

func NewFixture(employees *employee.Repository) *Fixture {
	return &Fixture{employees: employees}
}

func (f *Fixture) FindById(id int64) (*employee.Employee, error) {
	return f.employees.FindById(id)
}

func (f *Fixture) FindByIds(ids []int64) ([]*employee.Employee, error) {
	return f.employees.FindByIds(ids)
}

func (f *Fixture) FindAll() ([]*employee.Employee, error) {
	return f.employees.FindAll()
}

func (f *Fixture) CreateEmployee(name string) (*employee.Employee, error) {
	empl := &employee.Employee{Name: name}
	err := f.employees.Create(empl)

	return empl, err
}

func (f *Fixture) Remove(id int64) error {
	return f.employees.Remove(id)
}

func (f *Fixture) RemoveByIds(ids []int64) error {
	return f.employees.RemoveByIds(ids)
}
