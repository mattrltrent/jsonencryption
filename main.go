package main

//! CONTRIVED EXAMPLE USAGE

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/mattrltrent/jsonencryption/jsonencryption"
)

func init() {
	// you must set a key (probably should be loaded via env var)
	jsonencryption.SetKey("16-byte-key-1234")
}

type User struct {
	ID   jsonencryption.EncryptedID
	Name string
	Age  int
}

func main() {
	// create a new user with an encrypted ID
	user := User{ID: jsonencryption.NewEncryptedID(123), Name: "John", Age: 20}

	// simulate turning it to JSON
	jsonData, err := json.Marshal(user)
	if err != nil {
		log.Fatalf("Error occurred during marshaling. Error: %s", err.Error())
	}

	// look... it's encrypted!
	fmt.Println(string(jsonData))
}
