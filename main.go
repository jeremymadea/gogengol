package main

import (
    "github.com/montanaflynn/stats"
    "gogengol/rule"
    "gogengol/world"
    "fmt"
    "math/rand"
    "time"
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

func TheOriginal() *rule.Rule {
    r := rule.NewFromStrings(
        "0000000100010110000101100110100000010110011010000110100010000000000101100110100001101000100000000110100010000000100000000000000000010110011010000110100010000000011010001000000010000000000000000110100010000000100000000000000010000000000000000000000000000000",
        "0001011101111110011111101110100001111110111010001110100010000000011111101110100011101000100000001110100010000000100000000000000001111110111010001110100010000000111010001000000010000000000000001110100010000000100000000000000010000000000000000000000000000000")
    return r
}

func main() {
    neat := 0.0
    nrules := 100
    nruns  := 20
    ngens  := 200
    dens := make([]float64, ngens) // Storage for pop densities over a run. 

    means := make([]float64, nruns)
    stdvs := make([]float64, nruns)
    corrs := make([]float64, nruns)
    meansb := make([]float64, nruns)
    stdvsb := make([]float64, nruns)
    corrsb := make([]float64, nruns)

    line := make([]float64, ngens)
    for i:=0;i < ngens; i++ {
        line[i] = float64(i)
    }

    rand.Seed(time.Now().UnixNano())
    for i:=0;i < nrules;i++ { // rule loop
        //rdd := rand.Float64()/2 + 0.25
        //rdl := rand.Float64()/2 + 0.25
        rdd := 0.21875  // Original GOL densities. 
        rdl := 0.32815
        r := rule.NewRandom(rdd, rdl)
        //r := TheOriginal()
        neat = 0.0
//        fmt.Println()
//        fmt.Println("----------------------------------------")
//        fmt.Println(" Rule Density: rdd: ", rdd, " rdl: ", rdl);

        for j:=0;j < nruns;j++ { // run loop
            w := world.NewPopPatch(100,100, 20, rand.Float64()*.8+0.1);
            wbuf := world.NewEmpty(100,100);
            for i:=0;i < ngens;i++ { // generation loop
                dens[i] = apply(r, w, wbuf)
                w, wbuf = wbuf, w
            }
            mean, _ := stats.Mean(dens)
            std, _ := stats.StdDevP(dens)
            pearson, _ := stats.Pearson(line, dens)
            means[j] = mean
            stdvs[j] = std
            corrs[j] = pearson
            mean, _ = stats.Mean(dens[len(dens)/2:len(dens)-1])
            std, _ = stats.StdDevP(dens[len(dens)/2:len(dens)-1])
            pearson, _ = stats.Pearson(
                line[len(line)/2:len(line)-1], dens[len(dens)/2:len(dens)-1])
            meansb[j] = mean
            stdvsb[j] = std
            corrsb[j] = pearson
        }
        meana, _ := stats.Mean(means)
        meanb, _ := stats.Mean(meansb)
        stdvlo, _ := stats.Min(stdvs)
        stdvhi, _ := stats.Max(stdvs)
        stdvloB, _ := stats.Min(stdvsb)
        stdvhiB, _ := stats.Max(stdvsb)
        corrslo, _ := stats.Min(corrs)
        corrshi, _ := stats.Max(corrs)
        corrsloB, _ := stats.Min(corrsb)
        corrshiB, _ := stats.Max(corrsb)

        neat = (corrshiB - corrsloB) / 2
        if neat > 0.74 {
            fmt.Println("Mean A: ", meana, " Mean B: ", meanb)
            fmt.Println("Stdev ALL min: ", stdvlo, " max: ", stdvhi)
            fmt.Println("Stdev 1/2 min: ", stdvloB, " max: ", stdvhiB)
            fmt.Println("Corr ALL min: ", corrslo, " max: ", corrshi)
            fmt.Println("Corr 1/2 min: ", corrsloB, " max: ", corrshiB)
            fmt.Println("Neat: ", neat)
            r.Comment = fmt.Sprintf("rdd:%f rdl:%f neat:%f", rdd, rdl, neat)
            fmt.Println(r)
        }

    }
}
