package database

import (
	"fmt"
	"os"
	"testing"
)

func TestCreateDB(t *testing.T) {
	os.Remove("database.json")

	db, err := NewDB("database.json")
	if err != nil {
		t.Error(err)
	}
	actual, err := db.createChirp("Test 1")
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("Create chirp: %v\n", actual)

	if actual.Id != 0 {
		t.Errorf("actual: %v, expected: ´0´", actual.Id)
	}
	if actual.Body != "Test 1" {
		t.Errorf("actual: %v, expected: ´Test 1´", actual.Body)
	}

	chirps, err := db.GetChirps()
	if err != nil {
		t.Error(err)
	}
	for _, cc := range chirps {
		fmt.Printf("%v\n", cc)
	}
}

// func TestGet(t *testing.T) {
//         interval := time.Millisecond * 10
//         cache := NewCache(interval)
//
//         cache.Add("key1", []byte("val1"))
//         actual, ok := cache.Get("key1")
//         if !ok {
//                 t.Errorf("actual: %v, expected: 'val1'", actual)
//         }
// }
