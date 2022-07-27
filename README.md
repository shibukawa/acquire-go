# acquire-go

[![Go Reference](https://pkg.go.dev/badge/github.com/shibukawa/acquire-go.svg)](https://pkg.go.dev/github.com/shibukawa/acquire-go)

Search target file/folder from child dir to parent dir.

## Usage

If you have the following folder structure and there are two test cases `cool_test.go` and `cmd/awesome/awesome_test.go` and both test uses `testdata/testconfig.json`:

```txt
awesome-your-tool
├── LICENSE
├── README.md
├── cmd
│   └── awesome
│      ├── awesome.go
│      └── awesome_test.go
├── go.mod
├── go.sum
├── cool.go
├── cool_test.go
└── testdata
    └── testconfig.json
```

You can get the config file path even if your current folder is in top or under cmd/awesome. This function search recursively from current folder to ancestor folders.

```go
paths, err := acquire.Acquire(acquire.File, "testdata/testconfig.json")
// paths = []string{"/abs/folder/testdata/testconfig.json"}
```

## Functions

### `Acquire(targetType Type, patterns ...string) (matches []string, err error)`

Basic function to search file/dir.

`targetType` should be one of the following constant:

* `acquire.File`
* `acquire.Dir`
* `acquire.All`

`patterns` should be file/dir name or glob patterns:

* `sample.txt`
* `*.json`

If no files/dirs match, it returns `acquire.ErrNotFound` error.

### `MustAcquire(targetType Type, patterns ...string) (matches []string)`

Call `panic()` if there is an error. It is good for test code.

### `AcquireUnder(targetType Type, under string, patterns ...string) (matches []string, err error)`

It searches only `under` folder. It is good for preventing directory traversal attack.

### `AcquireFromUnder(targetType Type, folder, under string, patterns ...string) (matches []string, err error)`

It is the most flexible version. It can specify search start folder.

## License

Apache 2

## Credit

Yoshiki Shibukawa