package role

import (
	assertpackage "github.com/stretchr/testify/assert"
	"idm/inner/database"
	"idm/inner/role"
	"testing"
)

func TestRoleRepository(t *testing.T) {
	assert := assertpackage.New(t)
	var db = database.ConnectDb()

	var clearDb = func() {
		db.MustExec("DELETE FROM roles")
	}

	defer func() {
		if r := recover(); r != nil {
			clearDb()
		}
	}()

	var roleRepository = role.NewRepository(db)
	var fixture = NewFixture(roleRepository)

	t.Run("we can create a role", func(t *testing.T) {
		got, err := fixture.Create("Admin")
		if err != nil {
			t.Logf("unexpected error while creating role: %v", err)
		}
		assert.Nil(err)
		assert.NotNil(got)
		assert.NotEmpty(got.Id)
		assert.NotEmpty(got.CreatedAt)
		assert.NotEmpty(got.UpdatedAt)
		assert.Equal("Admin", got.Name)

		clearDb()
	})

	t.Run("we can find a role by id", func(t *testing.T) {
		roleEntity, err := fixture.Create("Admin")
		if err != nil {
			t.Logf("unexpected error while creating role: %v", err)
		}
		got, err := fixture.FindById(roleEntity.Id)
		if err != nil {
			t.Logf("unexpected error while finding role by id: %v", err)
		}
		assert.Nil(err)
		assert.NotNil(got)
		assert.Equal(roleEntity.Id, got.Id)
		assert.Equal(roleEntity.CreatedAt, got.CreatedAt)
		assert.Equal(roleEntity.UpdatedAt, got.UpdatedAt)
		assert.Equal(roleEntity.Name, got.Name)

		clearDb()
	})

	t.Run("we can find all roles", func(t *testing.T) {
		roleEntity, err := fixture.Create("Admin")
		if err != nil {
			t.Logf("unexpected error while creating role: %v", err)
		}
		got, err := fixture.FindAll()
		if err != nil {
			t.Logf("unexpected error while finding all roles: %v", err)
		}
		assert.Nil(err)
		assert.NotNil(got)
		assert.Len(got, 1)
		assert.Equal(roleEntity.Id, got[0].Id)
		assert.Equal(roleEntity.CreatedAt, got[0].CreatedAt)
		assert.Equal(roleEntity.UpdatedAt, got[0].UpdatedAt)
		assert.Equal(roleEntity.Name, got[0].Name)

		clearDb()
	})

	t.Run("we can find roles by ids", func(t *testing.T) {
		roleEntity, err := fixture.Create("Admin")
		if err != nil {
			t.Logf("unexpected error while creating role: %v", err)
		}
		got, err := fixture.FindByIds([]int64{roleEntity.Id})
		if err != nil {
			t.Logf("unexpected error while finding roles by ids: %v", err)
		}
		assert.Nil(err)
		assert.NotNil(got)
		assert.Len(got, 1)
		assert.Equal(roleEntity.Id, got[0].Id)
		assert.Equal(roleEntity.CreatedAt, got[0].CreatedAt)
		assert.Equal(roleEntity.UpdatedAt, got[0].UpdatedAt)
		assert.Equal(roleEntity.Name, got[0].Name)

		clearDb()
	})

	t.Run("we can remove a role by id", func(t *testing.T) {
		roleEntity, err := fixture.Create("Admin")
		if err != nil {
			t.Logf("unexpected error while creating role: %v", err)
		}
		err = fixture.Remove(roleEntity.Id)
		if err != nil {
			t.Logf("unexpected error while removing role by id: %v", err)
		}
		assert.Nil(err)

		clearDb()
	})

	t.Run("we can remove roles by ids", func(t *testing.T) {
		roleEntity, err := fixture.Create("Admin")
		if err != nil {
			t.Logf("unexpected error while creating role: %v", err)
		}
		err = fixture.RemoveByIds([]int64{roleEntity.Id})
		if err != nil {
			t.Logf("unexpected error while removing roles by ids: %v", err)
		}
		assert.Nil(err)

		clearDb()
	})
}
