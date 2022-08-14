package world

import (
	"math"
	"math/rand"
)

type World struct {
	W, H int
	Grid [][]byte
}

// Returns a new world with width w, height h, and (approximate) population
// density pd.
func New(w, h int, pd float64) *World {
	wld := new(World)
	wld.W = w
	wld.H = h
	wld.Grid = make([][]byte, w)
	for i := range wld.Grid {
		wld.Grid[i] = make([]byte, h)
		for j := range wld.Grid[i] {
			if rand.Float64() < pd {
				wld.Grid[i][j] = 1
			} else {
				wld.Grid[i][j] = 0
			}
		}
	}
    return wld
}

// Returns a new empty world.
func NewEmpty(w, h int) *World {
	wld := new(World)
	wld.W = w
	wld.H = h
	wld.Grid = make([][]byte, w)
	for i := range wld.Grid {
		wld.Grid[i] = make([]byte, h)
	}
    return wld
}

// Returns a new world with a circular patch of the given radius populated
// with an approximate population density of pd.
func NewPopPatch(w, h, radius int, pd float64) *World {
	wld := new(World)
	wld.W = w
	wld.H = h
	cx := float64(w / 2)
	cy := float64(h / 2)
	r := float64(radius)
	if r > cx || r > cy {
		r = math.Min(cx, cy)
	}
	r *= r // use radius squared for distance calculations.
	wld.Grid = make([][]byte, w)
	for i := range wld.Grid {
		wld.Grid[i] = make([]byte, h)
		for j := range wld.Grid[i] {
			// We don't care about actual distance, just relative distance so
			// we can avoid square roots.
            cxf := cx - float64(i)
            cyf := cy - float64(j)
			if cxf*cxf+cyf*cyf < r {
				if rand.Float64() < pd {
					wld.Grid[i][j] = 1
				} else {
					wld.Grid[i][j] = 0
				}
			}
		}
	}
	return wld
}

func GetHood(w *World, x, y int) byte {
	hood := byte(0)
	Ly := y - 1
	Lx := x - 1
	Hy := y + 1
	Hx := x + 1
	if Hx >= w.W {
		Hx = 0
	}
	if Hy >= w.H {
		Hy = 0
	}
	if Lx <= 0 {
		Lx = w.W - 1
	}
	if Ly <= 0 {
		Ly = w.H - 1
	}
	hood |= w.Grid[Lx][Ly] & 1
	hood <<= 1
	hood |= w.Grid[x][Ly] & 1
	hood <<= 1
	hood |= w.Grid[Hx][Ly] & 1
	hood <<= 1
	hood |= w.Grid[Hx][y] & 1
	hood <<= 1
	hood |= w.Grid[Hx][Hy] & 1
	hood <<= 1
	hood |= w.Grid[x][Hy] & 1
	hood <<= 1
	hood |= w.Grid[Lx][Hy] & 1
	hood <<= 1
	hood |= w.Grid[Lx][y] & 1
	return hood
}
