package visiontac

import (
	"strings"
	"testing"
	"time"
)

func TestParseInt(t *testing.T) {
	i, _ := parseInt("23\x00\x00\x00\x00")
	if i != 23 {
		t.Errorf("wrong int parsed: %v", i)
	}
}

func TestParseInvalidInt(t *testing.T) {
	_, err := parseInt("23\x001\x00\x00")
	if err == nil {
		t.Errorf("expected error from int")
	}
}

func TestParseFloat(t *testing.T) {
	f, _ := parseFloat("1.9\x00\x00")
	if f != 1.9 {
		t.Errorf("wrong float parsed: %v", f)
	}
}

func TestParseInvalidFloat(t *testing.T) {
	_, err := parseFloat("23\x001\x00\x00")
	if err == nil {
		t.Errorf("expected error from float")
	}
}

func TestParseTag(t *testing.T) {
	i, _ := parseTag("T")
	if i != 'T' {
		t.Errorf("wrong tag parsed: %v", i)
	}
}

func TestParseInvalidTag(t *testing.T) {
	_, err := parseTag("AB")
	if err == nil {
		t.Errorf("expected error from tag")
	}
}

func TestParseTimestamp(t *testing.T) {
	ts, _ := parseTimestamp("111213", "185059")
	if ts != time.Date(2011, 12, 13, 18, 50, 59, 0, time.UTC) {
		t.Errorf("wrong timestamp parsed: %v", ts)
	}
}

func TestParseInvalidTimestamp(t *testing.T) {
	_, err := parseTimestamp("111213", "1850591")
	if err == nil {
		t.Errorf("expected error from timestamp")
	}
}

func TestParseNorthernLatitude(t *testing.T) {
	lat, _ := parseCoordinate("36.874506N")
	if lat != 36.874506 {
		t.Errorf("wrong latitude parsed: %v", lat)
	}
}

func TestParseSouthernLatitude(t *testing.T) {
	lat, _ := parseCoordinate("36.874506S")
	if lat != -36.874506 {
		t.Errorf("wrong latitude parsed: %v", lat)
	}
}

func TestParseEasternLongitude(t *testing.T) {
	lon, _ := parseCoordinate("174.779188E")
	if lon != 174.779188 {
		t.Errorf("wrong longitude parsed: %v", lon)
	}
}

func TestParseWesternLongitude(t *testing.T) {
	lon, _ := parseCoordinate("174.779188W")
	if lon != -174.779188 {
		t.Errorf("wrong longitude parsed: %v", lon)
	}
}

func TestParseInvalidLatitudeLongitude(t *testing.T) {
	_, err := parseCoordinate("174.779188X")
	if err == nil {
		t.Errorf("expected error from coordinate 'X'")
	}
}

func TestParseStandardLine(t *testing.T) {
	input := "23\x00\x00\x00\x00,T,090512,041041,41.302453S,174.778450E,2\x00\x00,3\x00\x00\x00,1\x00\x00,\x00\x00\x00\x00\x00\x00\x00\x00\x00"

	wrapper, err := parseStandard(strings.Split(input, ","))
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	rec := wrapper.StandardRecord()

	if rec.Index != 23 {
		t.Errorf("wrong index parsed: %v", rec.Index)
	}
	if rec.Tag != 'T' {
		t.Errorf("wrong tag parsed: %v", rec.Tag)
	}
	if rec.Timestamp != time.Date(2009, 5, 12, 4, 10, 41, 0, time.UTC) {
		t.Errorf("wrong timestamp parsed: %v", rec.Timestamp)
	}
	if rec.Latitude != -41.302453 {
		t.Errorf("wrong latitude parsed: %v", rec.Latitude)
	}
	if rec.Longitude != 174.778450 {
		t.Errorf("wrong longitude parsed: %v", rec.Longitude)
	}
	if rec.Height != 2 {
		t.Errorf("wrong height parsed: %v", rec.Height)
	}
	if rec.Speed != 3 {
		t.Errorf("wrong speed parsed: %v", rec.Speed)
	}
	if rec.Heading != 1 {
		t.Errorf("wrong heading parsed: %v", rec.Heading)
	}
	// if rec.FixMode != nil {
	// 	t.Errorf("wrong fix mode parsed: %v", rec.FixMode)
	// }
	// if rec.Valid != nil {
	// 	t.Errorf("wrong valid parsed: %v", rec.Valid)
	// }
	// if rec.Pdop != nil {
	// 	t.Errorf("wrong pdop parsed: %v", rec.Pdop)
	// }
	// if rec.Hdop != nil {
	// 	t.Errorf("wrong hdop parsed: %v", rec.Hdop)
	// }
	// if rec.Vdop != nil {
	// 	t.Errorf("wrong vdop parsed: %v", rec.Vdop)
	// }
	if rec.Vox != "" {
		t.Errorf("wrong vox parsed: %v", rec.Vox)
	}
}

func TestParseAdvancedLine(t *testing.T) {
	input := "1\x00\x00\x00\x00\x00,T,111213,185059,36.874506S,174.779188E,152\x00\x00,79\x00\x00,120,3D,SPS ,2.1\x00\x00,1.9\x00\x00,1.0\x00\x00,\x00\x00\x00\x00\x00\x00\x00\x00\x00"

	wrapper, err := parseAdvanced(strings.Split(input, ","))
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	advwrapper, ok := wrapper.(AdvancedRecord)
	if !ok {
		t.Errorf("advanced record not parsed")
	}
	rec := advwrapper.AdvancedRecord()

	if rec.Index != 1 {
		t.Errorf("wrong index parsed: %v", rec.Index)
	}
	if rec.Tag != 'T' {
		t.Errorf("wrong tag parsed: %v", rec.Tag)
	}
	if rec.Timestamp != time.Date(2011, 12, 13, 18, 50, 59, 0, time.UTC) {
		t.Errorf("wrong timestamp parsed: %v", rec.Timestamp)
	}
	if rec.Latitude != -36.874506 {
		t.Errorf("wrong latitude parsed: %v", rec.Latitude)
	}
	if rec.Longitude != 174.779188 {
		t.Errorf("wrong longitude parsed: %v", rec.Longitude)
	}
	if rec.Height != 152 {
		t.Errorf("wrong height parsed: %v", rec.Height)
	}
	if rec.Speed != 79 {
		t.Errorf("wrong speed parsed: %v", rec.Speed)
	}
	if rec.Heading != 120 {
		t.Errorf("wrong heading parsed: %v", rec.Heading)
	}
	if rec.FixMode != "3D" {
		t.Errorf("wrong fix mode parsed: %v", rec.FixMode)
	}
	if rec.Valid != "SPS " {
		t.Errorf("wrong valid parsed: %v", rec.Valid)
	}
	if rec.Pdop != 2.1 {
		t.Errorf("wrong pdop parsed: %v", rec.Pdop)
	}
	if rec.Hdop != 1.9 {
		t.Errorf("wrong hdop parsed: %v", rec.Hdop)
	}
	if rec.Vdop != 1.0 {
		t.Errorf("wrong vdop parsed: %v", rec.Vdop)
	}
	if rec.Vox != "" {
		t.Errorf("wrong vox parsed: %v", rec.Vox)
	}
}

func TestParseStandardFile(t *testing.T) {
	input :=
		"INDEX,TAG,DATE,TIME,LATITUDE N/S,LONGITUDE E/W,HEIGHT,SPEED,HEADING,VOX\n" +
			"1\x00\x00\x00\x00\x00,T,090512,041041,41.302453S,174.778450E,2\x00\x00,3\x00\x00\x00,1\x00\x00,\x00\x00\x00\x00\x00\x00\x00\x00\x00\n" +
			"2\x00\x00\x00\x00\x00,T,090512,041041,41.302453S,174.778450E,2\x00\x00,3\x00\x00\x00,1\x00\x00,\x00\x00\x00\x00\x00\x00\x00\x00\x00"

	p, err := NewParser(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	recs, err := p.ParseAll()
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	if len(recs) != 2 {
		t.Fatalf("expected 2 records not %d", len(recs))
	}

	if _, ok := recs[0].(advrec); ok {
		t.Errorf("parsed as advanced record")
	}
	rec0 := recs[0].StandardRecord()
	if rec0.Index != 1 {
		t.Errorf("wrong index parsed: %v", rec0.Index)
	}
	if _, ok := recs[1].(advrec); ok {
		t.Errorf("parsed as advanced record")
	}
	rec1 := recs[1].StandardRecord()
	if rec1.Index != 2 {
		t.Errorf("wrong index parsed: %v", rec1.Index)
	}
}

func TestParseAdvancedFile(t *testing.T) {
	input :=
		"INDEX,TAG,DATE,TIME,LATITUDE N/S,LONGITUDE E/W,HEIGHT,SPEED,HEADING,FIX MODE,VALID,PDOP,HDOP,VDOP,VOX\n" +
			"1\x00\x00\x00\x00\x00,T,090512,041041,41.302453S,174.778450E,2\x00\x00,3\x00\x00\x00,1\x00\x00,3D,SPS ,1.3\x00\x00,1.0\x00\x00,0.9\x00\x00,\x00\x00\x00\x00\x00\x00\x00\x00\x00\n" +
			"2\x00\x00\x00\x00\x00,T,090512,041041,41.302453S,174.778450E,2\x00\x00,3\x00\x00\x00,1\x00\x00,3D,SPS ,1.7\x00\x00,0.8\x00\x00,1.5\x00\x00,\x00\x00\x00\x00\x00\x00\x00\x00\x00"

	p, err := NewParser(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	recs, err := p.ParseAll()
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	if len(recs) != 2 {
		t.Fatalf("expected 2 records not %d", len(recs))
	}

	if _, ok := recs[0].(advrec); !ok {
		t.Errorf("not parsed as advanced record")
	}
	rec0 := recs[0].StandardRecord()
	if rec0.Index != 1 {
		t.Errorf("wrong index parsed: %v", rec0.Index)
	}
	if _, ok := recs[1].(advrec); !ok {
		t.Errorf("not parsed as advanced record")
	}
	rec1 := recs[1].StandardRecord()
	if rec1.Index != 2 {
		t.Errorf("wrong index parsed: %v", rec1.Index)
	}
}
