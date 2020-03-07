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

// Parse parses a single line into a PasswdEntry struct.  Errors
// are returned if the wrong number of fields are present in the input
// string, or if the string contains illegal characters such as
// newlines.
func (pe *PasswdEntry) Parse(s string) error {
	fields := strings.FieldsFunc(s, func(c rune) bool { return c == ':' })
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
type PasswdMap []*PasswdEntry

// Parse loads a specified reader into a password map for
// manipulation.
func (pm *PasswdMap) Parse(r io.Reader) error {
	*pm = []*PasswdEntry{}
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		t := new(PasswdEntry)
		if err := t.Parse(scanner.Text()); err != nil {
			return err
		}
		*pm = append(*pm, t)
	}
	return nil
}
