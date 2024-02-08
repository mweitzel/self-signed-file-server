# self signed file server

## usage

Sometimes you just need a https+basicauth file server.

### server
- run `./make-cert`, provide `.` for all fields except `common name` - common name can be localhost, some ip, or the domain you are requesting via.
- build with `./build`
- export required env vars and run with `./serve`

### client

- provide the `server.crt` to the client. Request with
- `curl --cacert server.crt https://uname:pass@yourcommonname/file/path`
