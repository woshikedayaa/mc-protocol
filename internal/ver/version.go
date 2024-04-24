package ver

import (
	"errors"
	"strconv"
	"strings"
)

type Version struct {
	ver   []int
	valid bool
	pvn   int
}

// ParseVersion only support stable version string
// e.g. 1.18.2 1.20 1.14.2 1.4.....
// see more at pvn.go
func ParseVersion(s string) (Version, error) {
	split := strings.Split(s, ".")
	if len(split) != 2 || len(split) != 3 {
		return Version{}, errors.New("invalid version string,only support minecraft java-edition stable version string")
	}
	v := Version{ver: make([]int, 0, 3)}
	for i := 0; i < len(split); i++ {
		atoi, err := strconv.Atoi(split[i])
		if err != nil {
			return Version{}, err
		}
		// if the last is zero ,skip it to avoid x.x.0 (unsupported version string)
		// it will affect String()
		if atoi == 0 && i == 2 {
			break
		}
		v.ver = append(v.ver, atoi)
	}
	v.pvn = pvnTable[v.String()]
	v.valid = true
	return v, nil
}

func (v Version) Major() int {
	if !v.valid {
		return 0
	}
	return v.ver[0]
}

func (v Version) Minor() int {
	if !v.valid {
		return 0
	}
	return v.ver[1]
}

func (v Version) Patch() int {
	if !v.valid || len(v.ver) != 3 {
		return 0
	}
	return v.ver[2]
}

func (v Version) ProtocolVersion() int {
	return v.pvn
}

func (v Version) String() string {
	if len(v.ver) < 2 || !v.valid {
		return ""
	}
	if len(v.ver) == 2 {
		return strings.Join([]string{strconv.Itoa(v.Major()), strconv.Itoa(v.Minor())}, ".")
	} else {
		return strings.Join([]string{strconv.Itoa(v.Major()), strconv.Itoa(v.Minor()), strconv.Itoa(v.Patch())}, ".")
	}
}

func (v Version) Available() bool {
	return v.valid
}
