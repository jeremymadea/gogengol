package world

import (
	"math"
	"math/rand"
)

type World struct {
	w, h int
	grid [][]byte
}

// Returns a new world with width w, height h, and (approximate) population
// density pd.
func New(w, h int, pd float64) *World {
	w := new(World)
	w.w = w
	w.h = h
	w.grid = make([][]byte, w)
	for i := range w.grid {
		w.grid[i] = make([]byte, h)
		for j := range w.grid[i] {
			if rand.Float32() < pd {
				w.grid[i][j] = 1
			} else {
				w.grid[i][j] = 0
			}
		}
	}
    return w
}

// Returns a new empty world.
func NewEmpty(w, h int) *World {
	w := new(World)
	w.w = w
	w.h = h
	w.grid = make([][]byte, w)
	for i := range w.grid {
		w.grid[i] = make([]byte, h)
	}
    return w
}

// Returns a new world with a circular patch of the given radius populated
// with an approximate population density of pd.
func NewPopPatch(w, h, radius int, pd float64) *World {
	w := new(World)
	w.w = w
	w.h = h
	cx := w / 2
	cy := h / 2
	r := radius
	if r > cx || r > cy {
		r = int(math.Min(cx, cy))
	}
	r *= r // use radius squared for distance calculations.
	w.grid = make([][]byte, w)
	for i := range w.grid {
		w.grid[i] = make([]byte, h)
		for j := range w.grid[i] {
			// We don't care about actual distance, just relative distance so
			// we can avoid square roots.
			if (cx-i)*(cx-i)+(cy-j)*(cy-j) < r {
				if rand.Float32() < pd {
					w.grid[i][j] = 1
				} else {
					w.grid[i][j] = 0
				}
			}
		}
	}
	return w
}

func GetHood(w *World, x, y int) byte {
	hood := byte(0)
	Ly := y - 1
	Lx := x - 1
	Hy := y + 1
	Hx := x + 1
	if Hx >= w.w {
		Hx = 0
	}
	if Hy >= w.h {
		Hy = 0
	}
	if Lx <= 0 {
		Lx = w.w - 1
	}
	if Ly <= 0 {
		Ly = w.h - 1
	}
	hood |= w.grid[Lx][Ly] & 1
	hood <<= 1
	hood |= w.grid[x][Ly] & 1
	hood <<= 1
	hood |= w.grid[Hx][Ly] & 1
	hood <<= 1
	hood |= w.grid[Hx][y] & 1
	hood <<= 1
	hood |= w.grid[Hx][Hy] & 1
	hood <<= 1
	hood |= w.grid[x][Hy] & 1
	hood <<= 1
	hood |= w.grid[Lx][Hy] & 1
	hood <<= 1
	hood |= w.grid[Lx][y] & 1
	return hood
}
