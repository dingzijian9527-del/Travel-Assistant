package traveldata

type BudgetRequest struct {
	Total  int
	Days   int
	People int
}

func SplitBudget(req BudgetRequest) []BudgetItem {
	days := req.Days
	if days <= 0 {
		days = 1
	}
	people := req.People
	if people <= 0 {
		people = 1
	}
	total := req.Total
	if total <= 0 {
		total = days * people * 600
	}
	weights := []struct {
		label string
		rate  int
		note  string
	}{
		{label: "住宿", rate: 40, note: "优先按交通方便区域控制总价"},
		{label: "餐饮", rate: 25, note: "预留本地特色餐和小吃预算"},
		{label: "交通", rate: 20, note: "覆盖市内通勤和必要打车"},
		{label: "门票", rate: 10, note: "覆盖景点预约和体验项目"},
	}
	items := make([]BudgetItem, 0, 5)
	used := 0
	for _, item := range weights {
		amount := total * item.rate / 100
		used += amount
		items = append(items, BudgetItem{Label: item.label, Amount: amount, Note: item.note})
	}
	items = append(items, BudgetItem{Label: "机动", Amount: total - used, Note: "用于天气变化、临时交通和加餐"})
	return items
}
