export function shouldOfferTripCreation(userText, replyText) {
  const user = normalize(userText);
  const reply = normalize(replyText);
  if (!user || !reply) {
    return false;
  }
  if (!hasPlanIntent(user)) {
    return false;
  }
  if (hasClarifyingSignals(reply) || hasChoiceSignals(reply)) {
    return false;
  }
  return hasRealPlanStructure(reply);
}

function normalize(text) {
  return String(text || "").trim();
}

function hasPlanIntent(text) {
  return /行程|规划|路线|旅行计划|专属旅行|几日游|旅游|出行|安排/.test(text);
}

function hasClarifyingSignals(text) {
  return /告诉我|请告诉我|如果你能|还需要|请先|先告诉|补充|具体想去|具体城市|具体区域|哪个城市|哪个区域|几月出发|大概几月|你更想|你想先|回复我|选定后|确定后|再把/.test(text);
}

function hasChoiceSignals(text) {
  return /先给你.*方向|给你.*方向|供你挑|你挑一个|可以选|三个方向|几个方向|或者想串|你来挑|选一个/.test(text);
}

function hasRealPlanStructure(text) {
  const dayMatches = text.match(/第\s*[一二三四五六七八九十\d]+\s*天/g) || [];
  const sectionMatches = text.match(/(交通|住宿|预算|餐饮|美食|注意事项|提醒)[：:]/g) || [];
  return dayMatches.length >= 2 || (dayMatches.length >= 1 && sectionMatches.length >= 2 && text.length >= 120);
}
