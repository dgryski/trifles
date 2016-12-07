

.c.o:
	g++ -c -O3 $< -o $@

bench: main.cc tsip.o oaath.o oaathu.o
	g++ -O3 $^ -o $@ -lbenchmark -lpthread
