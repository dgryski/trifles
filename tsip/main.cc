
#include "benchmark/benchmark_api.h"

unsigned char buffer[8192];

#include "tsip.h"

uint32_t one_at_a_time_hard(const unsigned char *seed, const unsigned char *str, size_t len);
uint64_t oaathu(const unsigned char *seed, const unsigned char *message, size_t len);


void Hash_oaath(benchmark::State& state) {
    unsigned char seed[16];
    while (state.KeepRunning()) {
        benchmark::DoNotOptimize(
                one_at_a_time_hard(seed, buffer, state.range(0))
        );
    }
}

void Hash_oaathu(benchmark::State& state) {
    unsigned char seed[16];
    while (state.KeepRunning()) {
        benchmark::DoNotOptimize(
                oaathu(seed, buffer, state.range(0))
        );
    }
}

void Hash_tsip(benchmark::State& state) {
    unsigned char seed[16];
    while (state.KeepRunning()) {
        benchmark::DoNotOptimize(
                tsip(seed, buffer, state.range(0))
        );
    }
}


static void CustomArguments(benchmark::internal::Benchmark *b) {
    uint64_t sizes[] = {1,2,3,4,5,6,7,8,9,10,15,20,30,40,50,100, 0};
    for (int i=0; sizes[i] != 0; i++) {
        b->Arg(sizes[i]);
    }
}

static void Custom0To16(benchmark::internal::Benchmark *b) {
    for (int i=0; i <= 16; i++) {
        b->Arg(i);
    }
}

BENCHMARK(Hash_tsip)->Apply(CustomArguments);
BENCHMARK(Hash_oaath)->Apply(CustomArguments);
BENCHMARK(Hash_oaathu)->Apply(Custom0To16);

BENCHMARK_MAIN()
