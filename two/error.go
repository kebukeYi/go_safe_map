package two

import "fmt"

type outOfRangeError struct {
	upperIndex int
	lowerIndex int
}

func (out outOfRangeError) Error() string {
	return fmt.Sprintf("kv-db out of range with [%v:%v]", out.upperIndex, out.lowerIndex)
}

type notFoundError struct {
	key interface{}
}

func (not notFoundError) Error() string {
	return fmt.Sprintf("kv-db cannot find key:[%v]", not.key)
}

type keyExpiredError struct {
	key interface{}
}

func (e keyExpiredError) Error() string {
	return fmt.Sprintf("kv-db the key [%v] is expried ", e.key)
}
