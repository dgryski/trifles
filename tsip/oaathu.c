
#include <stdint.h>
#include <stdio.h>

uint32_t oaathu(const unsigned char *seed, const unsigned char *str, size_t len) {
  uint32_t hash = *((uint32_t *)seed) + len;

  switch (len) {
  case 16:
    hash += (hash << 10);
    hash ^= (hash >> 6);
    hash += str[15];
  case 15:
    hash += (hash << 10);
    hash ^= (hash >> 6);
    hash += str[14];
  case 14:
    hash += (hash << 10);
    hash ^= (hash >> 6);
    hash += str[13];
  case 13:
    hash += (hash << 10);
    hash ^= (hash >> 6);
    hash += str[12];
  case 12:
    hash += (hash << 10);
    hash ^= (hash >> 6);
    hash += str[11];
  case 11:
    hash += (hash << 10);
    hash ^= (hash >> 6);
    hash += str[10];
  case 10:
    hash += (hash << 10);
    hash ^= (hash >> 6);
    hash += str[9];
  case 9:
    hash += (hash << 10);
    hash ^= (hash >> 6);
    hash += str[8];
  case 8:
    hash += (hash << 10);
    hash ^= (hash >> 6);
    hash += str[7];
  case 7:
    hash += (hash << 10);
    hash ^= (hash >> 6);
    hash += str[6];
  case 6:
    hash += (hash << 10);
    hash ^= (hash >> 6);
    hash += str[5];
  case 5:
    hash += (hash << 10);
    hash ^= (hash >> 6);
    hash += str[4];
  case 4:
    hash += (hash << 10);
    hash ^= (hash >> 6);
    hash += str[3];
  case 3:
    hash += (hash << 10);
    hash ^= (hash >> 6);
    hash += str[2];
  case 2:
    hash += (hash << 10);
    hash ^= (hash >> 6);
    hash += str[1];
  case 1:
    hash += (hash << 10);
    hash ^= (hash >> 6);
    hash += str[0];
  case 0:
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
  }
  return (hash + (hash << 15));
}
