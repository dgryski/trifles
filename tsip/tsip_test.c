

#include <stdio.h>
#include <stdint.h>

#include "tsip.h"

int main() {

    const unsigned char seed[] = {0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f};

    size_t i;
    unsigned char buf[64];
    uint64_t h;

    for (i=0;i< 64;i++) {
        buf[i]=i;
        h = tsip(seed, (const unsigned char *)buf, i);
        printf("%016lx\n", h);
    }

    return 0;
}
