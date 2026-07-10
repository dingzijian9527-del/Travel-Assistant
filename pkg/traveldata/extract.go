package traveldata

import "regexp"

func InferDestination(message string) string {
	patterns := []*regexp.Regexp{
		regexp.MustCompile(`(?:去|到|玩|游|规划|安排)([\p{Han}]{2,8})(?:三天|两天|一天|四天|五天|六天|七天|旅游|旅行|行程|怎么玩|怎么安排|$)`),
		regexp.MustCompile(`目的地[:：\s]*([\p{Han}]{2,8})`),
	}
	for _, pattern := range patterns {
		matches := pattern.FindStringSubmatch(message)
		if len(matches) >= 2 {
			return matches[1]
		}
	}
	for _, city := range []string{"成都", "杭州", "广州", "上海", "北京", "西安", "三亚", "大理", "重庆", "南京", "苏州", "厦门", "青岛"} {
		if regexp.MustCompile(city).MatchString(message) {
			return city
		}
	}
	return ""
}

func InferDays(message string) int {
	if matches := regexp.MustCompile(`(\d+)\s*天`).FindStringSubmatch(message); len(matches) >= 2 {
		return atoi(matches[1])
	}
	values := map[string]int{"一天": 1, "两天": 2, "二天": 2, "三天": 3, "四天": 4, "五天": 5, "六天": 6, "七天": 7}
	for text, value := range values {
		if regexp.MustCompile(text).MatchString(message) {
			return value
		}
	}
	return 0
}

func InferPeople(message string) int {
	if matches := regexp.MustCompile(`(\d+)\s*人`).FindStringSubmatch(message); len(matches) >= 2 {
		return atoi(matches[1])
	}
	return 0
}

func InferBudget(message string) int {
	matches := regexp.MustCompile(`(?:预算|花费|费用)?\s*(\d{3,6})\s*(?:元|块)`).FindStringSubmatch(message)
	if len(matches) >= 2 {
		return atoi(matches[1])
	}
	return 0
}

func atoi(value string) int {
	n := 0
	for _, ch := range value {
		if ch < '0' || ch > '9' {
			continue
		}
		n = n*10 + int(ch-'0')
	}
	return n
}
