package employee

import (
	"github.com/78bits/go-sqlmock-sqlx"
	assertpackage "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"idm/inner/database"
	"idm/inner/employee"
	"regexp"
	"testing"
	"time"
)

func TestEmployeeRepository(t *testing.T) {
	assert := assertpackage.New(t)
	var db = database.ConnectDb()

	var clearDb = func() {
		db.MustExec("DELETE FROM employees")
	}

	defer func() {
		if r := recover(); r != nil {
			clearDb()
		}
	}()

	var employeeRepository = employee.NewRepository(db)
	var fixture = NewFixture(employeeRepository)

	t.Run("we can create an employee", func(t *testing.T) {
		got, err := fixture.CreateEmployee("John Doe")
		if err != nil {
			t.Logf("unexpected error while creating employee: %v", err)
		}
		assert.Nil(err)
		assert.NotNil(got)
		assert.NotEmpty(got.Id)
		assert.NotEmpty(got.CreatedAt)
		assert.NotEmpty(got.UpdatedAt)
		assert.Equal("John Doe", got.Name)

		clearDb()
	})

	t.Run("we can find employee by id", func(t *testing.T) {
		emp, _ := fixture.CreateEmployee("John Doe")
		got, err := fixture.FindById(emp.Id)
		if err != nil {
			t.Logf("unexpected error while finding employee by id: %v", err)
		}
		assert.Nil(err)
		assert.NotNil(got)
		assert.Equal(emp.Id, got.Id)
		assert.Equal(emp.CreatedAt, got.CreatedAt)
		assert.Equal(emp.UpdatedAt, got.UpdatedAt)
		assert.Equal(emp.Name, got.Name)

		clearDb()
	})

	t.Run("we can find all employees", func(t *testing.T) {
		emp1, _ := fixture.CreateEmployee("John Doe")
		emp2, _ := fixture.CreateEmployee("Jane Doe")
		got, err := fixture.FindAll()
		if err != nil {
			t.Logf("unexpected error while finding employees: %v", err)
		}
		assert.Nil(err)
		assert.NotNil(got)
		assert.Len(got, 2)
		assert.Equal(emp1.Id, got[0].Id)
		assert.Equal(emp1.CreatedAt, got[0].CreatedAt)
		assert.Equal(emp1.UpdatedAt, got[0].UpdatedAt)
		assert.Equal(emp1.Name, got[0].Name)
		assert.Equal(emp2.Id, got[1].Id)
		assert.Equal(emp2.CreatedAt, got[1].CreatedAt)
		assert.Equal(emp2.UpdatedAt, got[1].UpdatedAt)
		assert.Equal(emp2.Name, got[1].Name)

		clearDb()
	})

	t.Run("we can find employees by ids", func(t *testing.T) {
		emp1, _ := fixture.CreateEmployee("John Doe")
		emp2, _ := fixture.CreateEmployee("Jane Doe")
		_, _ = fixture.CreateEmployee("Jose Doe")
		got, err := fixture.FindByIds([]int64{emp1.Id, emp2.Id})
		if err != nil {
			t.Logf("unexpected error while finding employees by ids: %v", err)
		}

		assert.Nil(err)
		assert.NotNil(got)
		assert.Len(got, 2)
		assert.Equal(emp1.Id, got[0].Id)
		assert.Equal(emp1.CreatedAt, got[0].CreatedAt)
		assert.Equal(emp1.UpdatedAt, got[0].UpdatedAt)
		assert.Equal(emp1.Name, got[0].Name)
		assert.Equal(emp2.Id, got[1].Id)
		assert.Equal(emp2.CreatedAt, got[1].CreatedAt)
		assert.Equal(emp2.UpdatedAt, got[1].UpdatedAt)
		assert.Equal(emp2.Name, got[1].Name)

		clearDb()
	})

	t.Run("we can remove an employee", func(t *testing.T) {
		emp, err := fixture.CreateEmployee("John Doe")
		if err != nil {
			t.Logf("unexpected error while creating employee: %v", err)
		}

		id := emp.Id
		assert.NotNil(emp)
		assert.NotEmpty(emp.Id)
		assert.NotEmpty(emp.CreatedAt)
		assert.NotEmpty(emp.UpdatedAt)
		assert.Equal("John Doe", emp.Name)

		err = fixture.Remove(emp.Id)
		if err != nil {
			t.Logf("unexpected error while removing employee: %v", err)
		}
		assert.Nil(err)

		got, err := fixture.FindById(id)
		if err != nil {
			assert.Error(database.ErrRecordNotFound)
		}
		assert.Nil(got)

		clearDb()
	})

	t.Run("we can remove employees by ids", func(t *testing.T) {
		emp1, _ := fixture.CreateEmployee("John Doe")
		emp2, _ := fixture.CreateEmployee("Jane Doe")

		ids := []int64{emp1.Id, emp2.Id}

		assert.NotNil(emp1)
		assert.NotNil(emp2)
		assert.NotEmpty(emp1.Id)
		assert.NotEmpty(emp1.CreatedAt)
		assert.NotEmpty(emp1.UpdatedAt)
		assert.Equal("John Doe", emp1.Name)
		assert.NotEmpty(emp2.Id)
		assert.NotEmpty(emp2.CreatedAt)
		assert.NotEmpty(emp2.UpdatedAt)
		assert.Equal("Jane Doe", emp2.Name)

		err := fixture.RemoveByIds(ids)
		if err != nil {
			t.Logf("unexpected error while removing employees by ids: %v", err)
		}
		assert.Nil(err)

		got, err := fixture.FindByIds(ids)
		if err != nil {
			t.Logf("unexpected error while finding employees by ids: %v", err)
		}
		assert.Error(database.ErrRecordNotFound)
		assert.Nil(got)

		clearDb()
	})
}

func TestRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	require.NoError(t, err)
	defer db.Close()

	r := employee.NewRepository(db)

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM employees WHERE name = \\$1").
		WithArgs("John Doe").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
	mock.ExpectQuery(regexp.QuoteMeta(
		"INSERT INTO employees (name) VALUES ($1) RETURNING id, created_at, updated_at")).
		WithArgs("John Doe").
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).AddRow(1, time.Now(), time.Now()))
	mock.ExpectCommit()

	err = r.Create(&employee.Employee{Name: "John Doe"})
	require.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
