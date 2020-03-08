package shadow

import (
	"io"
	"strings"
	"testing"
)

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
	}
	for i, c := range cases {
		gm := new(GroupMap)
		if err := gm.Parse(c.r); err != c.wantErr {
			t.Errorf("%d: Got %v; Want %v", i, err, c.wantErr)
		}
	}
}
