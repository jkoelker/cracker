//

package hash

import (
	"golang.org/x/crypto/bcrypt"
)

type BCrypt struct {
}

func NewBCrypt() *BCrypt {
	return &BCrypt{
	}
}

func (b *BCrypt) Check(password string, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

var _ Hash = (*BCrypt)(nil)
