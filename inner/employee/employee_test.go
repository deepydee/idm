package employee

import (
	"errors"
	"fmt"
	assertpackage "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

type MockRepo struct {
	mock.Mock
}

// StubRepo - это заглушка (stub) репозитория
type StubRepo struct {
	employees map[int64]*Employee
}

// NewStubRepo создает новую заглушку репозитория с предварительно заполненными данными
func NewStubRepo() *StubRepo {
	return &StubRepo{
		employees: map[int64]*Employee{
			1: {
				Id:        1,
				Name:      "John Doe",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			2: {
				Id:        2,
				Name:      "Jane Smith",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		},
	}
}

// Реализация метода FindById для заглушки
func (s *StubRepo) FindById(id int64) (*Employee, error) {
	employee, exists := s.employees[id]
	if !exists {
		return nil, errors.New("employee not found")
	}
	return employee, nil
}

// Другие методы репозитория могут быть также реализованы при необходимости
func (s *StubRepo) FindAll() ([]*Employee, error) {
	return nil, nil
}

func (s *StubRepo) FindByIds(ids []int64) ([]*Employee, error) {
	return nil, nil
}

func (s *StubRepo) Create(employee *Employee) error {
	return nil
}

func (s *StubRepo) Remove(id int64) error {
	return nil
}

func (s *StubRepo) RemoveByIds(ids []int64) error {
	return nil
}

func (m *MockRepo) FindById(id int64) (*Employee, error) {
	args := m.Called(id)
	return args.Get(0).(*Employee), args.Error(1)
}

func (m *MockRepo) FindAll() ([]*Employee, error) {
	args := m.Called()
	return args.Get(0).([]*Employee), args.Error(1)
}

func (m *MockRepo) FindByIds(ids []int64) ([]*Employee, error) {
	args := m.Called(ids)
	return args.Get(0).([]*Employee), args.Error(1)
}

func (m *MockRepo) Create(employee *Employee) error {
	args := m.Called(employee)
	return args.Error(0)
}

func (m *MockRepo) Remove(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockRepo) RemoveByIds(ids []int64) error {
	args := m.Called(ids)
	return args.Error(0)
}

func TestEmployeeService(t *testing.T) {
	assert := assertpackage.New(t)

	t.Run("FindById should return an employee (use stub)", func(t *testing.T) {
		stubRepo := NewStubRepo()
		service := NewService(stubRepo)

		response, err := service.FindById(1)

		assert.Nil(err)
		assert.NotNil(response)
		assert.Equal(int64(1), response.Id)
		assert.Equal("John Doe", response.Name)
	})

	t.Run("FindById should return an employee", func(t *testing.T) {
		repo := &MockRepo{}
		service := NewService(repo)

		employee := Employee{
			Id:        1,
			Name:      "John Doe",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		want := employee.ToResponse()

		repo.On("FindById", int64(1)).Return(&employee, nil)
		got, err := service.FindById(1)

		assert.Nil(err)
		assert.Equal(*want, got)
		assert.True(repo.AssertNumberOfCalls(t, "FindById", 1))
	})

	t.Run("FindById should return wrapped error", func(t *testing.T) {
		repo := &MockRepo{}
		service := NewService(repo)
		employee := Employee{}

		err := errors.New("database error")
		want := fmt.Errorf("error finding employee with id 1: %w", err)

		repo.On("FindById", int64(1)).Return(&employee, err)
		resp, got := service.FindById(1)

		assert.Empty(resp)
		assert.NotNil(got)
		assert.Equal(want, got)
		assert.True(repo.AssertNumberOfCalls(t, "FindById", 1))
	})

	t.Run("FindAll should return a list of employees", func(t *testing.T) {
		repo := &MockRepo{}
		service := NewService(repo)

		repo.On("FindAll").Return([]*Employee{{Name: "John"}, {Name: "Jane"}}, nil)
		got, err := service.FindAll()

		assert.Nil(err)
		assert.Equal(2, len(got))
		assert.Equal("John", got[0].Name)
		assert.Equal("Jane", got[1].Name)
		assert.True(repo.AssertNumberOfCalls(t, "FindAll", 1))
	})

	t.Run("FindAll should return error", func(t *testing.T) {
		repo := &MockRepo{}
		service := NewService(repo)

		err := errors.New("database error")
		want := fmt.Errorf("error finding all employees: %w", err)

		repo.On("FindAll").Return([]*Employee{}, err)

		responses, got := service.FindAll()

		assert.Empty(responses)
		assert.NotNil(got)
		assert.Equal(want, got)
		assert.True(repo.AssertNumberOfCalls(t, "FindAll", 1))
	})

	t.Run("FindByIds should return a list of employees", func(t *testing.T) {
		repo := &MockRepo{}
		service := NewService(repo)

		employees := []*Employee{
			{
				Id:        1,
				Name:      "John",
				CreatedAt: time.Now(),
			},
			{
				Id:        2,
				Name:      "Jane",
				CreatedAt: time.Now(),
			},
		}

		var expectedResponses []Response
		for _, emp := range employees {
			expectedResponses = append(expectedResponses, *emp.ToResponse())
		}

		repo.On("FindByIds", []int64{1, 2}).Return(employees, nil)
		got, err := service.FindByIds([]int64{1, 2})

		assert.Nil(err)
		assert.Equal(len(expectedResponses), len(got))
		assert.Equal(expectedResponses, got)
		assert.True(repo.AssertNumberOfCalls(t, "FindByIds", 1))
	})

	t.Run("FindByIds should return error", func(t *testing.T) {
		repo := &MockRepo{}
		service := NewService(repo)

		err := errors.New("database error")
		want := fmt.Errorf("error finding employees with ids %v: %w", []int64{1, 2}, err)

		repo.On("FindByIds", []int64{1, 2}).Return([]*Employee{}, err)

		responses, got := service.FindByIds([]int64{1, 2})

		assert.Empty(responses)
		assert.NotNil(got)
		assert.Equal(want, got)
		assert.True(repo.AssertNumberOfCalls(t, "FindByIds", 1))
	})

	t.Run("Create should create an employee", func(t *testing.T) {
		repo := &MockRepo{}
		service := NewService(repo)

		repo.On("Create", mock.AnythingOfType("*employee.Employee")).Return(nil)
		got, err := service.Create("John")

		assert.Nil(err)
		assert.NotNil(got)
		assert.Equal("John", got.Name)
		assert.True(repo.AssertNumberOfCalls(t, "Create", 1))
	})

	t.Run("Create should return error", func(t *testing.T) {
		repo := &MockRepo{}
		service := NewService(repo)

		err := errors.New("database error")
		want := fmt.Errorf("error creating employee: %w", err)

		repo.On("Create", mock.AnythingOfType("*employee.Employee")).Return(err)

		resp, got := service.Create("John")

		assert.Empty(resp)
		assert.NotNil(got)
		assert.Equal(want, got)
		assert.True(repo.AssertNumberOfCalls(t, "Create", 1))
	})

	t.Run("Remove should remove an employee", func(t *testing.T) {
		repo := &MockRepo{}
		service := NewService(repo)

		repo.On("Remove", int64(1)).Return(nil)
		err := service.Remove(1)

		assert.Nil(err)
		assert.True(repo.AssertNumberOfCalls(t, "Remove", 1))
		assert.True(repo.AssertNumberOfCalls(t, "Remove", 1))
	})

	t.Run("Remove by ids should remove employees", func(t *testing.T) {
		repo := &MockRepo{}
		service := NewService(repo)

		repo.On("RemoveByIds", []int64{1, 2}).Return(nil)
		err := service.RemoveByIds([]int64{1, 2})

		assert.Nil(err)
		assert.True(repo.AssertNumberOfCalls(t, "RemoveByIds", 1))
		assert.True(repo.AssertNumberOfCalls(t, "RemoveByIds", 1))
	})
}

func TestService_Create(t *testing.T) {
	assert := assertpackage.New(t)

	t.Run("unable to begin transaction", func(t *testing.T) {
		repo := &MockRepo{}
		service := NewService(repo)
		baseError := errors.New("error creating transaction")
		repo.On("Create", &Employee{Name: "John Doe"}).Return(baseError)
		_, err := service.Create("John Doe")

		assert.Error(err)
		assert.Contains(err.Error(), "error creating transaction")
		assert.True(errors.Is(err, baseError))
	})

	t.Run("error checking employee exists", func(t *testing.T) {
		repo := &MockRepo{}
		service := NewService(repo)
		baseError := errors.New("error checking employee exists")
		repo.On("Create", &Employee{Name: "John Doe"}).Return(baseError)
		_, err := service.Create("John Doe")

		assert.Error(err)
		assert.Contains(err.Error(), "error checking employee exists")
		assert.True(errors.Is(err, baseError))
	})

	t.Run("employee already exists", func(t *testing.T) {
		repo := &MockRepo{}
		service := NewService(repo)
		baseError := errors.New("employee already exists")
		repo.On("Create", &Employee{Name: "John Doe"}).Return(baseError)
		_, err := service.Create("John Doe")

		assert.Error(err)
		assert.Contains(err.Error(), "employee already exists")
		assert.True(errors.Is(err, baseError))
	})

	t.Run("error creating employee", func(t *testing.T) {
		repo := &MockRepo{}
		service := NewService(repo)
		baseError := errors.New("error creating employee")
		repo.On("Create", &Employee{Name: "John Doe"}).Return(baseError)
		_, err := service.Create("John Doe")

		assert.Error(err)
		assert.Contains(err.Error(), "error creating employee")
		assert.True(errors.Is(err, baseError))
	})

	t.Run("successful employee creation", func(t *testing.T) {
		repo := &MockRepo{}
		service := NewService(repo)
		repo.On("Create", mock.AnythingOfType("*employee.Employee")).Return(nil)
		got, err := service.Create("John")

		assert.Nil(err)
		assert.NotNil(got)
		assert.Equal("John", got.Name)
	})
}
