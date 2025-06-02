package role

import "fmt"

type Repo interface {
	FindAll() ([]*Role, error)
	FindById(id int64) (*Role, error)
	FindByIds(ids []int64) ([]*Role, error)
	Create(role *Role) error
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
	role, err := s.repo.FindById(id)
	if err != nil {
		return Response{}, fmt.Errorf("error finding role with id %d: %w", id, err)
	}

	return *role.ToResponse(), nil
}

func (s *Service) FindAll() ([]Response, error) {
	roles, err := s.repo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("error finding all roles: %w", err)
	}

	var responses []Response
	for _, role := range roles {
		responses = append(responses, *role.ToResponse())
	}

	return responses, nil
}

func (s *Service) FindByIds(ids []int64) ([]Response, error) {
	roles, err := s.repo.FindByIds(ids)
	if err != nil {
		return nil, fmt.Errorf("error finding roles with ids %v: %w", ids, err)
	}

	var responses []Response
	for _, role := range roles {
		responses = append(responses, *role.ToResponse())
	}

	return responses, nil
}

func (s *Service) Create(name string) (Response, error) {
	role := &Role{Name: name}
	err := s.repo.Create(role)
	if err != nil {
		return Response{}, fmt.Errorf("error creating role: %w", err)
	}

	return *role.ToResponse(), nil
}

func (s *Service) Remove(id int64) error {
	return s.repo.Remove(id)
}

func (s *Service) RemoveByIds(ids []int64) error {
	return s.repo.RemoveByIds(ids)
}
