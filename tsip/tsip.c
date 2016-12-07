
#include <stdint.h>
#include <stdio.h>

static uint64_t inline U8TO64_LE(const unsigned char *p) {
  return *(const uint64_t *)p;
}

static uint64_t rotl64(uint64_t x, int k) {
  return (((x) << (k)) | ((x) >> (64 - k)));
}

uint64_t tsip(unsigned char *seed, const unsigned char *m, size_t len) {

  uint64_t v0, v1;
  uint64_t mi, k0, k1;
  uint64_t last7;
  size_t i, blocks;

  k0 = U8TO64_LE(seed);
  k1 = U8TO64_LE(seed + 8);

  v0 = k0 ^ 0x736f6d6570736575ull;
  v1 = k1 ^ 0x646f72616e646f6dull;

  last7 = (uint64_t)(len & 0xff) << 56;

#define sipcompress()                                                          \
  do {                                                                         \
    v0 += v1;                                                                  \
    v1 = rotl64(v1, 13) ^ v0;                                                  \
    v0 = rotl64(v0, 32) + v1;                                                  \
    v1 = rotl64(v1, 17) ^ v0;                                                  \
  } while (0)

  for (i = 0, blocks = (len & ~7); i < blocks; i += 8) {
    mi = U8TO64_LE(m + i);
    v1 ^= mi;
    sipcompress();
    v0 ^= mi;
  }

  switch (len - blocks) {
  case 7:
    last7 |= (uint64_t)m[i + 6] << 48;
  case 6:
    last7 |= (uint64_t)m[i + 5] << 40;
  case 5:
    last7 |= (uint64_t)m[i + 4] << 32;
  case 4:
    last7 |= (uint64_t)m[i + 3] << 24;
  case 3:
    last7 |= (uint64_t)m[i + 2] << 16;
  case 2:
    last7 |= (uint64_t)m[i + 1] << 8;
  case 1:
    last7 |= (uint64_t)m[i + 0];
  case 0:
  default:;
  };

  v1 ^= last7;
  sipcompress();
  v0 ^= last7;

  // finalization
  v1 ^= 0xff;
  sipcompress();
  sipcompress();

  return v0 ^ v1;
}
