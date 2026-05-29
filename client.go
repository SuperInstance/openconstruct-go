package openconstruct

import (
	"crypto/rand"
	"fmt"
	"sync"
)

// Phase constants for the onboarding flow.
const (
	PhaseInit           = "init"
	PhaseStarted        = "started"
	PhaseAgentDeclared  = "agent_declared"
	PhaseModulesSelected = "modules_selected"
	PhaseInterfaceChosen = "interface_chosen"
	PhaseComplete       = "complete"
)

// Client holds session state for an OpenConstruct onboarding flow.
type Client struct {
	mu        sync.Mutex
	sessionID string
	phase     string
	agent     AgentIdentity
	modules   []Module
	iface     string
	registry  *ModuleRegistry
	closed    bool
}

// NewClient creates a new OpenConstruct client with a unique session ID.
func NewClient() *Client {
	return &Client{
		sessionID: generateSessionID(),
		phase:     PhaseInit,
		registry:  NewModuleRegistry(),
	}
}

// SessionID returns the unique session identifier.
func (c *Client) SessionID() string {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.sessionID
}

// Phase returns the current onboarding phase.
func (c *Client) Phase() string {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.phase
}

// Start begins the onboarding session.
func (c *Client) Start() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.phase == PhaseInit {
		c.phase = PhaseStarted
	}
}

// DeclareAgent registers the agent identity.
func (c *Client) DeclareAgent(identity AgentIdentity) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.agent = identity
	if c.phase == PhaseStarted {
		c.phase = PhaseAgentDeclared
	}
}

// ListModules returns available modules, optionally filtered.
func (c *Client) ListModules(opts ...ListModulesOption) []Module {
	c.mu.Lock()
	reg := c.registry
	c.mu.Unlock()
	return reg.List(opts...)
}

// SelectModules selects modules by ID.
func (c *Client) SelectModules(ids []string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	var selected []Module
	for _, id := range ids {
		if m, ok := c.registry.FindByID(id); ok {
			selected = append(selected, m)
		}
	}
	c.modules = selected
	if c.phase == PhaseAgentDeclared {
		c.phase = PhaseModulesSelected
	}
}

// ChooseInterface sets the interface type (e.g. "cli", "api", "sdk").
func (c *Client) ChooseInterface(iface string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.iface = iface
	if c.phase == PhaseModulesSelected {
		c.phase = PhaseInterfaceChosen
	}
}

// GenerateConfig produces the final OnboardingConfig.
func (c *Client) GenerateConfig() OnboardingConfig {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.phase != PhaseInterfaceChosen {
		return OnboardingConfig{Phase: c.phase, SessionID: c.sessionID}
	}
	c.phase = PhaseComplete
	return OnboardingConfig{
		Agent:     c.agent,
		Modules:   c.modules,
		Interface: c.iface,
		SessionID: c.sessionID,
		Phase:     PhaseComplete,
	}
}

// Close cleans up client resources.
func (c *Client) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.closed = true
}

func generateSessionID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return fmt.Sprintf("oc-%x", b)
}
