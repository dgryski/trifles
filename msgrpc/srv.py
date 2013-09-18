import msgpackrpc

class EchoHandler(object):

    def Echo123(self, msg):
        (msg1, msg2, msg3) = msg
        return ("1:%s 2:%s 3:%s" % (msg1, msg2, msg3))

    def EchoStruct(self, msg):
        return ("msg:%s" % msg)

if __name__ == '__main__':
        addr = msgpackrpc.Address('localhost', 9000)
        server = msgpackrpc.Server(EchoHandler())
        server.listen(addr)
        server.start()
