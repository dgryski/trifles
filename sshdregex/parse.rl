package main

func matchSSHD(data []byte) bool {

%% machine scanner;
%% write data;

	cs, p, pe, eof := 0, 0, len(data), len(data)

        _ = eof

	%%{
	    main := any* 'sshd[' digit{5} ']:' space* 'Failed' @{ return true } ;

	    write init;
	    write exec;
	}%%

        return false
}
