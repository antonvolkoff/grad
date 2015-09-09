package main

import (
	"testing"

	"github.com/boltdb/bolt"
	"github.com/stretchr/testify/assert"
)

func initGDB() *DB {
	db, _ := bolt.Open("./db_test", 0606, nil)

	gdb := &DB{DB: db}
	gdb.Bucket("main").Create()
	return gdb
}

func TestFindByObject(t *testing.T) {
	gdb := initGDB()
	defer gdb.DB.Close()
	b := gdb.Bucket("main")

	b.Save(&Triplet{"Anton", "follows", "Mike"})

	result := b.Find(&Triplet{Subject: "Anton", Predicate: "follows"})

	assert.Len(t, result, 1)
	assert.Equal(t, &Triplet{"Anton", "follows", "Mike"}, result[0])
}
