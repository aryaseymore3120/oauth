# OpenAI Responses client

This optional Go package uses the standard library and does not change OAuth login behavior. OpenAI API authentication is separate from OAuth.

Set OPENAI_API_KEY securely in your local environment; never commit it. Select a model available to your API project. API usage may incur charges.

```go
package main

import (
    "context"
    "fmt"
    "log"
    "github.com/cli/oauth/openai"
)

func main() {
    client, err := openai.NewClient()
    if err != nil { log.Fatal(err) }
    text, err := client.Generate(context.Background(), "gpt-6.2-astra", "Hello!")
    if err != nil { log.Fatal(err) }
    fmt.Println(text)
}
```

Run `go test ./openai` for local mocked tests; no real key or network API call is needed. Run `go test ./...` for the whole module.

See https://developers.openai.com/api/docs/quickstart for API setup.
