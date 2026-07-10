import assert from "node:assert/strict";
import { describe, it } from "node:test";
import { moneyToNumber, normalizeTripView } from "./tripViewModel.js";

describe("normalizeTripView", () => {
  it("returns an empty view when no saved trip exists", () => {
    const view = normalizeTripView(null);

    assert.equal(view.hasTrip, false);
    assert.equal(view.title, "还没有保存行程");
    assert.equal(view.subtitle, "先和旅行管家定下目的地与时间，再生成你的第一份行程。");
    assert.deepEqual(view.days, []);
    assert.deepEqual(view.budgetItems, []);
    assert.equal(view.budgetTotal, "");
    assert.deepEqual(view.alerts, []);
  });

  it("uses saved trip data for summary, days, budget and alerts", () => {
    const view = normalizeTripView({
      title: "杭州 / 3天",
      subtitle: "西湖 / 灵隐寺 / 河坊街",
      dateRange: "07月02日到07月04日",
      dayCount: 3,
      people: "2人",
      summary: { budget: "¥4,000" },
      days: [
        {
          day: 1,
          title: "西湖慢游",
          route: "湖滨银泰 -> 西湖 -> 南山路",
          food: "晚餐建议湖滨附近杭帮菜",
          hotel: "住湖滨或武林广场",
          weather: "多云，注意防晒",
        },
      ],
      budget: [
        { label: "住宿", amount: "¥1,600" },
        { label: "餐饮", amount: "¥900" },
        { label: "交通", amount: "¥500" },
      ],
      alerts: ["灵隐寺建议提前预约。"],
    });

    assert.equal(view.hasTrip, true);
    assert.equal(view.title, "杭州 / 3天");
    assert.equal(view.summary.date, "07月02日到07月04日");
    assert.equal(view.summary.days, "3天");
    assert.equal(view.summary.people, "2人");
    assert.equal(view.days[0].title, "西湖慢游");
    assert.equal(view.days[0].tips.length, 2);
    assert.equal(view.budgetItems[0].label, "住宿");
    assert.equal(view.budgetTotal, "¥3,000");
    assert.deepEqual(view.alerts, ["灵隐寺建议提前预约。"]);
  });

  it("does not invent fallback content when backend details are missing", () => {
    const view = normalizeTripView({
      destination: "厦门",
      dayCount: 2,
      summary: { budget: "¥2,800-3,200" },
      days: [{ day: 1 }],
    });

    assert.equal(view.title, "厦门 / 2天");
    assert.equal(view.days[0].title, "");
    assert.equal(view.days[0].route, "");
    assert.equal(view.days[0].weather, "");
    assert.deepEqual(view.days[0].tips, []);
    assert.deepEqual(view.budgetItems, []);
    assert.equal(view.budgetTotal, "");
    assert.deepEqual(view.alerts, []);
  });
});

describe("moneyToNumber", () => {
  it("parses regular money text and money ranges", () => {
    assert.equal(moneyToNumber("¥1,600"), 1600);
    assert.equal(moneyToNumber("¥2,800-3,200"), 3000);
    assert.equal(moneyToNumber("预算约3000元"), 3000);
  });
});
