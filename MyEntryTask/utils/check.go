package utils

import (
	"net/http"
	"regexp"
	"strconv"
	"time"
)

const (
	CSign      = "sign"
	CTimestamp = "timestamp"
)

func CheckUsername(str string) bool {
	if ok, _ := regexp.MatchString("^[a-zA-Z0-9]{4,16}$", str); !ok {
		return false
	}
	return true
}

func CheckPassword(str string) bool {
	if ok, _ := regexp.MatchString("^[a-zA-Z0-9]{4,}$", str); !ok {
		return false
	}
	return true
}

// CheckTimeStamp 校验时间戳 2min内请求有效
func CheckTimeStamp(str string) (bool, error) {
	now := time.Now()
	//t1 := now.Unix()
	//t2 := time.Unix(t1, 0)

	if str == "" {
		return false, nil
	}
	timestamp, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		Logs.Warn("strconv parse int err :[%v]\n", err)
		return false, err
	}
	reqTime := time.Unix(timestamp, 0)
	if now.Sub(reqTime) > 2*time.Minute {
		return false, nil
	}

	return true, nil
}

// CheckSignature 校验签名
func CheckSignature(sign string, args []string) (bool, error) {
	serverSign, err := MakeSignature(args)
	if err != nil {
		Logs.Warn("make signature err :[%v]\n", err)
		return false, err
	}

	if serverSign != sign {
		Logs.Warn("sign don't match front:[%s] backend:[%s]\n", sign, serverSign)
		return false, nil
	}
	return true, nil
}

func CheckSignAndTimestamp(r *http.Request, reqArgs ...interface{}) (bool, error) {
	sign := r.URL.Query().Get(CSign)
	timestamp := r.URL.Query().Get(CTimestamp)
	args := make([]string, len(reqArgs))
	for i, reqArg := range reqArgs {
		args[i] = reqArg.(string)
	}

	ok, err := CheckTimeStamp(timestamp)
	if err != nil {
		Logs.Warn("check timestamp err :[%v]\n", err)
		return false, err
	}
	if !ok {
		Logs.Warn("timestamp incorrect\n")
		return false, nil
	}

	ok, err = CheckSignature(sign, args)
	if err != nil {
		Logs.Warn("check sign err :[%v]\n", err)
		return false, err
	}

	return ok, nil
}

func checkImg(b *[]byte, ext string) bool {
	switch ext {
	case "jpg", "jpeg":
		if (*b)[0] != 0xff || (*b)[1] != 0xd8 {
			return false
		}
		return true
	case "png":
		if (*b)[0] != 0x89 || (*b)[1] != 0x50 ||
			(*b)[2] != 0x4E || (*b)[3] != 0x47 {
			return false
		}
		return true
	case "gif":
		if (*b)[0] != 0x47 || (*b)[1] != 0x49 || (*b)[2] != 0x46 {
			return false
		}
		return true
	default:
		return false
	}
}

func CheckImageCode(b *[]byte) (bool, string) {
	ok := checkImg(b, "jpg")
	if ok {
		return true, ".jpg"
	}
	ok = checkImg(b, "jpeg")
	if ok {
		return true, ".jpeg"
	}
	ok = checkImg(b, "gif")
	if ok {
		return true, ".gif"
	}
	ok = checkImg(b, "png")
	return ok, ".png"
}
