package main

import (
	"encoding/json"
	"fmt"
)

type User struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Age      *int    `json:"age,omitempty"`
	IsActive *bool   `json:"is_active,omitempty"`
	MidName  *string `json:"mid_name"`
}

func Ptr[T any](v T) *T {
	return &v
}
func main() {
	user := User{
		ID:       "test",
		Name:     "test_name",
		Age:      Ptr(30),
		IsActive: Ptr(true),
	}
	buf, _ := json.Marshal(user)
	fmt.Println(string(buf))

	user = User{
		ID:   "test",
		Name: "test_name",
	}
	buf, _ = json.Marshal(user)
	fmt.Println(string(buf))

	user = User{
		ID:       "test",
		Name:     "test_name",
		Age:      nil,
		IsActive: nil,
	}
	buf, _ = json.Marshal(user)
	fmt.Println(string(buf))
}
