package cache

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestCache_Save(t *testing.T) {
	db := newDbMock()
	instance := NewCache(db)

	db.
		On("Exec",
			"INSERT INTO cache VALUES(?,?,?,?);", "key", "val",
			mock.Anything, mock.Anything).
		Return(nil)

	err := instance.Save("key", "val", time.Second*10)

	assert.Nil(t, err)

	db.
		On("Exec",
			"INSERT INTO cache VALUES(?,?,?,?);", "key2", "val2",
			mock.Anything, mock.Anything).
		Return(nil)

	err2 := instance.Save("key", "val", 0)

	assert.Nil(t, err2)
}

func TestCache_SaveWithUpdate(t *testing.T) {
	db := newDbMock()
	instance := NewCache(db)

	db.
		On("Exec",
			"INSERT INTO cache VALUES(?,?,?,?);",
			"key", "val", mock.Anything, mock.Anything).
		Return(errKeyAlreadyExists)
	db.
		On("Exec",
			"UPDATE cache SET value=?, createdAt=?, deleteAt=? WHERE key=?;",
			"val", mock.Anything, mock.Anything, "key").
		Return(nil)

	err := instance.Save("key", "val", time.Second*10)

	assert.Nil(t, err)
}

func TestCache_SaveWithError(t *testing.T) {
	db := newDbMock()
	instance := NewCache(db)

	customError := errors.New("custom error")

	db.
		On("Exec",
			"INSERT INTO cache VALUES(?,?,?,?);",
			"key", "val", mock.Anything, mock.Anything).
		Return(customError.Error())

	err := instance.Save("key", "val", time.Second*10)

	assert.Equal(t, customError, err)
}

func TestCache_Get(t *testing.T) {
	db := newDbMock()
	instance := NewCache(db)

	db.On("IsEmptyResult", nil).Return(false)

	db.
		On("QueryRow",
			&item{},
			"SELECT key, value, createdAt, deleteAt FROM cache WHERE key=?",
			"key").
		Return(nil).
		Run(func(args mock.Arguments) {
			arg := args.Get(0).(*item)
			arg.Key = "key"
			arg.Value = "val"
			arg.CreatedAt = time.Now()
			arg.DeleteAt = time.Now().Add(time.Minute * 10)
		})

	val, err := instance.Get("key")

	assert.Nil(t, err)
	assert.Equal(t, "val", val)
}

func TestCache_Get_NotFound(t *testing.T) {
	db := newDbMock()
	instance := NewCache(db)

	db.On("IsEmptyResult", ErrIDNotFound).Return(true)

	db.
		On("QueryRow",
			&item{},
			"SELECT key, value, createdAt, deleteAt FROM cache WHERE key=?",
			"key").
		Return(ErrIDNotFound)

	_, err := instance.Get("key")

	assert.Equal(t, ErrIDNotFound, err)
}

func TestCache_Get_FoundButShouldBeDeleted(t *testing.T) {
	db := newDbMock()
	instance := NewCache(db)

	db.On("IsEmptyResult", nil).Return(false)

	db.On("Exec", "DELETE FROM cache WHERE key=?;", "key").Return(nil)

	db.
		On("QueryRow",
			&item{},
			"SELECT key, value, createdAt, deleteAt FROM cache WHERE key=?",
			"key").
		Return(nil).
		Run(func(args mock.Arguments) {
			arg := args.Get(0).(*item)
			arg.Key = "key"
			arg.Value = "val"
			arg.CreatedAt = time.Now()
			arg.DeleteAt = time.Now().AddDate(0, 0, -1)
		})

	val, err := instance.Get("key")

	assert.Equal(t, "", val)
	assert.Equal(t, ErrIDNotFound, err)
}

func TestCache_Get_FoundButShouldBeDeleted_DeleteErr(t *testing.T) {
	db := newDbMock()
	instance := NewCache(db)

	customeError := errors.New("custom Error")

	db.On("IsEmptyResult", nil).Return(false)

	db.On("Exec", "DELETE FROM cache WHERE key=?;", "key").Return(customeError.Error())

	db.
		On("QueryRow",
			&item{},
			"SELECT key, value, createdAt, deleteAt FROM cache WHERE key=?",
			"key").
		Return(nil).
		Run(func(args mock.Arguments) {
			arg := args.Get(0).(*item)
			arg.Key = "key"
			arg.Value = "val"
			arg.CreatedAt = time.Now()
			arg.DeleteAt = time.Now().AddDate(0, 0, -1)
		})

	val, err := instance.Get("key")

	assert.Equal(t, "", val)
	assert.Equal(t, customeError, err)
}

func TestCache_Exists(t *testing.T) {
	db := newDbMock()
	instance := NewCache(db)

	db.On("Exists", "key").Return(true, errKeyAlreadyExists)

	db.
		On("QueryRow",
			&item{},
			"SELECT key, value, createdAt, deleteAt FROM cache WHERE key=?",
			"key").
		Return(nil).
		Run(func(args mock.Arguments) {
			arg := args.Get(0).(*item)
			arg.Key = "key"
			arg.Value = "val"
			arg.CreatedAt = time.Now()
			arg.DeleteAt = time.Now().AddDate(0, 0, 1)
		})

	db.On("IsEmptyResult", nil).Return(false)

	exists, err := instance.Exists("key")

	assert.Nil(t, err)
	assert.Equal(t, true, exists)
}

func TestCache_ClearOutdated(t *testing.T) {
	db := newDbMock()
	instance := NewCache(db)

	db.On("Exec", "DELETE FROM cache WHERE deleteAt > date('now');").Return(nil)

	err := instance.ClearOutdated()

	assert.Nil(t, err)
}

func TestCache_ClearAll(t *testing.T) {
	db := newDbMock()
	instance := NewCache(db)

	db.On("Exec", "DELETE FROM cache;").Return(nil)

	err := instance.ClearAll()

	assert.Nil(t, err)
}
