package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

const iso8601 = "2006-01-02T15:04:05-07:00"

func checkErr(tag string, err error, fatal bool) bool {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", tag, err.Error())
		if fatal {
			os.Exit(1)
		}
		return true
	}
	return false
}

func init() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %[1]s [OPTION...] FILE...\n", os.Args[0])
		fmt.Fprintln(flag.CommandLine.Output(), "Changes timestamps of files.")
		flag.PrintDefaults()
	}
}

func main() {
	ref := flag.String("r", "", "reference file, current time is used if not reference file present")
	changeModified := flag.Bool("m", true, "change modified time")
	changeAccessed := flag.Bool("a", true, "change accessed time")
	newDateStr := flag.String("d", "", fmt.Sprintf("new date in ISO-8601 (%s)", iso8601))
	flag.Parse()

	files := flag.Args()
	if len(files) == 0 {
		fmt.Fprintln(os.Stderr, "please specify at least one file")
		os.Exit(2)
	}

	var refTimeDef *timeDef
	if *ref != "" {
		r, err := getDate(*ref)
		checkErr("get ref time", err, true)
		refTimeDef = r
	} else {
		now := time.Now()
		refTimeDef = &timeDef{
			Accessed: now,
			Modified: now,
		}
	}

	if *newDateStr == "" {
		*newDateStr = "+0s"
	}
	var newDate *timeDef
	switch (*newDateStr)[0] {
	case '-', '+':
		// relative time duration
		dur, err := time.ParseDuration(*newDateStr)
		checkErr("resolve relative time", err, true)
		newDate = &timeDef{
			Accessed: refTimeDef.Accessed.Add(dur),
			Modified: refTimeDef.Modified.Add(dur),
		}
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		t, err := time.Parse(iso8601, *newDateStr)
		checkErr("resolve ISO8601 time", err, true)
		newDate = &timeDef{
			Accessed: t,
			Modified: t,
		}
	default:
		checkErr("resolve time", fmt.Errorf("unknown time format: %s", *newDateStr), true)
	}

	hasError := false
	for _, fn := range files {
		newDateForThisFile := *newDate
		if !*changeAccessed || !*changeModified {
			oldDate, err := getDate(fn)
			if checkErr(fn, err, false) {
				hasError = true
				continue
			}
			if !*changeAccessed {
				newDateForThisFile.Accessed = oldDate.Accessed
			}
			if !*changeModified {
				newDateForThisFile.Modified = oldDate.Modified
			}
		}
		if checkErr(fn, updateDate(fn, newDateForThisFile), false) {
			hasError = true
		}
	}

	if hasError {
		os.Exit(3)
	}
}
