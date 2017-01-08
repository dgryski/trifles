
CFLAGS=-O2

.c.o:
	g++ -c $(CFLAGS) $< -o $@

bench: main.cc tsip.o oaath.o oaathu.o
	g++ $(CFLAGS) $^ -o $@ -lbenchmark -lpthread
