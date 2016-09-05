package visiontac

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Record struct {
	Index     int
	Tag       byte
	Timestamp time.Time
	Latitude  float32
	Longitude float32
	Height    int
	Speed     int
	Heading   int
	// FixMode   string
	// Valid     string
	// Pdop      float32
	// Hdop      float32
	// Vdop      float32
	Vox string
}

func parseInt(s string) (int, error) {
	return strconv.Atoi(strings.TrimRight(s, "\x00"))
}

func parseTag(s string) (byte, error) {
	if len(s) != 1 {
		return 0, errors.New("Expected tag of length 1")
	}
	return s[0], nil
}

func parseTimestamp(s1, s2 string) (time.Time, error) {
	return time.Parse("060102150405", s1+s2)
}

func parseLatitude(s string) (float32, error) {
	real := s[:len(s)-1]
	half := s[len(s)-1]
	var mult float32
	switch half {
	case 'N':
		mult = 1
	case 'S':
		mult = -1
	default:
		return 0, errors.New(fmt.Sprintf("unexpected direction %c", half))
	}
	result, err := strconv.ParseFloat(real, 32)
	return float32(result) * mult, err
}

func parseLongitude(s string) (float32, error) {
	real := s[:len(s)-1]
	half := s[len(s)-1]
	var mult float32
	switch half {
	case 'E':
		mult = 1
	case 'W':
		mult = -1
	default:
		return 0, errors.New(fmt.Sprintf("unexpected direction %c", half))
	}
	result, err := strconv.ParseFloat(real, 32)
	return float32(result) * mult, err
}

func Parse(s string) (Record, error) {
	vals := strings.Split(s, ",")
	fmt.Println(vals)
	rec := Record{}

	index, err := parseInt(vals[0])
	if err != nil {
		return rec, err
	}
	rec.Index = index

	tag, err := parseTag(vals[1])
	if err != nil {
		return rec, err
	}
	rec.Tag = tag

	ts, err := parseTimestamp(vals[2], vals[3])
	if err != nil {
		return rec, err
	}
	rec.Timestamp = ts

	lat, err := parseLatitude(vals[4])
	if err != nil {
		return rec, err
	}
	rec.Latitude = lat

	lon, err := parseLongitude(vals[5])
	if err != nil {
		return rec, err
	}
	rec.Longitude = lon

	height, err := parseInt(vals[6])
	if err != nil {
		return rec, err
	}
	rec.Height = height

	speed, err := parseInt(vals[7])
	if err != nil {
		return rec, err
	}
	rec.Speed = speed

	return rec, nil
}
