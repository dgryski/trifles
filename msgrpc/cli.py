import msgpackrpc

address = msgpackrpc.Address('localhost', 9000)
client = msgpackrpc.Client(address, unpack_encoding='utf-8')

# go server
#print client.call("Echo.Echo123", ["foo1", "foo2", "foo3"])
#print client.call("Echo.EchoStruct", {"A" :"a string", "B":"b string", "C":"c string"})
# python server
print client.call("Echo123", ["foo1", "foo2", "foo3"])
print client.call("EchoStruct", {"A" :"a string", "B":"b string", "C":"c string"})
