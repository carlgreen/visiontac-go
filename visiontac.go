package visiontac

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

type StandardRecord interface {
	StandardRecord() stdrec
}

type AdvancedRecord interface {
	StandardRecord
	AdvancedRecord() advrec
}

type stdrec struct {
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

func (r stdrec) StandardRecord() stdrec {
	return r
}

type advrec struct {
	stdrec
	FixMode string
	Valid   string
	Pdop    float32
	Hdop    float32
	Vdop    float32
}

func (r advrec) StandardRecord() stdrec {
	return r.stdrec
}

func (r advrec) AdvancedRecord() advrec {
	return r
}

func parseInt(s string) (int, error) {
	return strconv.Atoi(strings.TrimRight(s, "\x00"))
}

func parseFloat(s string) (float32, error) {
	result, err := strconv.ParseFloat(strings.TrimRight(s, "\x00"), 32)
	if err != nil {
		return 0, err
	}
	return float32(result), err
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

func parseCoordinate(s string) (float32, error) {
	real := s[:len(s)-1]
	half := s[len(s)-1]
	var mult float32
	switch half {
	case 'N', 'E':
		mult = 1
	case 'S', 'W':
		mult = -1
	default:
		return 0, errors.New(fmt.Sprintf("unexpected direction %c", half))
	}
	result, err := strconv.ParseFloat(real, 32)
	return float32(result) * mult, err
}

func parseStandard(vals []string) (StandardRecord, error) {
	rec := stdrec{}

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

	lat, err := parseCoordinate(vals[4])
	if err != nil {
		return rec, err
	}
	rec.Latitude = lat

	lon, err := parseCoordinate(vals[5])
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

	heading, err := parseInt(vals[8])
	if err != nil {
		return rec, err
	}
	rec.Heading = heading

	return rec, nil
}

func parseAdvanced(vals []string) (AdvancedRecord, error) {
	advrec := advrec{}
	rec, err := parseStandard(vals)
	if err != nil {
		return advrec, nil
	}
	advrec.stdrec = rec.StandardRecord()

	fixMode := vals[9]
	advrec.FixMode = fixMode

	valid := vals[10]
	advrec.Valid = valid

	pdop, err := parseFloat(vals[11])
	if err != nil {
		return advrec, err
	}
	advrec.Pdop = pdop

	hdop, err := parseFloat(vals[12])
	if err != nil {
		return advrec, err
	}
	advrec.Hdop = hdop

	vdop, err := parseFloat(vals[13])
	if err != nil {
		return advrec, err
	}
	advrec.Vdop = vdop

	return advrec, nil
}

type Parser struct {
	s *bufio.Scanner
}

func NewParser(r io.Reader) *Parser {
	return &Parser{
		s: bufio.NewScanner(r),
	}
}

func (r *Parser) Parse() (record StandardRecord, err error) {
	for r.s.Scan() {
		record, err = parse(r.s.Text())
		if record != nil {
			break
		}
		if err != nil {
			return nil, err
		}
	}
	if err := r.s.Err(); err != nil {
		return nil, err
	}

	return record, nil
}

func (r *Parser) ParseAll() (records []StandardRecord, err error) {
	for {
		record, err := r.Parse()
		if err != nil {
			return nil, err
		}
		if record == nil {
			return records, nil
		}
		records = append(records, record)
	}
}

func parse(s string) (StandardRecord, error) {
	vals := strings.Split(s, ",")
	switch len(vals) {
	case 10:
		return parseStandard(vals)
	case 15:
		return parseAdvanced(vals)
	default:
		return nil, errors.New("unexpected number of fields")
	}
}
