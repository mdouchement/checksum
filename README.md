# Checksum

Checksum is a simple tool that compute checksums' algorithms on the given file.

## Requirements

- Golang 1.9.x or above

## Installation

```
go get github.com/mdouchement/checksum
```

or in the configured folder

```
dep ensure -vendor-only
go install
```

## Usage

```
$ checksum -h

$ checksum /bin/ls
Checksums:
       crc32: ffd27df2
         md5: d77c1dd5bb8e39c2dd27c96c3fd2263e
        sha1: e332cf8e1a78427f1368a5a0a67946ad1e7c8e28
      sha256: 5abf61c361e5ef91582e70634dfbf2214fbdb6f29c949160b69f27ae947d919d
      sha512: 6695ad89f1a2ed7b54d358399c7b78e48f5259d7fb2c4e93d4b2b547d8b6e74c116a7bb0e41f2e56f0f29ba4bf2cc325c3a2ebdea0c0021d0788a886ecd37224
     blake2b: 8670dbbb9fa5da4aaa15e9aab7adf0097064867406ee7937d67a8720b0e9d466
  blake2b512: 3671ff8f2c9001fe6240773318d9f3f6957c7d4649e1cb435176836091c239e64a794faaeda253d5da16f510fba0c70e01b314b991fe2bd6b8402336e269c0b9

$ checksum --algs sha256,blake2b /bin/ls
Checksums:
      sha256: 5abf61c361e5ef91582e70634dfbf2214fbdb6f29c949160b69f27ae947d919d
     blake2b: 8670dbbb9fa5da4aaa15e9aab7adf0097064867406ee7937d67a8720b0e9d466
```

## License

**MIT**


## Contributing

All PRs are welcome.

1. Fork it
2. Create your feature branch (git checkout -b my-new-feature)
3. Commit your changes (git commit -am 'Add some feature')
5. Push to the branch (git push origin my-new-feature)
6. Create new Pull Request
