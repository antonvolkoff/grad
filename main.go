package main

import (
	"fmt"

	"github.com/boltdb/bolt"
)

func main() {
	db, _ := bolt.Open("./db", 0606, nil)
	defer db.Close()

	gdb := &DB{DB: db}
	gdb.Bucket("main").Create()

	fmt.Println("grad v.0.1.0")

	t1 := &Triplet{Subject: "Anton", Predicate: "likes", Object: "golang"}
	gdb.Bucket("main").Save(t1)
	t2 := &Triplet{Subject: "Anton", Predicate: "likes", Object: "ruby"}
	gdb.Bucket("main").Save(t2)

	// Find by subject
	fmt.Println("Find by subject")
	results := gdb.Bucket("main").Find(&Triplet{Subject: "Anton"})
	for _, result := range results {
		fmt.Println(result)
	}

	// Find by predicate
	fmt.Println("Find by predicate")
	results = gdb.Bucket("main").Find(&Triplet{Predicate: "likes"})
	for _, result := range results {
		fmt.Println(result)
	}

	// Find by object
	fmt.Println("Find by object")
	results = gdb.Bucket("main").Find(&Triplet{Object: "golang"})
	for _, result := range results {
		fmt.Println(result)
	}
}
