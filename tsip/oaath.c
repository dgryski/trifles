
#include <stdint.h>
#include <stdio.h>

uint32_t one_at_a_time_hard(const unsigned char *seed, const unsigned char *str,
                            size_t len) {
  const unsigned char *const end = (const unsigned char *)str + len;
  uint32_t hash = *((uint32_t *)seed) + len;

  while (str < end) {
    hash += (hash << 10);
    hash ^= (hash >> 6);
    hash += *str++;
  }

  hash += (hash << 10);
  hash ^= (hash >> 6);
  hash += seed[4];

  hash += (hash << 10);
  hash ^= (hash >> 6);
  hash += seed[5];

  hash += (hash << 10);
  hash ^= (hash >> 6);
  hash += seed[6];

  hash += (hash << 10);
  hash ^= (hash >> 6);
  hash += seed[7];

  hash += (hash << 10);
  hash ^= (hash >> 6);

  hash += (hash << 3);
  hash ^= (hash >> 11);
  return (hash + (hash << 15));
}
