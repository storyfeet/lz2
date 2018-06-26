package lz2

import (
	"fmt"
	"os"
)

func Help(onzero bool, items ...string) bool {
	return HelpArgs(os.Args[1:], onzero, items...)
}

//HelpArgs returns if the firstitem was a request for help, and if it output the help
func HelpArgs(args []string, onzero bool, items ...string) bool {
	if len(args) > 1 {
		return false
	}
	if len(args) == 0 && !onzero {
		return false
	}
	if len(args) > 0 {
		found := false
		for _, v := range []string{"-help", "-h", "help", "h"} {
			if args[0] == v {
				found = true
			}
		}
		if !found {
			return false
		}
	}
	for k, v := range items {
		if k == 0 {
			fmt.Println(v)
			continue
		}
		fmt.Printf("\t%s\n", v)
	}
	return true
}

func (cf Config) Help(onzero bool, message string) bool {

	items := []string{message}
	for _, v := range cf.Helpdata {
		items = append(items, v)
	}
	return HelpArgs(cf.Args, onzero, items...)
}
