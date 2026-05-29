# openconstruct-go — Go Bindings for OpenConstruct

Go client for the [OpenConstruct](https://github.com/SuperInstance/OpenConstruct) ecosystem. Built for backend services, Kubernetes operators, and cloud-native agent deployments.

## What This Gives You

- **5-phase onboarding** — Start → DeclareAgent → SelectModules → ChooseInterface → GenerateConfig
- **Module registry** — built-in catalog with domain filtering
- **Idiomatic Go** — `Close()` for cleanup, struct-based types, error returns
- **Zero external dependencies** — just the standard library

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/superinstance/openconstruct-go"
)

func main() {
    client := openconstruct.NewClient()
    defer client.Close()

    client.Start()

    identity := openconstruct.AgentIdentity{
        Name:         "my-agent",
        Model:        "glm-5.1",
        Capabilities: []string{"code_generation"},
    }
    client.DeclareAgent(identity)

    modules := client.ListModules(openconstruct.WithDomain("math"))
    client.SelectModules([]string{"spectral-graph-core", "plato-room"})
    client.ChooseInterface("sdk")

    config := client.GenerateConfig()
    fmt.Printf("Config: %+v\n", config)
}
```

## Installation

```bash
go get github.com/superinstance/openconstruct-go
```

## How It Fits

One of the [polyglot OpenConstruct bindings](https://github.com/SuperInstance/OpenConstruct). See [openconstruct-examples](https://github.com/SuperInstance/openconstruct-examples) for a Go onboarding example.

## Testing

```bash
go test ./...
```

## License

MIT
