package module_verification

import (
	"crypto/rand"
	"errors"
	"io"
	"time"

	module_helpers "github.com/VDiPaola/go-backend-module/module_helpers"
	module_models "github.com/VDiPaola/go-backend-module/module_models"
)

func GenerateCode(length int, expirationLengthMins int64) module_models.Code {
	var nums = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}
	b := make([]byte, length)
	n, err := io.ReadAtLeast(rand.Reader, b, length)
	if n != length {
		panic(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = nums[int(b[i])%len(nums)]
	}

	return module_models.Code{
		Value:     string(b),
		ExpiresAt: module_helpers.UnixNanoToJS(time.Now().Add(time.Minute * time.Duration(expirationLengthMins)).UnixNano()),
	}
}

func VerifyCode(code module_models.Code, sentCode string) error {

	//check that codes match
	if code.Value != sentCode {
		return errors.New("code doesn't match")
	}

	//check code hasnt expired
	if module_helpers.NowJS() > code.ExpiresAt {
		return errors.New("code has expired")
	}

	return nil
}
