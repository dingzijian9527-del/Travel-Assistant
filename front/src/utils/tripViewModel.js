const EMPTY_TEXT = "";

export function normalizeTripView(trip) {
  if (!isValidTrip(trip)) {
    return {
      hasTrip: false,
      title: "还没有保存行程",
      subtitle: "先和旅行管家定下目的地与时间，再生成你的第一份行程。",
      summary: createEmptySummary(),
      days: [],
      budgetItems: [],
      budgetTotal: "",
      alerts: [],
    };
  }

  const days = buildDays(trip);
  const summary = buildSummary(trip, days);
  const budget = buildBudget(trip);

  return {
    hasTrip: true,
    title: buildTitle(trip, summary),
    subtitle: buildSubtitle(trip, summary),
    summary,
    days,
    budgetItems: budget.items,
    budgetTotal: budget.total,
    alerts: buildAlerts(trip, days),
  };
}

export function moneyToNumber(value) {
  const text = String(value || "")
    .replace(/,/g, "")
    .replace(/[¥￥]/g, "")
    .replace(/元/g, "")
    .trim();
  const matches = text.match(/\d+(?:\.\d+)?/g);
  if (!matches?.length) {
    return 0;
  }

  const numbers = matches.map((item) => Number(item)).filter((item) => Number.isFinite(item));
  if (!numbers.length) {
    return 0;
  }
  if (numbers.length >= 2 && /[-~至到]/.test(text)) {
    return Math.round((numbers[0] + numbers[1]) / 2);
  }
  return Math.round(numbers[0]);
}

function isValidTrip(trip) {
  return !!trip && typeof trip === "object" && (hasText(trip.title) || hasText(trip.destination) || hasItems(trip.days));
}

function buildTitle(trip, summary) {
  if (hasText(trip.title)) {
    return trip.title.trim();
  }
  const destination = textOr(getField(trip, "destination", "destination_name"), "");
  if (destination && summary.days) {
    return `${destination} / ${summary.days}`;
  }
  return destination || summary.days || "";
}

function buildSubtitle(trip, summary) {
  if (hasText(trip.subtitle)) {
    return trip.subtitle.trim();
  }
  const parts = [summary.date, summary.people, summary.budget].filter(hasText);
  if (parts.length) {
    return parts.join(" / ");
  }
  if (hasText(trip.sourceQuestion)) {
    return trip.sourceQuestion.trim().slice(0, 34);
  }
  return "";
}

function buildSummary(trip, days) {
  const summary = trip.summary && typeof trip.summary === "object" ? trip.summary : {};
  const dayCount = Number(getField(trip, "dayCount", "day_count")) || days.length || 0;

  return {
    date: textOr(summary.date, textOr(getField(trip, "dateRange", "date_range"), EMPTY_TEXT)),
    days: textOr(summary.days, dayCount > 0 ? `${dayCount}天` : EMPTY_TEXT),
    people: textOr(summary.people, textOr(trip.people, EMPTY_TEXT)),
    budget: textOr(summary.budget, textOr(getField(trip, "budgetLevel", "budget_level"), EMPTY_TEXT)),
  };
}

function buildDays(trip) {
  if (!hasItems(trip.days)) {
    return [];
  }
  return trip.days
    .map((item, index) => normalizeDay(item, index))
    .filter(Boolean)
    .sort((a, b) => a.day - b.day);
}

function normalizeDay(item, index) {
  if (!item || typeof item !== "object") {
    return null;
  }
  const day = Number(item.day) || index + 1;
  return {
    day,
    title: textOr(item.title, ""),
    route: textOr(item.route, ""),
    food: textOr(item.food, ""),
    hotel: textOr(item.hotel, ""),
    tips: normalizeTips(item),
    weather: textOr(item.weather, ""),
  };
}

function normalizeTips(day) {
  const explicitTips = Array.isArray(day.tips)
    ? day.tips
    .map((item) => ({
      icon: textOr(item?.icon, "提"),
      title: textOr(item?.title, "行程提示"),
      text: textOr(item?.text, ""),
    }))
    .filter((item) => hasText(item.text))
    : [];

  if (explicitTips.length) {
    return explicitTips;
  }

  const builtTips = [];
  if (hasText(day.food)) {
    builtTips.push({ icon: "食", title: "餐饮建议", text: day.food.trim() });
  }
  if (hasText(day.hotel)) {
    builtTips.push({ icon: "住", title: "住宿建议", text: day.hotel.trim() });
  }
  return builtTips;
}

function buildBudget(trip) {
  const items = (Array.isArray(trip.budget) ? trip.budget : [])
    .map((item) => ({
      label: textOr(item?.label, ""),
      amount: normalizeMoneyText(item?.amount),
      value: moneyToNumber(item?.amount),
    }))
    .filter((item) => item.value > 0);

  if (!items.length) {
    return { items: [], total: "" };
  }
  return withBudgetPercent(items);
}

function withBudgetPercent(items) {
  const total = items.reduce((sum, item) => sum + item.value, 0);
  return {
    items: items.map((item) => ({
      label: item.label,
      amount: item.amount,
      percent: Math.max(8, Math.round((item.value / total) * 100)),
    })),
    total: formatMoney(total),
  };
}

function buildAlerts(trip) {
  if (!Array.isArray(trip.alerts)) {
    return [];
  }
  return trip.alerts.filter(hasText).map((item) => item.trim()).slice(0, 4);
}

function createEmptySummary() {
  return {
    date: EMPTY_TEXT,
    days: EMPTY_TEXT,
    people: EMPTY_TEXT,
    budget: EMPTY_TEXT,
  };
}

function getField(source, camelKey, snakeKey) {
  return source?.[camelKey] ?? source?.[snakeKey];
}

function hasItems(value) {
  return Array.isArray(value) && value.length > 0;
}

function hasText(value) {
  return typeof value === "string" && value.trim().length > 0;
}

function textOr(value, fallback) {
  return hasText(value) ? value.trim() : fallback;
}

function normalizeMoneyText(value) {
  const number = moneyToNumber(value);
  return number > 0 ? formatMoney(number) : textOr(value, "");
}

function formatMoney(value) {
  const amount = Math.round(Number(value) || 0);
  return `¥${amount.toLocaleString("zh-CN")}`;
}
