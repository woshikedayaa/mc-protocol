package ping

import (
	"errors"
	"strconv"
	"strings"
)

type version struct {
	ver []int
}

// only support stable version string
// e.g. 1.18.2, 1.20 1.14.2 1.4
func newVersion(s string) (version, error) {
	split := strings.Split(s, ".")
	if len(split) <= 1 || len(split) > 3 {
		return version{}, errors.New("invalid version string")
	}
	v := version{ver: make([]int, 3)}
	for i := 0; i < len(split); i++ {
		atoi, err := strconv.Atoi(split[i])
		if err != nil {
			return version{}, err
		}
		v.ver[i] = atoi
	}
	return v, nil
}

func (v *version) Major() int {
	return v.ver[0]
}

func (v *version) Minor() int {
	return v.ver[1]
}

func (v *version) Patch() int {
	return v.ver[2]
}

var (
	versionTable = []struct {
		ver   int
		build func() IPing
	}{
		{7, func() IPing {
			return new(c17)
		}},
		{6, func() IPing {
			return new(c16)
		}},
		{4, func() IPing {
			return new(c1415)
		},
		},
	}
)

// chooseImpl will return a Client implementation
// choose by version.Minor
// more see versionTable
// if input invalid version ,it will return a noop implementation
func chooseImpl(v version) IPing {
	for i := 0; i < len(versionTable); i++ {
		if v.Minor() >= versionTable[i].ver {
			return versionTable[i].build()
		}
	}
	// return noop client
	return new(cNoop)
}
