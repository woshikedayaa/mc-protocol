package ping

import (
	"fmt"
	"testing"
)

func TestNewVersion(t *testing.T) {
	tests := []string{
		"1.18.2", "1.7", "1.19555.5", "1.16", "1.4",
	}
	for i := 0; i < len(tests); i++ {
		v, err := newVersion(tests[i])
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(v)
	}
}
