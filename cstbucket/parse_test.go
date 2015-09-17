package main

import (
	"bytes"
	"testing"
)

func TestParseLine(t *testing.T) {

	q := []byte(`2015/09/15 16:49:17 fetch: served "metric.name.with.lots.__of_dots__.and.underscores" from 1442324820 to 1442328420`)

	t0, m, e1, e2, err := parseLine(q)

	var want = struct {
		t0 []byte
		m  []byte
		e1 int64
		e2 int64
	}{
		[]byte("2015/09/15 16:49:17"),
		[]byte(`"metric.name.with.lots.__of_dots__.and.underscores"`),
		1442324820,
		1442328420,
	}

	if !bytes.Equal(t0, want.t0) || !bytes.Equal(m, want.m) || e1 != want.e1 || e2 != want.e2 || err != nil {
		t.Errorf("parseLine(q)=(%v,%v,%v,%v,%v), want (%v,%v,%v,%v,%v)", string(t0), string(m), e1, e2, err, string(want.t0), string(want.m), want.e1, want.e2, nil)
	}
}
