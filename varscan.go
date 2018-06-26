package lz2

import (
	"os"
	"path"
	"strings"
)

type Replacer func(string) string

type varBlock struct {
	S   string
	isV bool
}

func (vn varBlock) String() string {
	if vn.isV {
		return "Var:" + vn.S
	}
	return "Str:" + vn.S
}

func varScan(s, op, cl string) []varBlock {
	mode := 0
	skip := 0
	cstring := ""
	res := []varBlock{}
	for k, v := range s {
		if skip > k {
			continue
		}
		if mode == 0 && strings.HasPrefix(s[k:], op) {
			skip = k + len(op)
			mode = 1
			if cstring != "" {
				res = append(res, varBlock{cstring, false})
			}
			cstring = ""
			continue
		}

		if mode == 1 && strings.HasPrefix(s[k:], cl) {
			skip = k + len(cl)
			mode = 0
			res = append(res, varBlock{cstring, true})
			cstring = ""
			continue
		}
		cstring += string(v)
	}
	if cstring != "" {
		res = append(res, varBlock{cstring, false})
	}
	return res
}

func varReplace(vv []varBlock, f Replacer) string {
	res := ""
	for _, v := range vv {
		if !v.isV {
			res += v.S
			continue
		}
		res += f(v.S)
	}
	return res
}

//ScanReplace is the main method to call for a replace method
func ScanReplace(s, op, cl string, f Replacer) string {
	t := varScan(s, op, cl)
	return varReplace(t, f)
}

func EnvReplace(s string) string {
	return ScanReplace(s, "{", "}", os.Getenv)
}

func MapReplacer(m map[string]string) Replacer {
	return func(s string) string {
		r, _ := m[s]
		return r
	}
}

func PlusPathEnv(s ...string) string {
	return PlusPathRep(os.Getenv, s...)
}

func PlusPathRep(f Replacer, s ...string) string {

	ss := []string{}
	for _, v := range s {
		ss = append(ss, ScanReplace(v, "{", "}", f))
	}
	return PlusPath(ss...)

}

//PlusPath Takes a set of strings to join as path.
//if first char of later strings is '/'
func PlusPath(ss ...string) string {

	res := ""
	//Work backwards
	for i := len(ss) - 1; i >= 0; i-- {
		v := ss[i]
		if len(v) == 0 {
			continue
		}
		res = path.Join(v, res)
		if v[0] == '/' {
			return res
		}
	}

	return res
}
