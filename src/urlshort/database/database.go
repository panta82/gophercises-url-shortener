package database

import (
	"github.com/boltdb/bolt"
	"time"
	"path/filepath"
	"os"
	. "urlshort/types"
)

const defaultDatabaseFilename = "data.db"
const redirectsBucket = "redirects"

type Database struct {
	Path string
}

func GetDefaultDatabasePath() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic("Failed to determine directory of the executable")
	}

	return filepath.Join(dir, defaultDatabaseFilename)
}

func NewDatabase(path string) (Database, error) {
	db := Database{ Path: path}
	err := initializeDatabase(db)
	return db, err
}

func initializeDatabase(db Database) error {
	return db.exec(func(boltDB bolt.DB) error {
		return boltDB.Update(func(tx *bolt.Tx) error {
			_, err := tx.CreateBucketIfNotExists([]byte(redirectsBucket))
			return err
		})
	})
}

func (db Database) exec (fn func (boltDB bolt.DB) error) error {
	boltDB, err := bolt.Open(db.Path, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return err
	}

	defer boltDB.Close()

	return fn(*boltDB)
}

func (db Database) GetUrlForPath(path string) (string, error) {
	result := ""
	err := db.exec(func(boltDB bolt.DB) error {
		return boltDB.View(func(tx *bolt.Tx) error {
			bucket := tx.Bucket([]byte(redirectsBucket))
			bResult := bucket.Get([]byte(path))
			if bResult != nil {
				result = string(bResult[:])
			}

			return nil
		})
	})

	if err != nil {
		return "", err
	}

	return result, nil
}

func (db Database) ListAllRedirects() ([]Redirect, error) {
	var results []Redirect
	err := db.exec(func(boltDB bolt.DB) error {
		return boltDB.View(func(tx *bolt.Tx) error {
			bucket := tx.Bucket([]byte(redirectsBucket))

			return bucket.ForEach(func(k, v []byte) error {
				results = append(results, Redirect{
					Path: string(k[:]),
					Url: string(v[:]),
				})
				return nil
			})

			return nil
		})
	})

	return results, err
}

func (db Database) SetUrlForPath(path string, url string) (error) {
	return db.exec(func(boltDB bolt.DB) error {
		return boltDB.Update(func(tx *bolt.Tx) error {
			bucket := tx.Bucket([]byte(redirectsBucket))
			return bucket.Put([]byte(path), []byte(url))
		})
	})
}

func (db Database) RemoveUrlForPath(path string) (error) {
	return db.exec(func(boltDB bolt.DB) error {
		return boltDB.Update(func(tx *bolt.Tx) error {
			bucket := tx.Bucket([]byte(redirectsBucket))
			return bucket.Delete([]byte(path))
		})
	})
}

