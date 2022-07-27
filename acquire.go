package acquire

import (
	"errors"
	"os"
	"path/filepath"
)

type TargetType int

const (
	File TargetType = iota + 1
	Dir
	All
)

var ErrNotFound = errors.New("not found")

func Acquire(t TargetType, patterns ...string) (matches []string, err error) {
	return AcquireFromUnder(t, ".", "", patterns...)
}

func MustAcquire(t TargetType, patterns ...string) (matches []string) {
	matches, err := Acquire(t, patterns...)
	if err != nil {
		panic(err)
	}
	return matches
}

func AcquireUnder(t TargetType, under string, patterns ...string) (matches []string, err error) {
	return AcquireFromUnder(t, ".", under, patterns...)
}

func AcquireFromUnder(t TargetType, folder, under string, patterns ...string) (matches []string, err error) {
	if !filepath.IsAbs(folder) {
		folder, err = filepath.Abs(folder)
		if err != nil {
			return
		}
	}
	if under != "" && !filepath.IsAbs(under) {
		under, err = filepath.Abs(under)
		if err != nil {
			return
		}
	}
	for {
		for _, p := range patterns {
			tmpMatches, err := filepath.Glob(filepath.Join(folder, p))
			if err != nil {
				return nil, err
			}
			for _, m := range tmpMatches {
				if t == All {
					matches = append(matches, m)
				} else {
					info, err := os.Stat(m)
					if err != nil {
						panic(err) // should not be happened
					}
					if (t == File && !info.IsDir()) || (t == Dir && info.IsDir()) {
						matches = append(matches, m)
					}
				}
			}
		}
		if len(matches) > 0 {
			return
		}
		if folder == under {
			err = ErrNotFound
			return
		}
		next := filepath.Dir(folder)
		if next == folder {
			err = ErrNotFound
			return
		}
		folder = next
	}
}
