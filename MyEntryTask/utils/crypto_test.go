package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"testing"
	"time"
)

var data = "123456"

func TestEncryData(t *testing.T) {
	starttime := time.Now()
	hash, err := tEncryData(data)
	if err != nil {
		t.Errorf("encrydata err")
	}
	fmt.Println("encrpt time :", time.Since(starttime))
	fmt.Println(hash)
}

func tEncryData(str string) (string, error) {
	hash := sha256.New()
	_, err := hash.Write([]byte(str))
	// 加盐
	hash.Write([]byte("seatalk"))
	if err != nil {
		Logs.Warn("EncryData data hash write err :[%v]\n}", err)
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}
