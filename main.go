// Tool for automated search of non-isotropic GOL-like CAs for interesting 
// specimens. 
package main

import (
    "flag"
    "fmt"
    "github.com/montanaflynn/stats"
    "gogengol/rule"
    "gogengol/world"
    "math/rand"
    "os"
    "time"
)

type RunParams struct {
    nrules uint
    nruns uint
    ngens uint
    rdd float64
    rdl float64
    neatness_threshold float64
    fh *os.File
}

func check(err error) {
    if err != nil {
        panic(err)
    }
}

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

func run(p RunParams) int {
    neat := 0.0
    nfound := 0

    rdd := p.rdd
    rdl := p.rdl
    //rdd := rand.Float64()/2 + 0.25
    //rdl := rand.Float64()/2 + 0.25

    dens := make([]float64, p.ngens) // Storage for pop densities over a run. 
    means := make([]float64, p.nruns)
    stdvs := make([]float64, p.nruns)
    corrs := make([]float64, p.nruns)
    meansb := make([]float64, p.nruns)
    stdvsb := make([]float64, p.nruns)
    corrsb := make([]float64, p.nruns)

    line := make([]float64, p.ngens)
    for i:=0;i < int(p.ngens); i++ {
        line[i] = float64(i)
    }

    for i:=0;i < int(p.nrules);i++ { // rule loop
        if p.rdd <= 0.0 {
            rdd = rand.Float64()
        }
        if p.rdl <= 0.0 {
            rdl = rand.Float64()
        }

        r := rule.NewRandom(rdd, rdl)
        //r := TheOriginal()
        neat = 0.0

        for j:=0;j < int(p.nruns);j++ { // run loop
            w := world.NewPopPatch(100,100, 20, rand.Float64()*.8+0.1);
            wbuf := world.NewEmpty(100,100);
            for i:=0;i < int(p.ngens);i++ { // generation loop
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
        if neat > p.neatness_threshold {
            nfound++
            fmt.Println("Mean A: ", meana, " Mean B: ", meanb)
            fmt.Println("Stdev ALL min: ", stdvlo, " max: ", stdvhi)
            fmt.Println("Stdev 1/2 min: ", stdvloB, " max: ", stdvhiB)
            fmt.Println("Corr ALL min: ", corrslo, " max: ", corrshi)
            fmt.Println("Corr 1/2 min: ", corrsloB, " max: ", corrshiB)
            fmt.Println("Neat: ", neat)
            r.Comment = fmt.Sprintf("rdd:%f rdl:%f neat:%f", rdd, rdl, neat)
            if nfound > 1 {
                p.fh.WriteString(",")
            }
            p.fh.WriteString("\n" + fmt.Sprint(r))
        }
    }
    return nfound
}

func main() {
    rand.Seed(time.Now().UnixNano())

    var nrules uint = 100
    flag.UintVar(&nrules, "rules", nrules, "The number of rulesets to try.")
    var nruns uint = 20
    flag.UintVar(&nruns, "runs", nruns, "The number of runs for each ruleset.")
    var ngens uint = 200
    flag.UintVar(&ngens, "gens", ngens, "The number of generations per run.")

    outfn := "output.json"
	flag.StringVar(&outfn, "o", outfn, "Set the output file.")

    rdd := 0.21875  // Original GOL densities. 
    rdl := 0.32815
    flag.Float64Var(&rdd, "rdd", rdd, "Rule density for dead cells.")
    flag.Float64Var(&rdl, "rdl", rdl, "Rule density for live cells.")

    neat := 0.74
    flag.Float64Var(&neat, "neat", neat, "Neatness threshold.")

    flag.Parse()

    var params RunParams

    if _, err := os.Stat(outfn); err == nil {
        fmt.Println("File already exists: " + outfn)
        os.Exit(1)
    }

    fh, err := os.Create(outfn)
    check(err)

    _, err = fh.WriteString("[\n")
    check(err)

    params.nrules = nrules
    params.nruns = nruns
    params.ngens = ngens
    params.rdd = rdd
    params.rdl = rdl
    params.neatness_threshold = neat
    params.fh = fh

    nfound := run(params)

    _, err = fh.WriteString("\n]\n")
    check(err)

    fh.Sync()
    fh.Close()
    fmt.Println("Found", nfound, "neat rulesets.")
    fmt.Println("Output written to " + outfn + ".")

}
