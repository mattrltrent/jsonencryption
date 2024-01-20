# Golang package: `jsonencryption`

![unit tests](https://github.com/mattrltrent/jsonencryption/actions/workflows/unit_tests.yml/badge.svg)

[import "github.com/mattrltrent/jsonencryption/jsonencryption"](https://github.com/mattrltrent/jsonencryption)

### What does it do?

It takes a struct like this:

```go
user := User{ID: 123, Name: "John", Age: 20}
```

That normally serializes as:

```json
{
  "ID": 123 // <-- bad! this reveals your sensitive ID
  "Name": "John",
  "Age": 20
}
```

And allows you to instead populate the struct like this:

```go
user := User{ID: jsonencryption.NewEncryptedID(123), Name: "John", Age: 20}
```

Such that it now serializes like this with encrypted values:

```json
{
  "ID": {
    "hash": "pmWkWSBCL51Bfkhn79xPuKBKHz__H6B-mY6G9_eieuM",
    "mask": "GpvBdoyRDZR3vU1BC0ptrkvEhQ=="
  },
  "Name": "John",
  "Age": 20
}
```

### Serialized fields: `hash` & `mask`

**`"hash"`**: An irreversible string. Two of the same `uint` values will always have the same hash. Useful to ensure uniqueness on the client-side between different recieved JSON documents. 

**`"mask"`**: A string *only* reversible by your code since you set the 16-byte secret string that was used to mask it. Two of the same `uint` values will have different masks. Useful to have an opaque mask of the true `uint` that only you can read.

### Working with: `EncryptedID` & `uint`

```go
// create a new EncryptedID (providing some uint argument)
encryptedID := jsonencryption.NewEncryptedID(123) // type: jsonencryption.EncryptedID

// get the uint back
value := encryptedID.Val // type: uint
```

### Re: Databases

When you're using a struct as a model class with an `EncryptedID`, it should automatically send to the database just the inner `uint`.

### Model class usecase example

You have a model class that you put `EncryptedID` into as the primary key:

```go
type BlogArticle struct {
	ID         EncryptedID `json:"id"`
	PostedDate time.Time   `json:"posted_date"`
	Title      string      `json:"title"`
	Body       string      `json:"body"`
	Author     string      `json:"author"`
	Slug       string      `json:"slug"`
}
```

You can then both use this to send to the client (as JSON) and interact with the database. It will conceal the inner `uint` value of the `EncryptedID` from the client (with a "mask" and "hash"), but use the real `uint` with the database.

### Simple example usage

**TLDR:** Set your PRIVATE AND SENSITIVE 16-byte encryption key (`jsonencryption.SetKey("16-byte-key-1234")`) and then use `jsonencryption.EncryptedID` however you please.

```go
package main

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
```

### Other methods

Check out the [encryption_test.go](https://github.com/mattrltrent/jsonencryption/blob/main/jsonencryption/encryption_test.go) file to see what other methods can be utilized.

### Unit tests

This will run all this package's unit tests:

```sh
./scripts/test ./...
```

### Disclaimer

Made for fun as a non-security professional. Use at your own risk.