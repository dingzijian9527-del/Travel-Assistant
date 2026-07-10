package rpcaiagent

import (
	"strings"

	"github.com/dingzijian9527-del/Travel-Assistant/pkg/repository"
)

type chatMessageModel = repository.ChatMessage

func buildSuggestions(message string) []string {
	if strings.TrimSpace(message) == "" || !isTravelRelated(message) {
		return []string{"定制 5 天大理行程", "制作西安旅游攻略", "回答三亚亲子游问题"}
	}
	if containsAny(message, []string{"酒店", "民宿", "住宿"}) {
		return []string{"按预算筛选酒店", "比较民宿和酒店", "推荐交通方便的住宿区域"}
	}
	if containsAny(message, []string{"攻略", "指南", "避坑"}) {
		return []string{"生成每日路线", "整理门票预约提醒", "补充当地美食清单"}
	}
	return []string{"继续细化预算", "增加美食安排", "生成每日路线"}
}

func isTravelRelated(message string) bool {
	return containsAny(message, []string{
		"旅游", "旅行", "出行", "行程", "攻略", "目的地", "景点", "路线", "酒店", "民宿",
		"住宿", "机票", "高铁", "交通", "天气", "美食", "餐厅", "小吃", "门票", "签证",
		"预算", "亲子游", "海岛", "古镇", "徒步", "自驾", "大理", "上海", "三亚", "北京",
		"西安", "成都", "杭州", "广州", "新疆", "云南",
	})
}

func containsAny(message string, keywords []string) bool {
	normalized := strings.ToLower(strings.TrimSpace(message))
	for _, keyword := range keywords {
		if strings.Contains(normalized, strings.ToLower(keyword)) {
			return true
		}
	}
	return false
}
