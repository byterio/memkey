# memkey

memkey is a small, in-memory key-value store for Go.

## Features

- Simple in-memory key-value store (map-backed).
- Optional TTL per key.
- Background cleanup to remove expired keys.
- Concurrent-safe.

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

    if mk.Has("ping") {
        v, _ := mk.Get("ping")
        fmt.Println(string(v)) // prints "pong"
    }

    time.Sleep(6 * time.Second)
    if !mk.Has("ping") {
        fmt.Println("expired")
    }
}
```

## Feedback and Contributions

If you encounter any issues or have suggestions for improvement, please [open an issue](https://github.com/byterio/memkey/issues) on GitHub.

We welcome contributions! Fork the repository, make your changes, and submit a pull request.

## Support

If you enjoy using memkey, please consider giving it a star! Your support helps others discover the project and encourages further development.

## License

memkey is open-source software released under the Apache License, Version 2.0. You can find a copy of the license in the [LICENSE](https://github.com/byterio/memkey/blob/main/LICENSE) file.
