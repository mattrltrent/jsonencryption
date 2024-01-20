package jsonencryption

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"
)

var secret *Secret

func NewEncryptedID(val uint) EncryptedID {
	return EncryptedID{Val: val}
}

func SetKey(key string) {
	secret = NewSecret(key)
}

type EncryptedID struct {
	Val uint
}

func (i EncryptedID) ToString() string {
	return fmt.Sprintf("%d", i.Val)
}

func (i EncryptedID) ToInt() int {
	return int(i.Val)
}

func (i EncryptedID) ToMasked() string {
	masked, _ := secret.Mask(i.Val)
	return masked
}

func (i EncryptedID) MarshalJSON() ([]byte, error) {
	// Check if the secret object itself is nil
	if secret == nil {
		return nil, ErrNoSecret
	}

	// Check if the secret.val is nil or not the expected length
	if secret.val == nil || len(secret.val) != 16 {
		return nil, ErrInvalidKey
	}

	masked, err := secret.Mask(i.Val)
	if err != nil {
		return nil, err
	}

	hashed := Hash(i.Val)
	data := map[string]interface{}{
		"mask": masked, // volatile, reversible
		"hash": hashed, // consistent, not reversible
	}

	return json.Marshal(data)
}

func (i *EncryptedID) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	switch v := value.(type) {
	case int64:
		i.Val = uint(v)
	case []byte:
		// convert []byte to uint
		var intValue int64
		var err error
		if intValue, err = strconv.ParseInt(string(v), 10, 64); err != nil {
			return err
		}
		i.Val = uint(intValue)
	case string:
		// convert string to uint
		var intValue int64
		var err error
		if intValue, err = strconv.ParseInt(v, 10, 64); err != nil {
			return err
		}
		i.Val = uint(intValue)
	default:
		return fmt.Errorf("unsupported scan value type: %T", value)
	}

	return nil
}

func (i EncryptedID) Value() (driver.Value, error) {
	// Convert the uint to int64, which is a common type for IDs in databases
	return int64(i.Val), nil
}
