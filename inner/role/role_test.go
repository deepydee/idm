package role

import (
	"errors"
	"fmt"
	assertpackage "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type MockRepo struct {
	mock.Mock
}

func (m *MockRepo) FindById(id int64) (*Role, error) {
	args := m.Called(id)
	return args.Get(0).(*Role), args.Error(1)
}

func (m *MockRepo) FindAll() ([]*Role, error) {
	args := m.Called()
	return args.Get(0).([]*Role), args.Error(1)
}

func (m *MockRepo) FindByIds(ids []int64) ([]*Role, error) {
	args := m.Called(ids)
	return args.Get(0).([]*Role), args.Error(1)
}

func (m *MockRepo) Create(role *Role) error {
	args := m.Called(role)
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

func TestRoleService(t *testing.T) {
	assert := assertpackage.New(t)

	t.Run("FindById should return a role", func(t *testing.T) {
		repo := &MockRepo{}
		service := NewService(repo)

		repo.On("FindById", int64(1)).Return(&Role{Name: "admin"}, nil)
		role, err := service.FindById(1)

		assert.NoError(err)
		assert.Equal("admin", role.Name)
	})

	t.Run("FindById should return wrapped error", func(t *testing.T) {
		repo := &MockRepo{}
		service := NewService(repo)
		role := Role{}

		err := errors.New("database error")
		want := fmt.Errorf("error finding role with id %d: %w", 1, err)

		repo.On("FindById", int64(1)).Return(&role, err)
		resp, got := service.FindById(1)

		assert.Empty(resp)
		assert.NotNil(got)
		assert.Equal(want, got)
		assert.True(repo.AssertNumberOfCalls(t, "FindById", 1))
	})

	t.Run("FindAll should return a list of roles", func(t *testing.T) {
		repo := &MockRepo{}
		service := NewService(repo)

		repo.On("FindAll").Return([]*Role{{Name: "admin"}, {Name: "user"}}, nil)
		roles, err := service.FindAll()

		assert.NoError(err)
		assert.Len(roles, 2)
		assert.Equal("admin", roles[0].Name)
		assert.Equal("user", roles[1].Name)
		assert.True(repo.AssertNumberOfCalls(t, "FindAll", 1))
	})

	t.Run("FindByIds should return a list of roles", func(t *testing.T) {
		repo := &MockRepo{}
		service := NewService(repo)

		repo.On("FindByIds", []int64{1, 2}).Return([]*Role{{Name: "admin"}, {Name: "user"}}, nil)
		roles, err := service.FindByIds([]int64{1, 2})

		assert.NoError(err)
		assert.Len(roles, 2)
		assert.Equal("admin", roles[0].Name)
		assert.Equal("user", roles[1].Name)
		assert.True(repo.AssertNumberOfCalls(t, "FindByIds", 1))
	})

	t.Run("Create should create a role", func(t *testing.T) {
		repo := &MockRepo{}
		service := NewService(repo)

		repo.On("Create", &Role{Name: "admin"}).Return(nil)
		role, err := service.Create("admin")

		assert.NoError(err)
		assert.NotNil(role)
		assert.Equal("admin", role.Name)
		assert.True(repo.AssertNumberOfCalls(t, "Create", 1))
	})

	t.Run("Create should return wrapped error", func(t *testing.T) {
		repo := &MockRepo{}
		service := NewService(repo)

		err := errors.New("database error")
		want := fmt.Errorf("error creating role: %w", err)

		repo.On("Create", &Role{Name: "admin"}).Return(err)
		resp, got := service.Create("admin")

		assert.Empty(resp)
		assert.NotNil(got)
		assert.Equal(want, got)
		assert.True(repo.AssertNumberOfCalls(t, "Create", 1))
	})

	t.Run("Remove should remove a role", func(t *testing.T) {
		repo := &MockRepo{}
		service := NewService(repo)

		repo.On("Remove", int64(1)).Return(nil)
		err := service.Remove(1)

		assert.NoError(err)
		assert.True(repo.AssertNumberOfCalls(t, "Remove", 1))
	})

	t.Run("Remove by ids should remove a list of roles", func(t *testing.T) {
		repo := &MockRepo{}
		service := NewService(repo)

		repo.On("RemoveByIds", []int64{1, 2}).Return(nil)
		err := service.RemoveByIds([]int64{1, 2})

		assert.NoError(err)
		assert.True(repo.AssertNumberOfCalls(t, "RemoveByIds", 1))
	})
}
