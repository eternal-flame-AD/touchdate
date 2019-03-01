package main

import (
	"errors"
	"path/filepath"
	"time"

	"golang.org/x/sys/unix"
)

type timeDef struct {
	Accessed time.Time
	Modified time.Time
}

func updateDate(f string, d timeDef) error {
	f, err := filepath.Abs(f)
	if err != nil {
		return err
	}
	return unix.UtimesNano(f, []unix.Timespec{
		unix.Timespec{
			Sec:  d.Accessed.Unix(),
			Nsec: d.Accessed.UnixNano() - d.Accessed.Unix()*1000000000,
		},
		unix.Timespec{
			Sec:  d.Modified.Unix(),
			Nsec: d.Modified.UnixNano() - d.Modified.Unix()*1000000000,
		},
	})
}

func getDate(f string) (*timeDef, error) {
	if f == "" {
		return nil, errors.New("file not specified")
	}
	f, err := filepath.Abs(f)
	if err != nil {
		return nil, err
	}
	stat := new(unix.Stat_t)
	if err := unix.Stat(f, stat); err != nil {
		return nil, err
	}
	return &timeDef{
		Accessed: time.Unix(stat.Atim.Sec, stat.Atim.Nsec),
		Modified: time.Unix(stat.Mtim.Sec, stat.Mtim.Nsec),
	}, nil
}
