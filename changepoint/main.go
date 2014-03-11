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

type streamDetector struct {
	windowSize int
	blockSize  int

	data []float64

	buffer []float64
	bufidx int

	items int

	detector *change.Detector
}

func New(windowSize int, minSample int, blockSize int) *streamDetector {
	return &streamDetector{
		windowSize: windowSize,
		blockSize:  blockSize,
		data:       make([]float64, windowSize),
		buffer:     make([]float64, blockSize),
		detector: &change.Detector{
			MinSampleSize: minSample,
			TConf:         change.Conf99p5,
		},
	}
}

func (s *streamDetector) Push(item float64) *change.ChangePoint {
	s.buffer[s.bufidx] = item

	s.bufidx++
	s.items++

	if s.bufidx < s.blockSize {
		return nil
	}

	copy(s.data[0:], s.data[s.blockSize:])
	copy(s.data[s.windowSize-s.blockSize:], s.buffer)
	s.bufidx = 0

	return s.detector.Check(s.data)
}

func main() {
	windowSize := flag.Int("w", 120, "window size")
	minSample := flag.Int("ms", 30, "min sample size")
	blockSize := flag.Int("bs", 10, "block size")
	compressPoints := flag.Int("cp", 10, "compress points for graph display")
	fname := flag.String("f", "", "file name")
	ymin := flag.Int("ymin", 0, "minimum y value for graph")

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

	s := New(*windowSize, *minSample, *blockSize)

	type graphPoints [2]float64
	var graphData []graphPoints
	var last []float64

	var changePoints []int

	var items int

	for scanner.Scan() {
		item, err := strconv.ParseFloat(scanner.Text(), 64)
		if err != nil {
			fmt.Printf("error parsing <%s>: %s\n", scanner.Text(), err)
			continue
		}

		last = append(last, item)
		items++
		if items > 0 && items%*compressPoints == 0 {
			sort.Float64s(last)
			median := last[*compressPoints/2]
			last = last[:0]

			graphData = append(graphData, [2]float64{float64(items), median})
		}

		r := s.Push(item)
		if r != nil {
			diff := math.Abs(r.Difference / r.Before.Mean())
			if r.Difference != 0 && diff > 0.06 {
				log.Printf("difference found at offset=%d: %f %v\n", s.items-s.windowSize+r.Index, diff, r)
				changePoints = append(changePoints, items-s.windowSize+r.Index)
			}
		}
	}

	reportTmpl.Execute(os.Stdout, struct {
		YMin         int
		GraphData    []graphPoints
		ChangePoints []int
	}{
		*ymin,
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
             yaxis: { min: {{ .YMin }} },
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
