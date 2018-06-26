package lz2

import "fmt"

type ErrGroup []error

func (eg ErrGroup) Error() string {
	res := ""
	for _, v := range eg {
		res += v.Error() + "\n"
	}
	return res
}

func (eg ErrGroup) NErrs() int {
	return len(eg)
}

type LineErr struct {
	s string
	l int
}

func (l LineErr) Error() string {
	return fmt.Sprintf("Error on line %d: %s", l.l, l.s)
}

func (l LineErr) Line() int {
	return l.l
}
