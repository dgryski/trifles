#ifdef __cplusplus
extern "C" {
#endif

#include <stdint.h>
#include <stdio.h>


#define TSIP_STATIC_INLINE static inline
#define STRLEN int

TSIP_STATIC_INLINE
uint64_t U8TO64_LE(const unsigned char *p) {
  return *(const uint64_t *)p;
}

TSIP_STATIC_INLINE
uint64_t rotl64(uint64_t x, int k) {
  return (((x) << (k)) | ((x) >> (64 - k)));
}

TSIP_STATIC_INLINE
void tsip_seed_state(const unsigned char *seed, unsigned char *state_ch) {
  uint64_t *state= (uint64_t *)state_ch;
  state[0] = U8TO64_LE(seed+0) ^ 0x736f6d6570736575ull;
  state[1] = U8TO64_LE(seed+8) ^ 0x646f72616e646f6dull;
}

TSIP_STATIC_INLINE
uint64_t tsip_hash_with_state(const unsigned char *state, const unsigned char *m, size_t len) {

  uint64_t v0, v1;
  uint64_t mi;
  uint64_t last7;

  v0 = U8TO64_LE(state);
  v1 = U8TO64_LE(state + 8);

  last7 = (uint64_t)(len & 0xff) << 56;

#define sipcompress()                                                          \
  do {                                                                         \
    v0 += v1;                                                                  \
    v1 = rotl64(v1, 13) ^ v0;                                                  \
    v0 = rotl64(v0, 35) + v1;                                                  \
    v1 = rotl64(v1, 17) ^ v0;                                                  \
    v0 = rotl64(v0, 21);                                                       \
  } while (0)

  const unsigned char *end = m + (len & ~7);

  while (m < end) {
    mi = U8TO64_LE(m);
    v1 ^= mi;
    sipcompress();
    v0 ^= mi;
    m += 8;
  }

  switch (len & 7) {
  case 7:
    last7 |= (uint64_t)m[6] << 48;
  case 6:
    last7 |= (uint64_t)m[5] << 40;
  case 5:
    last7 |= (uint64_t)m[4] << 32;
  case 4:
    last7 |= (uint64_t)m[3] << 24;
  case 3:
    last7 |= (uint64_t)m[2] << 16;
  case 2:
    last7 |= (uint64_t)m[1] << 8;
  case 1:
    last7 |= (uint64_t)m[0];
  case 0:
  default:;
  };

  v1 ^= last7;
  sipcompress();
  v0 ^= last7;

  // finalization
  v1 ^= 0xff;
  sipcompress();
  v1 = rotl64(v1, 32);
  sipcompress();
  v1 = rotl64(v1, 32);
  sipcompress();
  v1 = rotl64(v1, 32);

  return v0 ^ v1;
}

TSIP_STATIC_INLINE
uint64_t tsip_hash(const unsigned char *seed, const unsigned char *m, size_t len) {
  uint64_t state[2];
  tsip_seed_state(seed,(unsigned char*)state);
  return tsip_hash_with_state((unsigned char*)state, m, len);
}


void tsip_seed_state_smhasher_test(int in_bits, const void *seed, void *state) {
    tsip_seed_state((unsigned char*)seed,(unsigned char *)state);
}

void tsip_hash_with_state_smhasher_test(const void *key, STRLEN len, const void *state, void *out) {
    *((uint64_t *)out)= tsip_hash_with_state((unsigned char*)state, (unsigned char *)key, len);
}

#ifdef __cplusplus
}
#endif

