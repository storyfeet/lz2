package lz2

import (
	"bufio"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type LZ struct {
	Name  string
	Deets map[string]string
}

// ParseLZ takes a ":" separated string, and puts the first section as the name, and the rest as ordered details belonging to the LZ object returned
//Used as part of the main Parsing Method, for files
func ParseLZ(s string, lcase bool) LZ {
	s = strings.TrimSpace(s)
	ss := strings.Split(s, ":")
	res := LZ{strings.TrimSpace(ss[0]), map[string]string{}}
	for k, v := range ss[1:] {
		res.Deets["ex"+strconv.Itoa(k)] = strings.TrimSpace(v)
	}
	return res
}

//Read Takes a reader in .lz format and converts it to a [] of LZ
//Errors are normally in case of missing colons etc.
func Read(r io.Reader, lcase bool) ([]LZ, error) {
	sc := bufio.NewScanner(r)
	res := []LZ{}
	var curr LZ

	errs := ErrGroup{}

	line := 0
	for sc.Scan() {
		line++
		t := sc.Text()
		tr := strings.TrimSpace(t)
		if len(tr) == 0 {
			continue
		}
		if tr[0] == '#' {
			continue
		}
		if tr[0] == t[0] {
			//New Entry
			curr = ParseLZ(tr, lcase)
			res = append(res, curr)
			continue
		}

		//Deets

		if curr.Deets == nil {
			errs = append(errs, LineErr{"No Object Defined", line})
			continue
		}

		ss := strings.SplitN(tr, ":", 2)
		if len(ss) != 2 {
			errs = append(errs, LineErr{"No Colon", line})
			continue
		}
		s := strings.TrimSpace(ss[0])
		if lcase {
			s = strings.ToLower(s)
		}
		curr.Deets[s] = strings.TrimSpace(ss[1])
	}

	if len(errs) > 0 {
		return res, errs
	}
	return res, sc.Err()
}

//ReadFile is a wrapper for Read, which takes a filename instead of a Reader
//Errors on read error as well format errors.
//If the error is a format error it will fulfil the interface{NErrs()int}
func ReadFile(fname string, lcase ...bool) (*Config, error) {
	lc := false
	if len(lcase) > 0 {
		lc = lcase[0]
	}
	f, err := os.Open(fname)
	if err != nil {
		return &Config{}, err
	}
	defer f.Close()
	lz, err := Read(f, lc)
	return &Config{Location: fname, LL: lz}, nil
}

//PString returns a string for the value matching the first string in ns from its properties. if none are found, this will return an error.
func (lz LZ) PString(ns ...string) (string, error) {
	for _, v := range ns {
		res, ok := lz.Deets[v]
		if ok {
			return res, nil
		}
	}
	return "", errors.New("Item not found")
}

//PStringAr, will return a slice of strings, for the propertyname followed by increasing increments of ns0, ns1, ns2 ...
//eg if ns is "holding", then "holding0", "holding1", etc until one is not found.
func (lz LZ) PStringAr(ns ...string) []string {
	res := []string{}
	for _, v := range ns {
		i := 0
		for {
			s, ok := lz.Deets[v+strconv.Itoa(i)]
			if !ok {
				if i != 0 {
					break
				}
				s, ok = lz.Deets[v]
				if !ok {
					break
				}

			}
			res = append(res, s)
			i++
		}
	}
	return res
}

//PInt will try to convert the result of PString into an Int, errors on not found or on conversion error
func (lz LZ) PInt(ns ...string) (int, error) {
	s, err := lz.PString(ns...)
	if err != nil {
		return 0, err
	}
	conv, err := strconv.Atoi(s)
	if err != nil {
		return 0, errors.Wrap(err, "Could not convert Item")
	}

	return conv, nil
}

//PBool will try to convert the result of PString into an Bool, errors on not found or on conversion error
func (lz LZ) PBool(ns ...string) (bool, error) {
	s, err := lz.PString(ns...)
	if err != nil {
		return false, err
	}

	conv, err := strconv.ParseBool(s)
	if err != nil {
		return false, errors.Wrap(err, "Could not convert Item")
	}

	return conv, nil
}

//PFloat will try to convert the result of PString into an float64, errors on not found or on conversion error
func (lz LZ) PFloat(ns ...string) (float64, error) {
	s, err := lz.PString(ns...)
	if err != nil {
		return 0, err
	}

	conv, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, errors.Wrap(err, "Could not convert Item")
	}

	return conv, nil
}

//PStringD takes and returns a default result in case of error otherwise acts as PString
func (lz LZ) PStringD(def string, ns ...string) string {
	r, err := lz.PString(ns...)
	if err != nil {
		return def
	}
	return r
}

//PIntD takes and returns a default result in case of error otherwise acts as PInt
func (lz LZ) PIntD(def int, ns ...string) int {
	r, err := lz.PInt(ns...)
	if err != nil {
		return def
	}
	return r
}

//PBoolD takes and returns a default result in case of error otherwise acts as PBool
func (lz LZ) PBoolD(def bool, ns ...string) bool {
	r, err := lz.PBool(ns...)
	if err != nil {
		return def
	}
	return r
}

//PFloatD takes and returns a default result in case of error otherwise acts as PFloat
func (lz LZ) PFloatD(def float64, ns ...string) float64 {
	r, err := lz.PFloat(ns...)
	if err != nil {
		return def
	}
	return r
}

//Finds an item in the list by name
func ByName(ll []LZ, s string) (LZ, bool) {
	s = strings.ToLower(s)
	for _, v := range ll {
		if strings.ToLower(v.Name) == s {
			return v, true
		}
	}
	return LZ{}, false
}
