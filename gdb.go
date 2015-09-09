package main

import (
	"bytes"
	"encoding/json"

	"github.com/boltdb/bolt"
)

// Triplet struct
type Triplet struct {
	Subject   string `json:"subject"`
	Predicate string `json:"predicate"`
	Object    string `json:"object"`
}

// DB struct
type DB struct {
	DB *bolt.DB
}

// Bucket struct
type Bucket struct {
	db   *bolt.DB
	Name string
}

// Bucket returns reference to a bucket
func (db *DB) Bucket(name string) *Bucket {
	return &Bucket{
		db:   db.DB,
		Name: name,
	}
}

// Create inits a bucket unless it exists already
func (b *Bucket) Create() {
	db := b.db
	db.Update(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists([]byte(b.Name))
		return nil
	})
}

// Save stores triplet inside a bucket
func (b *Bucket) Save(t *Triplet) {
	db := b.db

	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("main"))

		value, _ := json.Marshal(t)
		keySPO := "spo::" + t.Subject + "::" + t.Predicate + "::" + t.Object
		keyPOS := "pos::" + t.Predicate + "::" + t.Object + "::" + t.Subject
		keyOSP := "osp::" + t.Object + "::" + t.Subject + "::" + t.Predicate
		keyOPS := "ops::" + t.Object + "::" + t.Predicate + "::" + t.Subject

		b.Put([]byte(keySPO), value)
		b.Put([]byte(keyPOS), value)
		b.Put([]byte(keyOSP), value)
		b.Put([]byte(keyOPS), value)

		return nil
	})
}

// Find finds triplet with query object
func (b *Bucket) Find(query *Triplet) []*Triplet {
	db := b.db
	results := []*Triplet{}

	var key string
	switch {
	case query.Subject != "" && query.Predicate != "" && query.Object != "":
		key = "spo::" + query.Subject + "::" + query.Predicate + "::" + query.Object
	case query.Subject != "" && query.Predicate != "":
		key = "spo::" + query.Subject + "::" + query.Predicate + "::"
	case query.Subject != "":
		key = "spo::" + query.Subject
	case query.Predicate != "" && query.Object != "":
		key = "pos::" + query.Predicate + "::" + query.Object + "::"
	case query.Predicate != "":
		key = "pos::" + query.Predicate + "::"
	case query.Object != "" && query.Subject != "":
		key = "osp::" + query.Object + "::" + query.Subject + "::"
	case query.Object != "" && query.Predicate != "":
		key = "ops::" + query.Object + "::" + query.Predicate + "::"
	case query.Object != "":
		key = "osp::" + query.Object + "::"
	default:
		key = ""
	}

	if key == "" {
		return results
	}

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("main"))
		c := b.Cursor()

		prefix := []byte(key)

		for k, v := c.Seek(prefix); bytes.HasPrefix(k, prefix); k, v = c.Next() {
			var t Triplet
			json.Unmarshal(v, &t)
			results = append(results, &t)
		}

		return nil
	})

	return results
}
