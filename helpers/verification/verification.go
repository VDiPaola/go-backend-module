package verification

import (
	"crypto/rand"
	"errors"
	"io"
	"time"

	"github.com/VDiPaola/go-backend-module/helpers"
	"github.com/VDiPaola/go-backend-module/models"
)

func GenerateCode(length int, expirationLengthMins int64) models.Code {
	var nums = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}
	b := make([]byte, length)
	n, err := io.ReadAtLeast(rand.Reader, b, length)
	if n != length {
		panic(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = nums[int(b[i])%len(nums)]
	}

	return models.Code{
		Value:     string(b),
		ExpiresAt: helpers.UnixNanoToJS(time.Now().Add(time.Minute * time.Duration(expirationLengthMins)).UnixNano()),
	}
}

func VerifyCode(code models.Code, sentCode string) error {

	//check that codes match
	if code.Value != sentCode {
		return errors.New("code doesn't match")
	}

	//check code hasnt expired
	if helpers.NowJS() > code.ExpiresAt {
		return errors.New("code has expired")
	}

	return nil
}
