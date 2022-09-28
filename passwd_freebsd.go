package shadow

import (
	"strconv"
	"strings"
)

type PasswdEntry struct {
	Login    string
	Password string
	UID      int
	GID      int
	Class    string
	Change   string
	Expire   string
	Comment  string
	Home     string
	Shell    string
}

func (pe PasswdEntry) String() string {
	return pe.Login + ":" +
		pe.Password + ":" +
		strconv.Itoa(pe.UID) + ":" +
		strconv.Itoa(pe.GID) + ":" +
		pe.Class + ":" +
		pe.Change + ":" +
		pe.Expire + ":" +
		pe.Comment + ":" +
		pe.Home + ":" +
		pe.Shell
}

func (pe *PasswdEntry) Parse(s string) error {
	fields := strings.Split(s, ":")
	if len(fields) != 10 {
		return ErrWrongNumFields
	}

	pe.Login = fields[0]
	pe.Password = fields[1]
	pe.Class = fields[4]
	pe.Change = fields[5]
	pe.Expire = fields[6]
	pe.Comment = fields[7]
	pe.Home = fields[8]
	pe.Shell = fields[9]

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
