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

func (ge GroupEntry) String() string {
	return ge.Name + ":" +
		ge.Password + ":" +
		strconv.Itoa(ge.GID) + ":" +
		strings.Join(ge.UserList, ",")
}

// Parse reads a single entry of the group map.  Parsing will fail if
// a group has too many members to load in a single pass.
func (ge *GroupEntry) Parse(s string) error {
	fields := strings.Split(s, ":")
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

func (gm GroupMap) String() string {
	out := new(strings.Builder)
	for _, l := range gm.lines {
		out.WriteString(l.String())
		out.WriteRune('\n')
	}
	return out.String()
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

// FilterGID applies a NumericFilter to the UID field of all loaded
// GroupEntry's and returns a list of all entries that matched.
func (gm *GroupMap) FilterGID(f NumericFilter) []*GroupEntry {
	ng := []*GroupEntry{}
	for _, l := range gm.lines {
		if !f(l.GID) {
			// Filter did not match.
			continue
		}
		ng = append(ng, l)
	}
	return ng
}

// Add adds new group entries to the existing map.  Uniqueness is not
// enforced.
func (gm *GroupMap) Add(a []*GroupEntry) {
	gm.lines = append(gm.lines, a...)
}

// Del iterates through the provided list and removes entities that
// are exactly the same from the existing map.  The provided set must
// not contain duplicate Login values, potentially necessitating two
// calls if you have entries that are identical except for login.
func (gm *GroupMap) Del(d []*GroupEntry) {
	checkMap := make(map[string]*GroupEntry, len(d))

	for _, e := range d {
		checkMap[e.Name] = e
	}

	out := []*GroupEntry{}
	for _, l := range gm.lines {
		e, doTest := checkMap[l.Name]
		if doTest && l.Name == e.Name && l.GID == e.GID {
			// The entity is a match and should be
			// removed.
			continue
		}
		// The entity is not an exact match, and should be
		// retained.
		out = append(out, l)
	}
	gm.lines = out
}
