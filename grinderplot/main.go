package main

import (
	"bufio"
	"bytes"
	"flag"
	"html/template"
	"os"
	"sort"
	"strconv"
)

func main() {

	flag.Parse()

	var data = make(map[int][]int)
	for _, fname := range flag.Args() {
		f, _ := os.Open(fname)
		b := bufio.NewReader(f)

		scanner := bufio.NewScanner(b)

		// skip first line with the header
		scanner.Scan()

		for scanner.Scan() {
			line := scanner.Bytes()
			parts := bytes.Split(line, []byte{','})
			epochms, reqms := parts[3], parts[4]
			epoch, _ := strconv.Atoi(string(bytes.TrimSpace(epochms)))
			epoch /= 1000
			latency, _ := strconv.Atoi(string(bytes.TrimSpace(reqms)))
			data[epoch] = append(data[epoch], latency)
		}
	}

	var epochs []int
	for k, v := range data {
		epochs = append(epochs, k)
		sort.Ints(v)
	}

	sort.Ints(epochs)

	var graphdata [][]int
	for _, epoch := range epochs {
		l := len(data[epoch])
		graphdata = append(graphdata, []int{epoch, data[epoch][(50*l)/100], data[epoch][(75*l)/100], data[epoch][(95*l)/100], data[epoch][(99*l)/100]})
	}

	reportTmpl.Execute(os.Stdout, graphdata)
}

var reportTmpl = template.Must(template.New("report").Parse(`
<html>
<script src="//cdnjs.cloudflare.com/ajax/libs/jquery/2.0.3/jquery.min.js"></script>
<script src="//cdnjs.cloudflare.com/ajax/libs/flot/0.8.2/jquery.flot.min.js"></script>
<script src="//cdnjs.cloudflare.com/ajax/libs/flot/0.8.2/jquery.flot.time.min.js"></script>

<script type="text/javascript">

    var data = {{ . }};

    var p50 = data.map(function(e) { return [e[0], e[1]] });
    var p75 = data.map(function(e) { return [e[0], e[2]] });
    var p95 = data.map(function(e) { return [e[0], e[3]] });
    var p99 = data.map(function(e) { return [e[0], e[4]] });

    $(document).ready(function() {
        $.plot($("#placeholder"), [p50, p75, p95, p99 ], {
                lines: { show: false },
                points: { show: true, fillColor: false },
                xaxis: { mode: "time" }
            })
    })

</script>

<body>

<div id="placeholder" style="width:1200px; height:800px"></div>

</body>
</html>
`))
