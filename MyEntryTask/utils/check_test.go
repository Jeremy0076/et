package utils

import (
	"log"
	"testing"
)

func TestCheckUsername(t *testing.T) {
	if ok := CheckUsername("aaaa"); !ok {
		t.Errorf("check username error")
	}
}

func TestCheckPassword(t *testing.T) {
	if ok := CheckPassword("aaaa"); !ok {
		t.Errorf("check username error")
	}
}

func TestCheckSignAndTimestamp(t *testing.T) {
	token := "test01"
	username := "test01"
	strs := []string{token, username}
	testSign, err := MakeSignature(strs)
	if err != nil {
		t.Errorf("make sign error")
		return
	}
	log.Printf("sign=%s\n", testSign)
}
