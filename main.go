package main

import (
    "github.com/montanaflynn/stats"
    "gogengol/rule"
    "gogengol/world"
    "fmt"
    "math/rand"
    "time"
//    "strings"
)

func apply(r *rule.Rule, w, buf *world.World) float64 {
    livecells := 0.0
    for i := range w.Grid {
        for j := range w.Grid[i] {
            if w.Grid[i][j] == 0 {
                buf.Grid[i][j] = r.Dead[world.GetHood(w,i,j)]
            } else {
                buf.Grid[i][j] = r.Live[world.GetHood(w,i,j)]
                livecells++
            }
        }
    }
    return 100.0*livecells/float64(w.W*w.H)
}

func main() {
    densities := make([]float64, 1000);
    rand.Seed(time.Now().UnixNano())
    w := world.NewPopPatch(100,100, 20, 0.5);
    wbuf := world.NewEmpty(100,100);
    r := rule.NewRandom(.5, .5);
    fmt.Println(r)
    for i:=0;i<1000;i++ {
        densities[i] = apply(r, w, wbuf)
        fmt.Println(densities[i])
        w, wbuf = wbuf, w
    }
    fmt.Println(stats.StdDevP(densities))
    fmt.Println(stats.Variance(densities))
}
