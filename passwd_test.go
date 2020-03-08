package shadow

import (
	"io"
	"strings"
	"testing"
)

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
