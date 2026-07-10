package traveldata

import "testing"

func TestSplitBudgetKeepsTotalAmount(t *testing.T) {
	items := SplitBudget(BudgetRequest{
		Total:  3000,
		Days:   3,
		People: 2,
	})

	if len(items) != 5 {
		t.Fatalf("expected five budget items, got %d", len(items))
	}
	total := 0
	labels := map[string]bool{}
	for _, item := range items {
		total += item.Amount
		labels[item.Label] = true
	}
	if total != 3000 {
		t.Fatalf("budget total mismatch: %d", total)
	}
	for _, label := range []string{"住宿", "餐饮", "交通", "门票", "机动"} {
		if !labels[label] {
			t.Fatalf("missing budget label %s in %#v", label, items)
		}
	}
}

func TestSplitBudgetUsesDailyFallbackWhenNoTotalProvided(t *testing.T) {
	items := SplitBudget(BudgetRequest{Days: 2, People: 1})

	total := 0
	for _, item := range items {
		total += item.Amount
	}
	if total != 1200 {
		t.Fatalf("expected fallback budget 1200, got %d", total)
	}
}
