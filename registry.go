package openconstruct

// ModuleRegistry holds available modules.
type ModuleRegistry struct {
	modules []Module
}

// NewModuleRegistry returns a registry preloaded with built-in modules.
func NewModuleRegistry() *ModuleRegistry {
	return &ModuleRegistry{
		modules: []Module{
			{ID: "spectral-graph-core", Name: "Spectral Graph Core", Description: "Graph spectral analysis engine", Domain: "math"},
			{ID: "plato-room", Name: "Plato Room", Description: "Collaborative reasoning space", Domain: "math"},
			{ID: "code-gen-engine", Name: "Code Generation Engine", Description: "Autonomous code generation", Domain: "code"},
			{ID: "test-runner", Name: "Test Runner", Description: "Automated test execution", Domain: "code"},
			{ID: "doc-writer", Name: "Doc Writer", Description: "Documentation generation", Domain: "docs"},
			{ID: "api-designer", Name: "API Designer", Description: "REST/GraphQL API scaffolding", Domain: "code"},
			{ID: "tensor-flow-lite", Name: "Tensor Flow Lite", Description: "Lightweight ML inference", Domain: "math"},
			{ID: "deploy-agent", Name: "Deploy Agent", Description: "Cloud deployment automation", Domain: "infra"},
		},
	}
}

// List returns all modules matching the optional filters.
func (r *ModuleRegistry) List(opts ...ListModulesOption) []Module {
	cfg := &listModulesConfig{}
	for _, o := range opts {
		o(cfg)
	}
	var result []Module
	for _, m := range r.modules {
		if cfg.domain != "" && m.Domain != cfg.domain {
			continue
		}
		result = append(result, m)
	}
	return result
}

// FindByID returns a module by its ID.
func (r *ModuleRegistry) FindByID(id string) (Module, bool) {
	for _, m := range r.modules {
		if m.ID == id {
			return m, true
		}
	}
	return Module{}, false
}
