package db

import (
	"encoding/binary"
	"errors"
	"github.com/boltdb/bolt"
	"log"
	"time"
)

const TASK_BUCKET = "tasks"

type DB struct {
	db *bolt.DB
}

func NewDB(path string) (DB, error) {
	db, err := bolt.Open("my.db", 0600, &bolt.Options{
		Timeout: 1 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(TASK_BUCKET))
		return err
	})
	if err != nil {
		return DB{}, err
	}
	return DB{db: db}, nil
}

type Task struct {
	Key   []byte
	Value string
}

func (d *DB) GetList() ([]Task, error) {
	var tasks []Task
	err := d.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(TASK_BUCKET))
		c := b.Cursor()
		for k, v := c.First(); c != nil; k, v = c.Next() {
			task := Task {k, string(v) }
			tasks = append(tasks, task)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

var Found = errors.New("found task")

func (d *DB) Delete(idx int) error {
	var t Task
	var i int
	err := d.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(TASK_BUCKET))
		c := b.Cursor()
		for k, v := c.First(); c != nil; k, v = c.Next() {
			if i == idx-1 {
				t.Key = k
				t.Value = string(v)
				return Found
			}
			i++
		}
		return nil
	})
	if err != nil {
		return err
	}
	err = d.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(TASK_BUCKET))
		return b.Delete(t.Key)
	})
	return err
}

func (d *DB) AddTasks(tasks ...string) error {
	return d.db.Update(func(tx *bolt.Tx) error {
		bkt := tx.Bucket([]byte(TASK_BUCKET))
		for _, task := range tasks {
			id, err := bkt.NextSequence()
			if err != nil {
				return err
			}
			if err = bkt.Put(itob(id), []byte(task)); err != nil {
				return err
			}
		}
		return nil
	})
}

func itob(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, v)
	return b
}

func btoi(b []byte) int {
	return int(binary.BigEndian.Uint64(b))
}
