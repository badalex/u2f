gou2f
======

Brain dead simple example.

Note you must generate some ssl certs:
go run /usr/lib/go/src/pkg/crypto/tls/generate_cert.go  --host gou2f.com

and add the below to your hosts file:
127.0.0.1 gou2f.com

Currently it does not seem to work if you point it at localhost :(
