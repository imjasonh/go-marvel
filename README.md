go-marvel
=========

[Go](https://golang.org) client library for the [Marvel REST API](https://developer.marvel.com/)

Usage:

```go
import (
  "fmt"
  
  "github.com/ImJasonH/go-marvel"
)

func main() {
  c := marvel.NewClient("my-public-key", "my-private-key")
  r, err := c.Series(2258, marvel.CommonRequest{})
  if err != nil {
    panic(err)
  }
  fmt.Println("%+v\n", r)
}
```
