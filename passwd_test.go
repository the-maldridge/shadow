package shadow

import (
	"io"
	"strings"
	"testing"
)

func TestPasswdEntryString(t *testing.T) {
	pe := &PasswdEntry{
		Login:    "foo",
		Password: "x",
		UID:      2,
		GID:      2,
		Comment:  "Foo",
		Home:     "/home/foo",
		Shell:    "/bin/fooshell",
	}

	want := "L: foo P: x U: 2 G: 2 C: Foo H: /home/foo S: /bin/fooshell"
	if pe.String() != want {
		t.Errorf("Want '%s'; Got '%s'", want, pe.String())
	}
}

func TestPasswdEntryParse(t *testing.T) {
	cases := []struct {
		line    string
		entry   PasswdEntry
		wantErr error
	}{
		{
			line:    "",
			entry:   PasswdEntry{},
			wantErr: ErrWrongNumFields,
		},
		{
			line: "maldridge:x:1000:1000:maldridge:/home/maldridge:/bin/bash",
			entry: PasswdEntry{
				Login:    "maldridge",
				Password: "x",
				UID:      1000,
				GID:      1000,
				Comment:  "maldridge",
				Home:     "/home/maldridge",
				Shell:    "/bin/bash",
			},
			wantErr: nil,
		},
		{
			line:    "maldridge:x:potato:1000:maldridge:/home/maldridge:/bin/bash",
			entry:   PasswdEntry{},
			wantErr: ErrNotANumber,
		},
		{
			line:    "maldridge:x:1000:potato:maldridge:/home/maldridge:/bin/bash",
			entry:   PasswdEntry{},
			wantErr: ErrNotANumber,
		},
	}

	for i, c := range cases {
		p := new(PasswdEntry)
		if err := p.Parse(c.line); err != c.wantErr {
			t.Errorf("%d: Got %v Want %v", i, err, c.wantErr)
		}
		if *p != c.entry {
			t.Errorf("%d: Got %v Want %v", i, p, c.entry)
		}
	}
}

func TestPasswdMapString(t *testing.T) {
	x := PasswdMap{
		lines: []*PasswdEntry{
			&PasswdEntry{
				Login:    "foo",
				Password: "x",
				UID:      2,
				GID:      2,
				Comment:  "Foo",
				Home:     "/home/foo",
				Shell:    "/bin/fooshell",
			},
			&PasswdEntry{
				Login:    "bar",
				Password: "x",
				UID:      3,
				GID:      3,
				Comment:  "Bar",
				Home:     "/home/foo",
				Shell:    "/bin/fooshell",
			},
		},
	}

	want := "L: foo P: x U: 2 G: 2 C: Foo H: /home/foo S: /bin/fooshell\nL: bar P: x U: 3 G: 3 C: Bar H: /home/foo S: /bin/fooshell\n"
	if x.String() != want {
		t.Errorf("Got: '%s'; Want '%s'", x.String(), want)
	}
}

func TestParsePasswdMap(t *testing.T) {
	cases := []struct {
		r       io.Reader
		wantErr error
	}{
		{
			r:       strings.NewReader("\nplaceholder\n"),
			wantErr: ErrWrongNumFields,
		},
		{
			r:       strings.NewReader("maldridge:x:1000:1000:maldridge:/home/maldridge:/bin/bash"),
			wantErr: nil,
		},
	}
	for i, c := range cases {
		if _, err := ParsePasswdMap(c.r); err != c.wantErr {
			t.Errorf("%d: Got %v; Want %v", i, err, c.wantErr)
		}
	}
}

func TestFilterUID(t *testing.T) {
	pm := &PasswdMap{
		lines: []*PasswdEntry{
			&PasswdEntry{
				Login: "login1",
				UID:   1,
			},
			&PasswdEntry{
				Login: "login2",
				UID:   2,
			},
		},
	}

	res := pm.FilterUID(func(i int) bool { return i == 2 })
	if len(res) != 1 || res[0].Login != "login2" {
		t.Error("Filter applied incorrectly!")
	}
}

func TestPasswdAdd(t *testing.T) {
	pm := &PasswdMap{
		lines: []*PasswdEntry{},
	}

	if len(pm.lines) > 0 {
		t.Error("Wrong base condition")
	}

	pm.Add([]*PasswdEntry{&PasswdEntry{Login: "foo"}})

	if len(pm.lines) != 1 {
		t.Error("Add failed")
	}
}

func TestPasswdDel(t *testing.T) {
	pm := &PasswdMap{
		lines: []*PasswdEntry{
			&PasswdEntry{
				Login: "login1",
				UID:   1,
			},
			&PasswdEntry{
				Login: "login2",
				UID:   2,
			},
		},
	}

	pm.Del([]*PasswdEntry{&PasswdEntry{Login: "login1", UID: 1}})
	if len(pm.lines) != 1 || pm.lines[0].Login != "login2" {
		t.Logf("%v", pm.lines)
		t.Error("Incorrect delete")
	}
}
