# grpc-client-server

Implementation of server-side streaming via GRPC 

The client on run reads URL from flags, then makes GRPC call to the server, which makes HTTP GET request to passed URL and returns the result to the client via GRPC stream with chunks in 1024 bytes.

### **Usage:**

1. `git clone https://github.com/dmi-kov/grpc-client-server.git`
2. `cd grpc-client-server`
3. `make build`
4. `./grpc-server`
5. in another terminal `./grpc-client -url=http://google.com`

Run `make help` to see all commands.