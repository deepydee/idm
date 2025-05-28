package employee

import "fmt"

type Repo interface {
	FindById(id int64) (*Employee, error)
	FindAll() ([]*Employee, error)
	FindByIds(ids []int64) ([]*Employee, error)
	Create(employee *Employee) error
	Remove(id int64) error
	RemoveByIds(ids []int64) error
}

// Service будет инкапсулировать бизнес-логику
type Service struct {
	repo Repo
}

func NewService(repository Repo) *Service {
	return &Service{repo: repository}
}

func (s *Service) FindById(id int64) (Response, error) {
	employee, err := s.repo.FindById(id)
	if err != nil {
		return Response{}, fmt.Errorf("error finding employee with id %d: %w", id, err)
	}

	return *employee.ToResponse(), nil
}

func (s *Service) FindAll() ([]Response, error) {
	employees, err := s.repo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("error finding all employees: %w", err)
	}

	var responses []Response
	for _, employee := range employees {
		responses = append(responses, *employee.ToResponse())
	}

	return responses, nil
}

func (s *Service) FindByIds(ids []int64) ([]Response, error) {
	employees, err := s.repo.FindByIds(ids)
	if err != nil {
		return nil, fmt.Errorf("error finding employees with ids %v: %w", ids, err)
	}

	var responses []Response
	for _, employee := range employees {
		responses = append(responses, *employee.ToResponse())
	}

	return responses, nil
}

func (s *Service) Create(name string) (Response, error) {
	employee := &Employee{Name: name}
	err := s.repo.Create(employee)
	if err != nil {
		return Response{}, fmt.Errorf("error creating employee: %w", err)
	}

	return *employee.ToResponse(), nil
}

func (s *Service) Remove(id int64) error {
	return s.repo.Remove(id)
}

func (s *Service) RemoveByIds(ids []int64) error {
	return s.repo.RemoveByIds(ids)
}
