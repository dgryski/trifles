package main

import (
    "errors"
    "strconv"
)

func parseLine(data []byte) ([]byte, []byte, int64, int64, error) {

%% machine scanner;
%% write data;

	cs, p, pe, eof := 0, 0, len(data), len(data)

        var n int

        var t0 []byte
        var m []byte
        var e1 int64
        var e2 int64

	var parsed bool

	// 2015/09/15 16:49:17 fetch: served "metric.name.with.lots.__of_dots__.and.underscores" from 1442324820 to 1442328420

	%%{
            datetime = digit{4} '/' digit{2} '/' digit{2} ' ' digit{2} ':' digit{2} ':' digit{2} ;
            epoch = digit{10} ;
            quoted = '"' [^"]+ '"' ;

            action start { n = fpc }
            action save_datetime { t0 = data[n:fpc] }
            action save_metric { m = data[n:fpc] }
            action save_epoch1 { e1, _ = strconv.ParseInt(string(data[n:fpc]), 10, 64) }
            action save_epoch2 { e2, _ = strconv.ParseInt(string(data[n:fpc]), 10, 64) }

	    main := (datetime >start %save_datetime) ' fetch: served '  ( quoted >start %save_metric ) ' from ' ( epoch >start %save_epoch1 ) ' to ' ( epoch >start %save_epoch2 ) @{ parsed = true };

	    write init;
	    write exec;
	}%%

	if !parsed {
	    return nil, nil, 0, 0, errors.New("parse error")
	}

	return t0, m, e1, e2, nil
}
