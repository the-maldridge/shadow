package shadow

import (
	"bufio"
	"io"
	"strconv"
	"strings"
)

// A PasswdEntry represents a single entry in the passwd map.  The
// entry uses the field names as found in `man 5 passwd`.
type PasswdEntry struct {
	Login    string
	Password string
	UID      int
	GID      int
	Comment  string
	Home     string
	Shell    string
}

func (pe PasswdEntry) String() string {
	return "L: " + pe.Login +
		" P: " + pe.Password +
		" U: " + strconv.Itoa(pe.UID) +
		" G: " + strconv.Itoa(pe.GID) +
		" C: " + pe.Comment +
		" H: " + pe.Home +
		" S: " + pe.Shell
}

// Parse parses a single line into a PasswdEntry struct.  Errors
// are returned if the wrong number of fields are present in the input
// string, or if the string contains illegal characters such as
// newlines.
func (pe *PasswdEntry) Parse(s string) error {
	fields := strings.Split(s, ":")
	if len(fields) != 7 {
		return ErrWrongNumFields
	}

	pe.Login = fields[0]
	pe.Password = fields[1]
	pe.Comment = fields[4]
	pe.Home = fields[5]
	pe.Shell = fields[6]

	uid, err := strconv.Atoi(fields[2])
	if err != nil {
		*pe = PasswdEntry{}
		return ErrNotANumber
	}
	pe.UID = uid

	gid, err := strconv.Atoi(fields[3])
	if err != nil {
		*pe = PasswdEntry{}
		return ErrNotANumber
	}
	pe.GID = gid

	return nil
}

// A PasswdMap is a complete set of passwd entries that can be written
// and used as a list of entities on a system.
type PasswdMap struct {
	lines []*PasswdEntry
}

func (pm PasswdMap) String() string {
	b := new(strings.Builder)
	for _, l := range pm.lines {
		b.WriteString(l.String())
		b.WriteRune('\n')
	}
	return b.String()
}

// ParsePasswdMap loads a specified reader into a password map for
// manipulation.
func ParsePasswdMap(r io.Reader) (*PasswdMap, error) {
	lines := []*PasswdEntry{}
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		t := new(PasswdEntry)
		if err := t.Parse(scanner.Text()); err != nil {
			return nil, err
		}
		lines = append(lines, t)
	}
	pm := new(PasswdMap)
	pm.lines = lines
	return pm, nil
}
