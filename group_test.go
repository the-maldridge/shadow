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

func TestFilterGID(t *testing.T) {
	gm := &GroupMap{
		lines: []*GroupEntry{
			&GroupEntry{
				Name: "group1",
				GID:  1,
			},
			&GroupEntry{
				Name: "group2",
				GID:  2,
			},
		},
	}

	res := gm.FilterGID(func(i int) bool { return i == 2 })
	if len(res) != 1 || res[0].Name != "group2" {
		t.Error("Filter applied incorrectly!")
	}
}

func TestGroupAdd(t *testing.T) {
	gm := &GroupMap{
		lines: []*GroupEntry{},
	}

	if len(gm.lines) > 0 {
		t.Error("Wrong base condition")
	}

	gm.Add([]*GroupEntry{&GroupEntry{Name: "foo"}})

	if len(gm.lines) != 1 {
		t.Error("Add failed")
	}
}

func TestGroupDel(t *testing.T) {
	gm := &GroupMap{
		lines: []*GroupEntry{
			&GroupEntry{
				Name: "group1",
				GID:  1,
			},
			&GroupEntry{
				Name: "group2",
				GID:  2,
			},
		},
	}

	gm.Del([]*GroupEntry{&GroupEntry{Name: "group1", GID: 1}})
	if len(gm.lines) != 1 || gm.lines[0].Name != "group2" {
		t.Logf("%v", gm.lines)
		t.Error("Incorrect delete")
	}
}
