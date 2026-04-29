package env

import (
	"testing"
)

func TestAddEmptyGroupReturnsError(t *testing.T) {
	g := NewGroupIndex()
	if err := g.Add("", "my-set"); err == nil {
		t.Fatal("expected error for empty group name")
	}
}

func TestAddEmptySetNameReturnsError(t *testing.T) {
	g := NewGroupIndex()
	if err := g.Add("my-group", ""); err == nil {
		t.Fatal("expected error for empty set name")
	}
}

func TestAddAndRetrieveMembers(t *testing.T) {
	g := NewGroupIndex()
	_ = g.Add("prod", "api")
	_ = g.Add("prod", "worker")
	members, err := g.Members("prod")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(members) != 2 {
		t.Fatalf("expected 2 members, got %d", len(members))
	}
}

func TestAddDuplicateSetIsIdempotent(t *testing.T) {
	g := NewGroupIndex()
	_ = g.Add("staging", "api")
	_ = g.Add("staging", "api")
	members, _ := g.Members("staging")
	if len(members) != 1 {
		t.Fatalf("expected 1 member after duplicate add, got %d", len(members))
	}
}

func TestRemoveSetDeletesEmptyGroup(t *testing.T) {
	g := NewGroupIndex()
	_ = g.Add("dev", "only-set")
	_ = g.Remove("dev", "only-set")
	groups := g.Groups()
	if len(groups) != 0 {
		t.Fatalf("expected no groups after removing last member, got %v", groups)
	}
}

func TestRemoveSetFromMissingGroupReturnsError(t *testing.T) {
	g := NewGroupIndex()
	if err := g.Remove("ghost", "api"); err == nil {
		t.Fatal("expected error removing from missing group")
	}
}

func TestRemoveMissingSetReturnsError(t *testing.T) {
	g := NewGroupIndex()
	_ = g.Add("prod", "api")
	if err := g.Remove("prod", "missing"); err == nil {
		t.Fatal("expected error removing missing set")
	}
}

func TestMembersUnknownGroupReturnsEmpty(t *testing.T) {
	g := NewGroupIndex()
	members, err := g.Members("nope")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(members) != 0 {
		t.Fatalf("expected empty slice, got %v", members)
	}
}

func TestGroupsReturnsSortedNames(t *testing.T) {
	g := NewGroupIndex()
	_ = g.Add("zebra", "s1")
	_ = g.Add("alpha", "s2")
	_ = g.Add("mango", "s3")
	groups := g.Groups()
	expected := []string{"alpha", "mango", "zebra"}
	for i, name := range expected {
		if groups[i] != name {
			t.Fatalf("expected %s at index %d, got %s", name, i, groups[i])
		}
	}
}
