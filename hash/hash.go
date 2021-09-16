//

package hash

type Hash interface {
	Check(password string, hash string) bool
}
