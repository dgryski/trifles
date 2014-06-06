# udp tests
import socket
import time

UDP_IP = "127.0.0.1"
UDP_PORT = 12233

sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)

message = "hello, world"

for i in xrange(100):
    for j in xrange(5000): sock.sendto(message, (UDP_IP, UDP_PORT))
    time.sleep(0.1)

