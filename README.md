# memkey

[![Go Version][GoVer-Image]][GoDoc-Url] [![License][License-Image]][License-Url] [![GoDoc][GoDoc-Image]][GoDoc-Url] [![Go Report Card][ReportCard-Image]][ReportCard-Url]

[GoVer-Image]: https://img.shields.io/badge/Go-1.26%2B-blue
[GoDoc-Url]: https://pkg.go.dev/github.com/byterio/memkey
[GoDoc-Image]: https://pkg.go.dev/badge/github.com/byterio/memkey.svg
[ReportCard-Url]: https://goreportcard.com/report/github.com/byterio/memkey
[ReportCard-Image]: https://goreportcard.com/badge/github.com/byterio/memkey?style=flat
[License-Url]: https://github.com/byterio/memkey/blob/main/LICENSE
[License-Image]: https://img.shields.io/github/license/byterio/memkey

memkey is a lightweight, in-memory key-value store for Go.

## Features

- 🚀 Fast, lightweight in-memory key-value store.
- ⏰ Optional TTL per key.
- ✨ Background cleanup to remove expired keys.
- 🔐 Concurrent-safe.

## Installation

```bash
go get -u github.com/byterio/memkey
```

## Usage

```go
package main

import (
    "fmt"
    "time"

    "github.com/byterio/memkey"
)

func main() {
    mk := memkey.New()
    defer mk.Close()

    mk.Set("ping", []byte("pong"), 5*time.Second)
    mk.Set("foo", []byte("bar"), 0)

    if mk.Has("ping") {
        v, _ := mk.Get("ping")
        fmt.Println(string(v)) // prints "pong"
    }

    time.Sleep(6 * time.Second)

    if mk.Has("foo") {
        v, _ := mk.Get("foo")
        fmt.Println(string(v)) // prints "bar"
    }

    if !mk.Has("ping") {
        fmt.Println("expired")
    }
}
```

## Benchmarks

### Benchmark on an Intel Core i7-6700HQ CPU @ 2.60GHz:

| Benchmark                       | Iterations | ns/op | B/op | allocs/op |
| ------------------------------- | ---------: | ----: | ---: | --------: |
| Benchmark_Memkey_Set-8          | 12,485,835 | 91.72 |    4 |         1 |
| Benchmark_Memkey_Get-8          | 30,052,064 | 39.77 |    0 |         0 |
| Benchmark_Memkey_SetAndDelete-8 |  6,038,889 | 194.4 |    4 |         1 |

## Feedback and Contributions

If you encounter any issues or have suggestions for improvement, please [open an issue](https://github.com/byterio/memkey/issues) on GitHub.

We welcome contributions! Fork the repository, make your changes, and submit a pull request.

## Support

If you enjoy using memkey, please consider giving it a star! Your support helps others discover the project and encourages further development.

## License

memkey is open-source software released under the Apache License, Version 2.0. You can find a copy of the license in the [LICENSE](https://github.com/byterio/memkey/blob/main/LICENSE) file.
