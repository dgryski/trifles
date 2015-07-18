import lz4
import siphashc
import sys
import struct

def unpack(stream):

    while True:
        data = stream.read(4)
        if len(data) != 4:
            break
        (pktSize,) = struct.unpack('<I', data)
        data = stream.read(8+2)
        if len(data) != 8+2:
            sys.stderr.write('short read')
            break
        (checksum, lenfname) = struct.unpack('<QH', data)
        fname = stream.read(lenfname)
        if len(fname) != lenfname:
            sys.stderr.write('short read')
            break
        data = stream.read(4)
        if len(data) != 4:
            sys.stderr.write('short read')
            break
        (fsize,) = struct.unpack('<I', data)
        compressedSize = pktSize - 4 - 8 - 2 - 4 - lenfname
        data = stream.read(compressedSize)
        if len(data) != compressedSize:
            sys.stderr.write('short read')
            break

        data = lz4.uncompress(data)
        got = siphashc.siphash('\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0', data)
        if got == checksum:
            sys.stderr.write('%s: %d -> %d\n' % (fname, compressedSize, fsize))
        else:
            sys.stderr.write('%s: checksum fail: got %d, want %d\n' % (fname, got, checksum))

unpack(sys.stdin)
