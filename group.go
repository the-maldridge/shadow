package shadow

import (
	"bufio"
	"io"
	"strconv"
	"strings"
)

// A GroupEntry is a single entry in the group map.  The entry uses
// the field names as found in `man 5 group`.
type GroupEntry struct {
	Name     string
	Password string
	GID      int
	UserList []string
}

// Parse reads a single entry of the group map.  Parsing will fail if
// a group has too many members to load in a single pass.
func (ge *GroupEntry) Parse(s string) error {
	fields := strings.FieldsFunc(s, func(c rune) bool { return c == ':' })
	if len(fields) != 4 {
		return ErrWrongNumFields
	}

	ge.Name = fields[0]
	ge.Password = fields[1]

	gid, err := strconv.Atoi(fields[2])
	if err != nil {
		*ge = GroupEntry{}
		return ErrNotANumber
	}
	ge.GID = gid

	ge.UserList = strings.FieldsFunc(fields[3], func(c rune) bool { return c == ',' })
	return nil
}

// A GroupMap is a complete list of groups that can be written and
// used by the system.
type GroupMap struct {
	lines []*GroupEntry
}

// ParseGroupMap loads from the specified reader into a list of
// GroupEntry.
func ParseGroupMap(r io.Reader) (*GroupMap, error) {
	lines := []*GroupEntry{}
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		t := new(GroupEntry)
		if err := t.Parse(scanner.Text()); err != nil {
			return nil, err
		}
		lines = append(lines, t)
	}
	gm := new(GroupMap)
	gm.lines = lines
	return gm, nil
}
