import test from "node:test";
import assert from "node:assert/strict";

import { shouldOfferTripCreation } from "./aiPlanIntent.js";

test("does not offer trip creation when assistant is still asking follow-up questions", () => {
  const userText = "帮我规划三天杭州行程，预算3000元，喜欢美食和轻松节奏";
  const replyText = `1. 浙江具体想去哪个城市或区域？
比如杭州（西湖、灵隐、龙井）、乌镇/西塘（古镇慢生活）、莫干山（山野度假）、宁波（海鲜+老城）、绍兴（黄酒+水乡），或者想串一条线？
2. 大概几月出发？
不同季节穿搭和出行提醒完全不同。

你告诉我这两个，我就能结合你喜欢的“美食”和“轻松节奏”，把三天的住宿、路线、避坑和预算分配都排出来。`;

  assert.equal(shouldOfferTripCreation(userText, replyText), false);
});

test("does not offer trip creation when assistant only gives generic suggestions", () => {
  const userText = "明天从宿迁到浙江游玩3天人数情侣两人预算";
  const replyText = `因为预算充裕，无论哪个方向都建议：
+- 交通：全程包车或网约专车，不挤地铁
+- 餐饮：提前1-2天预订黑珍珠或金梧桐餐厅，减少落空
+- 穿搭：浙江当前季节多为闷热或雨，建议轻便透气衣物+防晒+折叠伞

如果你能马上告诉我具体城市或区域，我可以立刻把这三天的酒店、三餐、路线和预算分配全部排出来。`;

  assert.equal(shouldOfferTripCreation(userText, replyText), false);
});

test("does not offer trip creation when assistant is still giving options for the user to choose", () => {
  const userText = "想去浙江玩三天，喜欢轻松和美食";
  const replyText = `我先给你三个方向，你挑一个我再给你排完整行程：
1. 杭州：西湖、龙井、黑珍珠餐厅，整体最省心。
2. 乌镇西塘：古镇慢生活，适合情侣拍照。
3. 莫干山：民宿和山野度假，更安静。

你选定一个城市或路线后，我再把每天的酒店、三餐、交通和预算细化出来。`;

  assert.equal(shouldOfferTripCreation(userText, replyText), false);
});

test("offers trip creation only when reply contains a real day-by-day itinerary", () => {
  const userText = "帮我规划苏州三天两晚旅行，预算4000元";
  const replyText = `第1天：上午到达平江路，下午拙政园和苏州博物馆，晚上观前街用餐。
交通：地铁加步行最方便。
住宿：建议住在姑苏区地铁口附近。
预算：首日约1200元。

第2天：上午留园，下午山塘街，晚上听评弹。
餐饮：午餐安排苏帮菜，晚餐安排小吃。
注意事项：热门园林建议提前预约。

第3天：上午金鸡湖，下午返程。
预算：总预算控制在4000元以内。`;

  assert.equal(shouldOfferTripCreation(userText, replyText), true);
});
