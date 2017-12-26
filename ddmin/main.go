package main

import (
	"bytes"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	ddmin "github.com/dgryski/go-ddmin"
)

func run(prog string, data []byte) ([]byte, error) {
	cmd := exec.Command(prog)
	cmd.Stdin = bytes.NewReader(data)
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

func main() {
	good := flag.String("good", "", "known good command")
	bad := flag.String("bad", "", "known bad command")
	flag.Parse()

	input, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalf("unable to read stdin: %v", err)
	}

	// cache of responses
	m := make(map[string]ddmin.Result)

	os.Stdout.Write(ddmin.Minimize(input, func(d []byte) ddmin.Result {
		if r, ok := m[string(d)]; ok {
			return r
		}

		g, err := run(*good, d)
		if err != nil {
			m[string(d)] = ddmin.Unresolved
			return ddmin.Unresolved
		}

		b, err := run(*bad, d)
		if err != nil {
			m[string(d)] = ddmin.Unresolved
			return ddmin.Unresolved
		}

		if bytes.Equal(b, g) {
			m[string(d)] = ddmin.Pass
			return ddmin.Pass
		}

		m[string(d)] = ddmin.Fail
		return ddmin.Fail
	}))
}
