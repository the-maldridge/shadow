package shadow

import (
	"io"
	"strings"
	"testing"
	"time"
)

func TestShadowEntryString(t *testing.T) {
	x := ShadowEntry{
		Login:    "foo",
		Password: "*",
	}

	want := "L: foo P: * LC: 01 Jan 01 00:00 +0000 mPA: 0 MPA: 0 WD: 0 ID: 0 E: 01 Jan 01 00:00 +0000 R: "
	if x.String() != want {
		t.Errorf("Got: '%s'; Want: '%s'", x.String(), want)
	}
}

func TestParseShadowEntry(t *testing.T) {
	cases := []struct {
		line      string
		wantEntry ShadowEntry
		wantErr   error
	}{
		{
			line:      "",
			wantEntry: ShadowEntry{},
			wantErr:   ErrWrongNumFields,
		},
		{
			line: "nobody:x:17518:0:99999:7:::",
			wantEntry: ShadowEntry{
				Login:              "nobody",
				Password:           "x",
				LastChanged:        time.Date(2017, time.December, 18, 0, 0, 0, 0, time.UTC),
				MaximumPasswordAge: 99999,
				WarningDays:        7,
				Expiration:         epochStart,
			},
			wantErr: nil,
		},
	}

	for i, c := range cases {
		se := new(ShadowEntry)
		if err := se.Parse(c.line); err != c.wantErr {
			t.Errorf("%d: Got %v; Want %v", i, err, c.wantErr)
		}
		if *se != c.wantEntry {
			t.Errorf("%d: Got \n%+v; Want \n%+v", i, *se, c.wantEntry)
		}
	}
}

func TestShadowMapString(t *testing.T) {
	x := ShadowMap{
		lines: []*ShadowEntry{
			&ShadowEntry{
				Login:    "foo",
				Password: "*",
			},
			&ShadowEntry{
				Login:    "bar",
				Password: "!",
			},
		},
	}

	want := "L: foo P: * LC: 01 Jan 01 00:00 +0000 mPA: 0 MPA: 0 WD: 0 ID: 0 E: 01 Jan 01 00:00 +0000 R: \nL: bar P: ! LC: 01 Jan 01 00:00 +0000 mPA: 0 MPA: 0 WD: 0 ID: 0 E: 01 Jan 01 00:00 +0000 R: \n"
	if x.String() != want {
		t.Errorf("Got: '%s'; Want: '%s'", x.String(), want)
	}
}

func TestParseShadowMap(t *testing.T) {
	cases := []struct {
		r       io.Reader
		wantErr error
	}{
		{
			r:       strings.NewReader("\nplaceholder\n"),
			wantErr: ErrWrongNumFields,
		},
		{
			r:       strings.NewReader("nobody:x:17518:0:99999:7:::"),
			wantErr: nil,
		},
	}

	for i, c := range cases {
		if _, err := ParseShadowMap(c.r); err != c.wantErr {
			t.Errorf("%d: Got %v; Want %v", i, err, c.wantErr)
		}
	}
}
