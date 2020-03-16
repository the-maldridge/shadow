package shadow

import (
	"io"
	"strings"
	"testing"
)

func TestGroupEntryString(t *testing.T) {
	x := GroupEntry{
		Name:     "group",
		Password: "x",
		GID:      42,
		UserList: []string{"foo", "bar"},
	}

	want := "N: group P: x G: 42 M: foo,bar"
	if x.String() != want {
		t.Errorf("Got: '%s'; Want: '%s'", x.String(), want)
	}
}

func TestParseGroupEntry(t *testing.T) {
	cases := []struct {
		line     string
		wantErr  error
		wantName string
	}{
		{
			line:     "",
			wantErr:  ErrWrongNumFields,
			wantName: "",
		},
		{
			line:     "kvm:x:24:maldridge,libvirt",
			wantErr:  nil,
			wantName: "kvm",
		},
		{
			line:     "kvm:x:potato:maldridge,libvirt",
			wantErr:  ErrNotANumber,
			wantName: "",
		},
	}

	for i, c := range cases {
		ge := new(GroupEntry)
		if err := ge.Parse(c.line); err != c.wantErr {
			t.Errorf("%d: Got %v; Want %v", i, err, c.wantErr)
		}
		if ge.Name != c.wantName {
			t.Errorf("%d: Got Name %s; Want name %s", i, ge.Name, c.wantName)
		}
	}
}

func TestGroupMapString(t *testing.T) {
	x := GroupMap{
		lines: []*GroupEntry{
			&GroupEntry{
				Name:     "group",
				Password: "x",
				GID:      42,
				UserList: []string{"foo", "bar"},
			},
			&GroupEntry{
				Name:     "ungroup",
				Password: "x",
				GID:      43,
				UserList: []string{"bar", "baz"},
			},
		},
	}

	want := "N: group P: x G: 42 M: foo,bar\nN: ungroup P: x G: 43 M: bar,baz\n"
	if x.String() != want {
		t.Errorf("Got: '%s'; Want: '%s'", x.String(), want)
	}
}

func TestParseGroupMap(t *testing.T) {
	cases := []struct {
		r       io.Reader
		wantErr error
	}{
		{
			r:       strings.NewReader("\nplaceholder\n"),
			wantErr: ErrWrongNumFields,
		},
		{
			r:       strings.NewReader("kvm:x:24:maldridge,libvirt"),
			wantErr: nil,
		},
		{
			r:       strings.NewReader("kvm:x:24:"),
			wantErr: nil,
		},
	}
	for i, c := range cases {
		if _, err := ParseGroupMap(c.r); err != c.wantErr {
			t.Errorf("%d: Got %v; Want %v", i, err, c.wantErr)
		}
	}
}
