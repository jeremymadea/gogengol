package rules

import (
	"fmt"
	"math/rand"
	"strings"
)

type Rule struct {
	dead [256]byte
	live [256]byte
}

func (r *Rule) String() string {
	var b strings.Builder
	b.WriteString("    \"A\": \"")
	for i := range r.dead {
		fmt.Fprintf(&b, "%d", r.dead[i])
	}
	b.WriteString("\"\n    \"B\": \"")
	for i := range r.live {
		fmt.Fprintf(&b, "%d", r.live[i])
	}
	b.WriteString("\"\n")
	return b.String()
}

// Returns a new random rule with the approximate dead/live densities.
func NewRandomRule(rdd, rdl float32) *Rule {
	r := new(Rule)
	for i := range r.dead {
		if rand.Float32() < rdd {
			r.dead[i] = 1
		} else {
			r.dead[i] = 0
		}
	}
	for i := range r.live {
		if rand.Float32() < rdl {
			r.live[i] = 1
		} else {
			r.live[i] = 0
		}
	}
	return r
}

func NewFromStrings(dead, live string) *Rule {
	if len(dead) != 256 || len(live) != 256 {
		panic("String length not 256.")
	}
	r = new(Rule)
	for i, c := range dead {
		if c == '0' {
			r.dead[i] = 0
		} else {
			r.dead[i] = 1
		}
	}
	for i, c := range live {
		if c == '0' {
			r.live[i] = 0
		} else {
			r.live[i] = 1
		}
	}
	return r
}
