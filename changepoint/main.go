package main

import (
	"bufio"
	"flag"
	"fmt"
	"html/template"
	"io"
	"log"
	"math"
	"os"
	"sort"
	"strconv"

	"github.com/dgryski/go-change"
)

func main() {
	window := flag.Int("w", 120, "window size")
	minSample := flag.Int("ms", 30, "min sample size")
	blockSize := flag.Int("bs", 10, "block size")
	compressPoints := flag.Int("cp", 10, "compress points for graph display")
	fname := flag.String("f", "", "file name")

	flag.Parse()

	var f io.Reader

	if *fname == "" {
		log.Println("reading from stdin")
		f = os.Stdin
	} else {
		var err error
		f, err = os.Open(*fname)
		if err != nil {
			fmt.Println("open failed:", err)
			return
		}
	}

	scanner := bufio.NewScanner(f)

	data := make([]float64, *window)
	buffer := make([]float64, *blockSize)
	var bufidx int
	var items int

	var detector = change.Detector{
		MinSampleSize: *minSample,
		TConf:         change.Conf99p5,
	}

	type graphPoints [2]float64

	var graphData []graphPoints
	var last []float64

	var changePoints []int

	for scanner.Scan() {
		item, err := strconv.ParseFloat(scanner.Text(), 64)
		if err != nil {
			fmt.Printf("error parsing <%s>: %s\n", scanner.Text(), err)
			continue
		}

		last = append(last, item)
		if items > 0 && items%*compressPoints == 0 {
			sort.Float64s(last)
			median := last[*compressPoints/2]
			last = last[:0]

			graphData = append(graphData, [2]float64{float64(items), median})
		}

		buffer[bufidx] = item

		bufidx++
		items++
		if bufidx == *blockSize {
			copy(data[0:], data[*blockSize:])
			copy(data[*window-*blockSize:], buffer)
			bufidx = 0

			if items >= *window {
				r := detector.Check(data)
				//				s := onlinestats.NewSlice(data)
				diff := math.Abs(r.Difference / r.Before.Mean())
				if r.Difference != 0 && diff > 0.06 {
					log.Printf("difference found at offset=%d: %f %v\n", items-*window+r.Index, diff, r)
					changePoints = append(changePoints, items-*window+r.Index)
				}
			}
		}
	}

	reportTmpl.Execute(os.Stdout, struct {
		GraphData    []graphPoints
		ChangePoints []int
	}{
		graphData,
		changePoints,
	})
}

var reportTmpl = template.Must(template.New("report").Parse(`
<html>
<script src="//cdnjs.cloudflare.com/ajax/libs/jquery/2.0.3/jquery.min.js"></script>
<script src="//cdnjs.cloudflare.com/ajax/libs/flot/0.8.2/jquery.flot.min.js"></script>
<script src="//cdnjs.cloudflare.com/ajax/libs/flot/0.8.2/jquery.flot.time.min.js"></script>

<script type="text/javascript">

    var data = {{ .GraphData }};

    $(document).ready(function() {
        $.plot($("#placeholder"), [data], {
             yaxis: { min: 2000 },
             grid: {
                markings: [
                  {{ range .ChangePoints }}{ color: '#000', lineWidth: 1, xaxis: { from: {{ . }}, to: {{ . }} } },
                  {{ end }}
                ]
              }
           })
        })

</script>

<body>

<div id="placeholder" style="width:1200px; height:400px"></div>

</body>
</html>
`))
