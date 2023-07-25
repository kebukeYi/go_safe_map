package one

import "fmt"

type outOfRangeError struct {
	upperIndex int
	lowerIndex int
}

func (e outOfRangeError) Error() string {
	return fmt.Sprintf("kv-db out of range: upper=%d, lower=%d", e.upperIndex, e.lowerIndex)
}

type notFoundError struct {
	key interface{}
}

type keyExpired struct {
	key interface{}
}

func (e notFoundError) Error() string {
	return fmt.Sprintf("kv-db key not found: %v", e.key)
}
