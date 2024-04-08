package ver

import (
	"errors"
	"strconv"
	"strings"
)

type Version struct {
	ver []int
}

// NewVersion only support stable version string
// e.g. 1.18.2 1.20 1.14.2 1.4.....
func NewVersion(s string) (Version, error) {
	split := strings.Split(s, ".")
	if len(split) <= 1 || len(split) > 3 {
		return Version{}, errors.New("invalid version string")
	}
	v := Version{ver: make([]int, 3)}
	for i := 0; i < len(split); i++ {
		atoi, err := strconv.Atoi(split[i])
		if err != nil {
			return Version{}, err
		}
		v.ver[i] = atoi
	}
	return v, nil
}

func (v Version) Major() int {
	return v.ver[0]
}

func (v Version) Minor() int {
	return v.ver[1]
}

func (v Version) Patch() int {
	return v.ver[2]
}

func (v Version) ProtocolVersion() int {

}

func (v Version) String() string {
	if len(v.ver) < 2 {
		return ""
	}
	if len(v.ver) == 2 {
		return strings.Join([]string{strconv.Itoa(v.Major()), strconv.Itoa(v.Minor())}, ".")
	} else {
		return strings.Join([]string{strconv.Itoa(v.Major()), strconv.Itoa(v.Minor()), strconv.Itoa(v.Patch())}, ".")
	}
}
