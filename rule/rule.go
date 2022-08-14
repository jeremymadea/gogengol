package rule

import (
	"fmt"
	"math/rand"
	"strings"
)

type Rule struct {
	Dead [256]byte
	Live [256]byte
}

func (r *Rule) String() string {
	var b strings.Builder
	b.WriteString("    \"A\": \"")
	for i := range r.Dead {
		fmt.Fprintf(&b, "%d", r.Dead[i])
	}
	b.WriteString("\"\n    \"B\": \"")
	for i := range r.Live {
		fmt.Fprintf(&b, "%d", r.Live[i])
	}
	b.WriteString("\"\n")
	return b.String()
}

// Returns a new random rule with the approximate dead/live densities.
func NewRandom(rdd, rdl float32) *Rule {
	r := new(Rule)
	for i := range r.Dead {
		if rand.Float32() < rdd {
			r.Dead[i] = 1
		} else {
			r.Dead[i] = 0
		}
	}
	for i := range r.Live {
		if rand.Float32() < rdl {
			r.Live[i] = 1
		} else {
			r.Live[i] = 0
		}
	}
	return r
}

func NewFromStrings(dead, live string) *Rule {
	if len(dead) != 256 || len(live) != 256 {
		panic("String length not 256.")
	}
	r := new(Rule)
	for i, c := range dead {
		if c == '0' {
			r.Dead[i] = 0
		} else {
			r.Dead[i] = 1
		}
	}
	for i, c := range live {
		if c == '0' {
			r.Live[i] = 0
		} else {
			r.Live[i] = 1
		}
	}
	return r
}
