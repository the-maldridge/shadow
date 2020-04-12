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

	want := "foo:*:::::::"
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
			// This struct is now wrong.  Needs to be
			// fixed to handle that there are now a ton of
			// flags.
			line: "nobody:x:17518:0:99999:7:::",
			wantEntry: ShadowEntry{
				LastChanged: time.Date(2017, time.December, 18, 0, 0, 0, 0, time.UTC),
				Expiration:  epochStart,

				Login:                 "nobody",
				Password:              "x",
				MinimumPasswordAge:    0,
				MaximumPasswordAge:    99999,
				WarningDays:           7,
				InactivityDays:        0,
				Reserved:              "",
				HasLastChanged:        true,
				HasMinimumPasswordAge: true,
				HasMaximumPasswordAge: true,
				HasWarningDays:        true,
				HasInactivityDays:     false,
				HasExpiration:         false,
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
			t.Errorf("%d: Got \n%v; Want \n%v", i, *se, c.wantEntry)
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

	want := "foo:*:::::::\nbar:!:::::::\n"
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

func TestFilterLogin(t *testing.T) {
	pm := &ShadowMap{
		lines: []*ShadowEntry{
			&ShadowEntry{
				Login: "login1",
			},
			&ShadowEntry{
				Login: "login2",
			},
		},
	}

	res := pm.FilterUID(func(s string) bool { return s == "login2" })
	if len(res) != 1 || res[0].Login != "login2" {
		t.Error("Filter applied incorrectly!")
	}
}

func TestShadowAdd(t *testing.T) {
	pm := &ShadowMap{
		lines: []*ShadowEntry{},
	}

	if len(pm.lines) > 0 {
		t.Error("Wrong base condition")
	}

	pm.Add([]*ShadowEntry{&ShadowEntry{Login: "foo"}})

	if len(pm.lines) != 1 {
		t.Error("Add failed")
	}
}

func TestShadowDel(t *testing.T) {
	pm := &ShadowMap{
		lines: []*ShadowEntry{
			&ShadowEntry{
				Login: "login1",
			},
			&ShadowEntry{
				Login: "login2",
			},
		},
	}

	pm.Del([]*ShadowEntry{&ShadowEntry{Login: "login1"}})
	if len(pm.lines) != 1 || pm.lines[0].Login != "login2" {
		t.Logf("%v", pm.lines)
		t.Error("Incorrect delete")
	}
}
