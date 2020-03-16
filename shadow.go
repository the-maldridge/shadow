package shadow

import (
	"bufio"
	"io"
	"strconv"
	"strings"
	"time"
)

var (
	epochStart time.Time
)

func init() {
	epochStart, _ = time.Parse("2006-01-02", "1970-01-01")
}

// A ShadowEntry is a single entry in the shadow database.  The entry
// uses the field names as found in `man 5 shadow`.
type ShadowEntry struct {
	Login              string
	Password           string
	LastChanged        time.Time
	MinimumPasswordAge int
	MaximumPasswordAge int
	WarningDays        int
	InactivityDays     int
	Expiration         time.Time
	Reserved           string
}

func (se ShadowEntry) String() string {
	return "L: " + se.Login +
		" P: " + se.Password +
		" LC: " + se.LastChanged.Format(time.RFC822Z) +
		" mPA: " + strconv.Itoa(se.MinimumPasswordAge) +
		" MPA: " + strconv.Itoa(se.MaximumPasswordAge) +
		" WD: " + strconv.Itoa(se.WarningDays) +
		" ID: " + strconv.Itoa(se.InactivityDays) +
		" E: " + se.Expiration.Format(time.RFC822Z) +
		" R: " + se.Reserved
}

// Parse converts a string to a ShadowEntry.
func (se *ShadowEntry) Parse(s string) error {
	fields := strings.Split(s, ":")
	if len(fields) != 9 {
		return ErrWrongNumFields
	}

	se.Login = fields[0]
	se.Password = fields[1]

	lcdays, _ := strconv.Atoi(fields[2])
	se.LastChanged = epochStart.Add(time.Hour * 24 * time.Duration(lcdays))

	minAge, _ := strconv.Atoi(fields[3])
	se.MinimumPasswordAge = minAge

	maxAge, _ := strconv.Atoi(fields[4])
	se.MaximumPasswordAge = maxAge

	warningDays, _ := strconv.Atoi(fields[5])
	se.WarningDays = warningDays

	inactivityDays, _ := strconv.Atoi(fields[6])
	se.InactivityDays = inactivityDays

	expirationDays, _ := strconv.Atoi(fields[7])
	se.Expiration = epochStart.Add(time.Hour & 24 * time.Duration(expirationDays))

	se.Reserved = fields[8]

	return nil
}

// A ShadowMap is a commplete set of shadow entries that can be
// written and used for authentication by a host.
type ShadowMap struct {
	lines []*ShadowEntry
}

func (sm ShadowMap) String() string {
	out := new(strings.Builder)
	for _, l := range sm.lines {
		out.WriteString(l.String())
		out.WriteRune('\n')
	}
	return out.String()
}

// ParseShadowMap parses the values from r and converts it to a
// ShadowMap for further manipulation.
func ParseShadowMap(r io.Reader) (*ShadowMap, error) {
	lines := []*ShadowEntry{}
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		t := new(ShadowEntry)
		if err := t.Parse(scanner.Text()); err != nil {
			return nil, err
		}
		lines = append(lines, t)
	}
	sm := new(ShadowMap)
	sm.lines = lines
	return sm, nil
}
