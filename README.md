# openconstruct-go

Go bindings for OpenConstruct backend services, Kubernetes operators, and cloud-native agent deployments.

## Install

```bash
go get github.com/superinstance/openconstruct-go
```

## Usage

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
    fmt.Println("Available math modules:", len(modules))

    client.SelectModules([]string{"spectral-graph-core", "plato-room"})
    client.ChooseInterface("sdk")

    config := client.GenerateConfig()
    fmt.Printf("Config: %+v\n", config)
}
```

## License

MIT
