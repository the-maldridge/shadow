package shadow

import (
	"bufio"
	"io"
	"strings"
)

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

// FilterUID applies a NumericFilter to the UID field of all loaded
// PasswdEntry's and returns a list of all entries that matched.
func (pm *PasswdMap) FilterUID(f NumericFilter) []*PasswdEntry {
	nl := []*PasswdEntry{}
	for _, l := range pm.lines {
		if !f(l.UID) {
			// Filter did not match.
			continue
		}
		nl = append(nl, l)
	}
	return nl
}

// Add adds new passwd entries to the existing map.  Uniqueness is not
// enforced.
func (pm *PasswdMap) Add(a []*PasswdEntry) {
	pm.lines = append(pm.lines, a...)
}

// Del iterates through the provided list and removes entities that
// are exactly the same from the existing map.  The provided set must
// not contain duplicate Login values, potentially necessitating two
// calls if you have entries that are identical except for login.
func (pm *PasswdMap) Del(d []*PasswdEntry) {
	checkMap := make(map[string]*PasswdEntry, len(d))

	for _, e := range d {
		checkMap[e.Login] = e
	}

	out := []*PasswdEntry{}
	for _, l := range pm.lines {
		e, doTest := checkMap[l.Login]
		if doTest && *l == *e {
			// The entity is an exact match and should be
			// removed.
			continue
		}
		// The entity is not an exact match, and should be
		// retained.
		out = append(out, l)
	}
	pm.lines = out
}
