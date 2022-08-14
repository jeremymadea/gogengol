package main

import (
    "github.com/montanaflynn/stats"
    "fmt"
    "math/rand"
    "time"
    "strings"
)

type Rule struct {
    dead [256]byte
    live [256]byte
}

func (r Rule) String() string {
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

type World struct {
    w int
    h int
    grid [][]byte
}

func NewRandomWorld(w, h int) World {
    wld := World{w, h, nil}
    wld.grid = make([][]byte, h)
    for i := range wld.grid {
        wld.grid[i] = make([]byte, w)
        for j := range wld.grid[i] {
            wld.grid[i][j] = byte(rand.Int() & 1);
        }
    }
    return wld
}

func NewWorld(w, h int) World {
    wld := World{w, h, nil}
    wld.grid = make([][]byte, h)
    for i := range wld.grid {
        wld.grid[i] = make([]byte, w)
    }
    return wld
}

func InitRandomPatch(w World, z int, pd float64) {
    for i := range w.grid {
        for j := range w.grid[i] {
            w.grid[i][j] = 0;
            if i > w.w-z/2 && i < w.w+z/2 {
                if j > w.h-z/2 && i < w.h+z/2 {
                    if rand.Float64() < pd {
                        w.grid[i][j] = 1
                    }
                }
            }
        }
    }
}

func ClearWorld(w World) {
    for i := range w.grid {
        for j := range w.grid[i] {
            w.grid[i][j] = 0;
        }
    }
}

func NewRandomRule() Rule {
    r := Rule{}
    for i := range r.dead { r.dead[i] = byte(rand.Int() & 1) }
    for i := range r.live { r.live[i] = byte(rand.Int() & 1) }
    return r
}

// Returns a byte representing one of the 256 possible Moore neighborhood
// states for the cell located at (x,y) in the passed world. 
func gethood(w World, x, y int) byte {
    hood := byte(0)
    Ly := y-1
    Lx := x-1
    Hy := y+1
    Hx := x+1
    if Hx >= w.w { Hx = 0 }
    if Hy >= w.h { Hy = 0 }
    if Lx <= 0   { Lx = w.w-1 }
    if Ly <= 0   { Ly = w.h-1 }
    hood |= w.grid[Lx][Ly] & 1
    hood <<= 1
    hood |= w.grid[x][Ly]  & 1
    hood <<= 1
    hood |= w.grid[Hx][Ly] & 1
    hood <<= 1
    hood |= w.grid[Hx][y]  & 1
    hood <<= 1
    hood |= w.grid[Hx][Hy] & 1
    hood <<= 1
    hood |= w.grid[x][Hy]  & 1
    hood <<= 1
    hood |= w.grid[Lx][Hy] & 1
    hood <<= 1
    hood |= w.grid[Lx][y]  & 1
    return hood
}

func apply(r Rule, w, buf World) float64 {
    livecells := 0.0
    for i := range w.grid {
        for j := range w.grid[i] {
            if w.grid[i][j] == 0 {
                buf.grid[i][j] = r.dead[gethood(w,i,j)]
            } else {
                buf.grid[i][j] = r.live[gethood(w,i,j)]
                livecells++
            }
        }
    }
    return 100.0*livecells/float64(w.w*w.h)
}

func main() {
    densities := make([]float64, 1000);
    rand.Seed(time.Now().UnixNano())
    w := NewRandomWorld(100,100);
    wbuf := NewWorld(100,100);
    InitRandomPatch(w, 20, 0.3);
    r := NewRandomRule();
    fmt.Println(r)
    for i:=0;i<1000;i++ {
        densities[i] = apply(r, w, wbuf)
        fmt.Println(densities[i])
        w, wbuf = wbuf, w
/*
        fmt.Println(w.grid[0][0], w.grid[0][1], w.grid[0][2])
        fmt.Println(w.grid[1][0], w.grid[1][1], w.grid[1][2])
        fmt.Println(w.grid[2][0], w.grid[2][1], w.grid[2][2]) 
        fmt.Println()  */
    }
    fmt.Println(stats.StdDevP(densities))
    fmt.Println(stats.Variance(densities))
}
