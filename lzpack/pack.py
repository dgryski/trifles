import lz4
import siphashc
import sys
import struct

total = 0
compressed = 0


def pack(fname):

    global total, compressed


    f = open(fname)
    data = f.read()
    f.close()
    checksum = siphashc.siphash('\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0', data)
    fsize = len(data)
    data = lz4.compress(data)

    # size of packet(4), checksum(8), fnamelen(2)+fname, uncompressed size(4), compressed data
    l = len(data)
    pktlen = 4 + 8 + 2 + len(fname) + 4 + l

    total += fsize
    compressed += len(data)

    sys.stderr.write("%s: %d -> %d\n"  %(fname, fsize, len(data)))
    sys.stdout.write( struct.pack('<IQH%dsI%ds' % (len(fname), l), pktlen, checksum, len(fname), fname,fsize,data))
    sys.stdout.flush()

for fn in sys.stdin:
    pack(fn.rstrip())

if total == 0:
    total = 1

sys.stderr.write("compressed %d -> %d (%f %%)\n" % (total, compressed, 100-(100.0*+compressed)/total))
