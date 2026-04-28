package env

import (
	"testing"

	"github.com/patchwork-deploy/internal/state"
)

func budgetBaseState() *state.State {
	st := state.New()
	st.AddEnvironment("prod")
	st.Record("prod", "001-init.sql")
	st.Record("prod", "002-index.sql")
	return st
}

func TestSetBudget_Success(t *testing.T) {
	st := budgetBaseState()
	if err := SetBudget(st, "prod", "001-init.sql", 3); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	limit, ok := GetBudget(st, "prod", "001-init.sql")
	if !ok || limit != 3 {
		t.Errorf("expected budget 3, got %d (ok=%v)", limit, ok)
	}
}

func TestSetBudget_ZeroLimitReturnsError(t *testing.T) {
	st := budgetBaseState()
	if err := SetBudget(st, "prod", "001-init.sql", 0); err == nil {
		t.Fatal("expected error for zero limit")
	}
}

func TestSetBudget_MissingEnvReturnsError(t *testing.T) {
	st := budgetBaseState()
	if err := SetBudget(st, "staging", "001-init.sql", 5); err == nil {
		t.Fatal("expected error for missing environment")
	}
}

func TestSetBudget_MissingPatchReturnsError(t *testing.T) {
	st := budgetBaseState()
	if err := SetBudget(st, "prod", "999-ghost.sql", 5); err == nil {
		t.Fatal("expected error for missing patch")
	}
}

func TestGetBudget_NoneSet(t *testing.T) {
	st := budgetBaseState()
	_, ok := GetBudget(st, "prod", "001-init.sql")
	if ok {
		t.Error("expected no budget to be set")
	}
}

func TestClearBudget_RemovesEntry(t *testing.T) {
	st := budgetBaseState()
	_ = SetBudget(st, "prod", "001-init.sql", 2)
	ClearBudget(st, "prod", "001-init.sql")
	_, ok := GetBudget(st, "prod", "001-init.sql")
	if ok {
		t.Error("expected budget to be cleared")
	}
}

func TestCheckBudget_NotExceeded(t *testing.T) {
	st := budgetBaseState()
	_ = SetBudget(st, "prod", "001-init.sql", 5)
	if err := CheckBudget(st, "prod", "001-init.sql"); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestCheckBudget_Exceeded(t *testing.T) {
	st := budgetBaseState()
	_ = SetBudget(st, "prod", "001-init.sql", 1)
	if err := CheckBudget(st, "prod", "001-init.sql"); err == nil {
		t.Error("expected budget exceeded error")
	}
}

func TestListBudgets_ReturnsAllEntries(t *testing.T) {
	st := budgetBaseState()
	_ = SetBudget(st, "prod", "001-init.sql", 2)
	_ = SetBudget(st, "prod", "002-index.sql", 4)
	budgets := ListBudgets(st, "prod")
	if len(budgets) != 2 {
		t.Errorf("expected 2 budgets, got %d", len(budgets))
	}
	if budgets["001-init.sql"] != 2 {
		t.Errorf("expected budget 2 for 001-init.sql")
	}
}
