package openconstruct

import (
	"strings"
	"testing"
)

func TestStart(t *testing.T) {
	c := NewClient()
	defer c.Close()
	if c.Phase() != PhaseInit {
		t.Fatalf("expected init, got %s", c.Phase())
	}
	c.Start()
	if c.Phase() != PhaseStarted {
		t.Fatalf("expected started, got %s", c.Phase())
	}
}

func TestDeclareAgent(t *testing.T) {
	c := NewClient()
	defer c.Close()
	c.Start()
	identity := AgentIdentity{Name: "test-agent", Model: "glm-5.1", Capabilities: []string{"code_generation"}}
	c.DeclareAgent(identity)
	if c.Phase() != PhaseAgentDeclared {
		t.Fatalf("expected agent_declared, got %s", c.Phase())
	}
}

func TestListModules(t *testing.T) {
	c := NewClient()
	defer c.Close()
	modules := c.ListModules()
	if len(modules) == 0 {
		t.Fatal("expected some modules")
	}
}

func TestListModulesFilter(t *testing.T) {
	c := NewClient()
	defer c.Close()
	modules := c.ListModules(WithDomain("math"))
	for _, m := range modules {
		if m.Domain != "math" {
			t.Fatalf("expected math domain, got %s", m.Domain)
		}
	}
	if len(modules) == 0 {
		t.Fatal("expected at least one math module")
	}
}

func TestSelectModules(t *testing.T) {
	c := NewClient()
	defer c.Close()
	c.Start()
	c.DeclareAgent(AgentIdentity{Name: "a", Model: "m"})
	c.SelectModules([]string{"spectral-graph-core", "plato-room"})
	if c.Phase() != PhaseModulesSelected {
		t.Fatalf("expected modules_selected, got %s", c.Phase())
	}
}

func TestChooseInterface(t *testing.T) {
	c := NewClient()
	defer c.Close()
	c.Start()
	c.DeclareAgent(AgentIdentity{Name: "a", Model: "m"})
	c.SelectModules([]string{"spectral-graph-core"})
	c.ChooseInterface("sdk")
	if c.Phase() != PhaseInterfaceChosen {
		t.Fatalf("expected interface_chosen, got %s", c.Phase())
	}
}

func TestGenerateConfig(t *testing.T) {
	c := NewClient()
	defer c.Close()
	c.Start()
	c.DeclareAgent(AgentIdentity{Name: "a", Model: "m"})
	c.SelectModules([]string{"spectral-graph-core"})
	c.ChooseInterface("cli")
	cfg := c.GenerateConfig()
	if cfg.Phase != PhaseComplete {
		t.Fatalf("expected complete, got %s", cfg.Phase)
	}
	if cfg.Agent.Name != "a" {
		t.Fatalf("expected agent name 'a', got %s", cfg.Agent.Name)
	}
	if cfg.Interface != "cli" {
		t.Fatalf("expected interface 'cli', got %s", cfg.Interface)
	}
}

func TestFullLifecycle(t *testing.T) {
	c := NewClient()
	defer c.Close()

	c.Start()
	if c.Phase() != PhaseStarted {
		t.Fatalf("phase: %s", c.Phase())
	}

	c.DeclareAgent(AgentIdentity{
		Name:         "lifecycle-agent",
		Model:        "glm-5.1",
		Capabilities: []string{"code_generation", "testing"},
	})
	if c.Phase() != PhaseAgentDeclared {
		t.Fatalf("phase: %s", c.Phase())
	}

	mods := c.ListModules(WithDomain("math"))
	if len(mods) < 2 {
		t.Fatalf("expected >=2 math modules, got %d", len(mods))
	}

	var ids []string
	for _, m := range mods {
		ids = append(ids, m.ID)
	}
	c.SelectModules(ids)
	if c.Phase() != PhaseModulesSelected {
		t.Fatalf("phase: %s", c.Phase())
	}

	c.ChooseInterface("api")
	if c.Phase() != PhaseInterfaceChosen {
		t.Fatalf("phase: %s", c.Phase())
	}

	cfg := c.GenerateConfig()
	if cfg.Phase != PhaseComplete {
		t.Fatalf("phase: %s", cfg.Phase)
	}
	if len(cfg.Modules) != len(mods) {
		t.Fatalf("expected %d modules in config, got %d", len(mods), len(cfg.Modules))
	}
	if !strings.HasPrefix(cfg.SessionID, "oc-") {
		t.Fatalf("session ID should start with 'oc-', got %s", cfg.SessionID)
	}
}

func TestUniqueSessionID(t *testing.T) {
	ids := make(map[string]bool)
	for i := 0; i < 100; i++ {
		c := NewClient()
		sid := c.SessionID()
		if ids[sid] {
			t.Fatalf("duplicate session ID: %s", sid)
		}
		ids[sid] = true
		c.Close()
	}
}
