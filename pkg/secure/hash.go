package secure

import (
	"crypto/sha1"
	"fmt"
)

const salt = "ersdfuli;hnml"

func HashPassword(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))
	hash.Sum([]byte(salt))
	return fmt.Sprintf("%x", hash)

}
