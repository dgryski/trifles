
CFLAGS=-O2

.c.o:
	g++ -c $(CFLAGS) $< -o $@

bench: main.cc tsip.o oaath.o oaathu.o
	g++ $(CFLAGS) $^ -o $@ -lbenchmark -lpthread

tsiptest: tsip.o tsip_test.o
	g++ $(CFLAGS) $^ -o $@
