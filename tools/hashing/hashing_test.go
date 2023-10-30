package hashing

import (
	"fmt"
	"testing"
)

func TestHashPassword(t *testing.T) {
	fmt.Println(HashPassword("12345678!!"))
}
