package lz2

import (
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"
)

type Config struct {
	Location string
	LL       []LZ
	Args     []string
	Helpdata []string
}

func LoadConfig(flagname string, lcase bool, locs ...string) (*Config, error) {
	return LoadConfigArgs(flagname, lcase, os.Args[1:], locs...)
}

func LoadConfigArgs(flagname string, lcase bool, args []string, locs ...string) (*Config, error) {
	fg, ok := FlagArgs(flagname, args...)
	if ok {
		cf, err := ReadFile(fg, lcase)
		cf.Args = args
		return cf, err
	}
	res, err := firstConfig(lcase, locs...)
	res.Args = args
	return res, err
}

func firstConfig(lcase bool, locs ...string) (*Config, error) {

	for _, v := range locs {
		floc := EnvReplace(v)
		cf, err := ReadFile(floc, lcase)
		if err != nil {
			if _, ok := err.(interface {
				NErrs() int
			}); ok {
				return cf, err
			}
			continue
		}
		//happy
		return cf, nil
	}

	//good let know what files were searched for in error
	estr := ""
	for _, v := range locs {
		floc := EnvReplace(v)
		estr += v + ":" + floc + ",\n"
	}

	return &Config{}, errors.Errorf("Config not found : " + estr)
}

func (cf *Config) GetS(fg string, cname string) (string, bool) {
	return cf.Flag(fg, "", cname)
}

func (cf *Config) Flag(fg, meaning string, cnames ...string) (string, bool) {
	res, ok := FlagArgs(fg, cf.Args...)
	cf.Helpdata = append(cf.Helpdata, fmt.Sprintf("-%s : %v : %s", fg, cnames, meaning))
	if ok {
		return res, true
	}

	for _, v := range cnames {
		res, ok := cf.Value(v)
		if ok {
			return res, true
		}
	}
	return "", false
}

func (cf *Config) FlagDef(fg, def, meaning string, cnames ...string) string {
	res, ok := cf.Flag(fg, fmt.Sprintf("%s [def = %s]", meaning, def), cnames...)
	if !ok {
		return def
	}
	return res
}

func (cf *Config) Value(s string) (string, bool) {
	if len(cf.LL) == 0 {
		return "", false
	}
	ss := strings.Split(s, ".")

	if len(ss) == 2 {
		for _, v := range cf.LL {
			if v.Name == ss[0] {
				res, ok := v.Deets[ss[1]]
				return res, ok
			}
		}
		return "", false
	}

	res, ok := cf.LL[0].Deets[s]
	return res, ok
}
