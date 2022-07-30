package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// EncryData 采用sha256哈希 用作签名和哈希password
func EncryData(str string) (string, error) {
	hash := sha256.New()
	_, err := hash.Write([]byte(str))
	// 加盐
	//hash.Write([]byte(Conf.Salt))
	hash.Write([]byte("seatalk"))
	if err != nil {
		Logs.Warn("EncryData data hash write err :[%v]\n}", err)
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// MakeSignature 根据参数生成签名
func MakeSignature(args []string) (string, error) {
	var str string
	for _, arg := range args {
		// 签名的规则 $arg1_$arg2_..._$salt
		str += arg
		str += "_"
	}

	sign, err := EncryData(str)
	if err != nil {
		Logs.Warn("EncryData data err :[%v]\n}", err)
		return "", err
	}
	return sign, nil
}

//// GenUUid 生成uuid
//func GenUUid() (string, error) {
//	// 借助类linux下的uuidgen生成uuid
//	out, err := exec.Command("uuidgen").Output()
//	if err != nil {
//		Logs.Warn("uuidgen err :[%v]\n", err)
//		return "", err
//	}
//
//	return string(out[:len(out)-2]), nil
//}

func GenUUid() (uuid string, err error) {
	u := new([16]byte)
	_, err = rand.Read(u[:])
	if err != nil {
		Logs.Warn("Cannot generate UUID", err)
		return "", err
	}
	u[8] = (u[8] | 0x40) & 0x7F
	u[6] = (u[6] & 0xF) | (0x4 << 4)
	uuid = fmt.Sprintf("%x-%x-%x-%x-%x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:])
	return
}
