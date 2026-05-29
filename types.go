package openconstruct

// AgentIdentity describes the agent being constructed.
type AgentIdentity struct {
	Name         string
	Model        string
	Capabilities []string
}

// Module represents a selectable onboarding module.
type Module struct {
	ID          string
	Name        string
	Description string
	Domain      string
}

// OnboardingConfig is the final generated configuration.
type OnboardingConfig struct {
	Agent       AgentIdentity
	Modules     []Module
	Interface   string
	SessionID   string
	Phase       string
}

// ListModulesOption configures module listing.
type ListModulesOption func(*listModulesConfig)

type listModulesConfig struct {
	domain string
}

// WithDomain filters modules by domain.
func WithDomain(domain string) ListModulesOption {
	return func(cfg *listModulesConfig) {
		cfg.domain = domain
	}
}
