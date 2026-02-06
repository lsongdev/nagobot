package thread

import (
	"testing"

	"github.com/linanwx/nagobot/agent"
)

func TestThreadCompositionConstructors(t *testing.T) {
	plain := NewPlain(nil, nil, nil)
	if plain == nil || plain.Thread == nil {
		t.Fatalf("plain thread should be initialized")
	}
	if plain.Type() != ThreadTypePlain {
		t.Fatalf("plain type mismatch: got %s", plain.Type())
	}
	if plain.sessionKey != "" {
		t.Fatalf("plain sessionKey should be empty, got %q", plain.sessionKey)
	}

	channel := NewChannel(nil, nil, " chat:user ", nil)
	if channel == nil || channel.PlainThread == nil || channel.Thread == nil {
		t.Fatalf("channel thread composition should be initialized")
	}
	if channel.Type() != ThreadTypeChannel {
		t.Fatalf("channel type mismatch: got %s", channel.Type())
	}
	if channel.sessionKey != "chat:user" {
		t.Fatalf("channel sessionKey mismatch: got %q", channel.sessionKey)
	}

	child := NewChild(nil, nil, nil)
	if child == nil || child.PlainThread == nil || child.Thread == nil {
		t.Fatalf("child thread composition should be initialized")
	}
	if child.Type() != ThreadTypeChild {
		t.Fatalf("child type mismatch: got %s", child.Type())
	}
	if child.allowSpawn {
		t.Fatalf("child thread must not allow spawn")
	}
}

func TestManagerGetOrCreateChannelReusesChannelThread(t *testing.T) {
	mgr := NewManager(&Config{})
	ag := agent.NewRawAgent("a", "prompt")

	first := mgr.GetOrCreateChannel("room:1", ag, nil)
	second := mgr.GetOrCreateChannel("room:1", nil, nil)
	if first != second {
		t.Fatalf("expected channel thread reuse")
	}
	if first.Type() != ThreadTypeChannel {
		t.Fatalf("expected channel type, got %s", first.Type())
	}
}

func TestManagerGetOrCreatePlainReturnsFreshThread(t *testing.T) {
	mgr := NewManager(&Config{})

	first := mgr.GetOrCreate("", nil, nil)
	second := mgr.GetOrCreate("", nil, nil)
	if first == second {
		t.Fatalf("plain threads should not be cached")
	}
	if first.Type() != ThreadTypePlain || second.Type() != ThreadTypePlain {
		t.Fatalf("plain thread type mismatch: %s / %s", first.Type(), second.Type())
	}
}
