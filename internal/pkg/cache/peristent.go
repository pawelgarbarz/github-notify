package cache

import (
	"errors"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

type database interface {
	Exec(query string, args ...any) error
	QueryRow(destination interface{}, query string, args ...any) error
	IsEmptyResult(err error) bool
}

type Cache interface {
	Save(key string, value string, ttl time.Duration) error
	Get(key string) (string, error)
	Exists(key string) (bool, error)
	Delete(key string) error
	ClearOutdated() error
	ClearAll() error
}

type cache struct {
	database database
}

type item struct {
	Key       string
	Value     string
	CreatedAt time.Time
	DeleteAt  time.Time
}

var errKeyAlreadyExists = "UNIQUE constraint failed: cache.key"
var ErrIDNotFound = errors.New("cache item not found")

func NewCache(db database) Cache {

	return &cache{database: db}
}

func (c *cache) Save(key string, value string, ttl time.Duration) error {
	now := time.Now().UTC()
	deleteAt := c.calculateDeleteAt(now, ttl)
	err := c.database.Exec("INSERT INTO cache VALUES(?,?,?,?);", key, value, now, deleteAt)
	if err != nil && err.Error() == errKeyAlreadyExists {
		err = c.update(key, value, ttl)
	}

	return err
}

func (c *cache) Get(key string) (string, error) {
	item := item{}
	err := c.database.QueryRow(&item, "SELECT key, value, createdAt, deleteAt FROM cache WHERE key=?", key)
	if c.database.IsEmptyResult(err) {
		return item.Value, ErrIDNotFound
	}

	if !item.DeleteAt.IsZero() && time.Now().After(item.DeleteAt) {
		err = c.Delete(key)
		if err != nil {
			return "", err
		}

		return "", ErrIDNotFound
	}

	return item.Value, err
}

func (c *cache) Exists(key string) (bool, error) {
	_, err := c.Get(key)

	return err != ErrIDNotFound, err
}

func (c *cache) Delete(key string) error {
	return c.database.Exec("DELETE FROM cache WHERE key=?;", key)
}

func (c *cache) ClearOutdated() error {
	return c.database.Exec("DELETE FROM cache WHERE deleteAt < date('now');")
}

func (c *cache) ClearAll() error {
	return c.database.Exec("DELETE FROM cache;")
}

func (c *cache) update(key string, value string, ttl time.Duration) error {
	now := time.Now().UTC()
	deleteAt := c.calculateDeleteAt(now, ttl)
	err := c.database.Exec("UPDATE cache SET value=?, createdAt=?, deleteAt=? WHERE key=?;", value, now, deleteAt, key)

	return err
}

func (c *cache) calculateDeleteAt(now time.Time, ttl time.Duration) time.Time {
	if ttl > 0 {
		return now.Add(ttl)
	}

	return time.Time{}
}
