package lz2

import "os"

func FlagArgsDef(fname, def string, args ...string) string {
	res, ok := FlagArgs(fname, args...)
	if !ok {
		return def
	}
	return res
}

func FlagArgs(fname string, args ...string) (string, bool) {
	if fname == "" {
		return "", false
	}
	compst := "-" + fname
	//TODO handle minus numbers / escape -
	found := false
	for _, v := range args {
		if found {
			return v, true
		}
		if compst == v {
			found = true
		}
	}
	return "", found
}

func Flag(fname string) (string, bool) {
	return FlagArgs(fname, os.Args[1:]...)
}

func FlagDef(fname, def string) string {
	return FlagArgsDef(fname, def, os.Args[1:]...)
}
