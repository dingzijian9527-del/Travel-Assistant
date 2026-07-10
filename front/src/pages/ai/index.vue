<template>
  <view class="page">
    <view class="phone-shell">
      <view class="chat-header">
        <view class="header-main">
          <view class="bot-avatar">旅</view>
          <view class="header-copy">
            <text class="bot-name">旅行管家</text>
            <text class="bot-status">在线陪你聊目的地、路线、预算和提醒</text>
          </view>
        </view>
        
      </view>

      <scroll-view class="tools-scroll" scroll-x show-scrollbar="false">
        <view class="quick-tools">
          <view v-for="item in abilities" :key="item.title" class="tool-chip" @tap="sendPrompt(item.prompt)">
            <text class="tool-icon">{{ item.icon }}</text>
            <view class="tool-copy">
              <text class="tool-title">{{ item.title }}</text>
              <text class="tool-desc">{{ item.desc }}</text>
            </view>
          </view>
        </view>
      </scroll-view>

      <view class="prompt-panel">
        <view class="panel-head">
          <view>
            <text class="panel-title">快速开始</text>
            <text class="panel-subtitle">选一句直接发，或者在下方自己输入</text>
          </view>
          <text class="panel-more" @tap="shufflePrompts">换一组</text>
        </view>
        <view class="prompt-list">
          <button v-for="item in visiblePrompts" :key="item" @tap="sendPrompt(item)">{{ item }}</button>
        </view>
        <view v-if="prompts.length > 2" class="prompt-toggle" @tap="togglePrompts">
          <text>{{ promptExpanded ? '收起建议' : '展开更多建议' }}</text>
        </view>
      </view>

      <scroll-view class="chat-scroll" scroll-y :scroll-into-view="lastId">
        <view v-for="item in messages" :key="item.id" :id="item.id" :class="['msg-wrap', item.role]">
          <view v-if="item.role === 'ai'" class="msg-side">
            <text class="msg-ai-icon">旅</text>
          </view>
          <view class="msg-content">
            <view v-if="item.role === 'ai' && item.id === 'msg-1'" class="welcome-card">
              <text class="welcome-kicker">专属服务入口</text>
              <text class="welcome-title">我是你的旅行管家</text>
              <text class="welcome-desc">告诉我目的地、时间、同行人数或预算，我会帮你整理路线、餐饮、住宿和注意事项。</text>
              <view class="welcome-tags">
                <text>多轮追问</text>
                <text>预算拆分</text>
                <text>真实行程</text>
              </view>
            </view>

            <view v-if="item.role === 'ai'" :class="['bubble', 'ai', item.id === 'msg-1' ? 'first-bubble' : '']">
              <view v-for="(block, index) in formatAiBlocks(item.text)" :key="`${item.id}-${index}`" :class="['ai-block', block.type]">
                <text v-if="block.type === 'title'" class="ai-block-title">{{ block.text }}</text>
                <text v-else class="ai-block-text">{{ block.text }}</text>
              </view>
            </view>
            <view v-else class="bubble user">{{ item.text }}</view>

            <text :class="['msg-time', item.role]">{{ item.time }}</text>
            <view v-if="item.card" class="plan-card-compact">
              <image class="plan-img" :src="item.card.image" mode="aspectFill" />
              <view class="plan-body">
                <view class="plan-rating">{{ item.card.score }}</view>
                <text class="plan-title">{{ item.card.title }}</text>
                <view class="plan-tags">
                  <text v-for="tag in item.card.tags" :key="tag">{{ tag }}</text>
                </view>
                <text class="plan-desc">{{ item.card.desc }}</text>
              </view>
            </view>
          </view>
        </view>
      </scroll-view>

      <view class="input-area">
        <view class="input-caption">直接描述你的旅行想法，越具体越好</view>
        <view class="input-bar">
          <input v-model="message" placeholder="输入你的旅行想法，比如去哪里、玩几天、预算多少" @confirm="send" />
          <button class="send-btn" :disabled="isSending" @tap="send">{{ isSending ? '...' : '发送' }}</button>
        </view>
      </view>
      <BottomNav active="ai" />
    </view>
  </view>
</template>

<script setup>
import { computed, nextTick, onMounted, ref, watch } from "vue";
import { onLoad } from "@dcloudio/uni-app";
import BottomNav from "../../components/BottomNav.vue";
import { createAiPageState } from "../../utils/aiPageState.js";
import { shouldOfferTripCreation } from "../../utils/aiPlanIntent.js";
import { createTripPageState } from "../../utils/tripPageState.js";

const apiBase = "http://127.0.0.1:8080";
const chatStreamUrl = `${apiBase}/api/v1/ai-stream`;
const createTripUrl = `${apiBase}/api/v1/trips`;
const pageState = createAiPageState();
const tripPageState = createTripPageState();
const defaultMessages = [
  { id: "msg-1", role: "ai", text: "你好，我是你的专属旅行管家。无论你想去哪里，我都能帮你整理国内出行路线、酒店建议、当地攻略和预算拆分。今天想聊聊哪里的旅行？", time: "10:30" },
];

const message = ref("");
const lastId = ref("msg-1");
const isSending = ref(false);
const pendingAutoQuestion = ref("");
const hasAutoSent = ref(false);
const pendingTripPlan = ref(null);
const promptGroupIndex = ref(0);
const promptExpanded = ref(false);
const messages = ref([...defaultMessages]);

const abilities = [
  { icon: "路", title: "规划行程", desc: "路线节奏", prompt: "请结合实时天气、路线通勤、景点顺序和我的预算，帮我规划一个三天两晚的国内旅行行程" },
  { icon: "住", title: "选住宿", desc: "区域预算", prompt: "请结合目的地交通、景点分布、天气和预算，帮我推荐住宿区域，并说明适合人群和避坑点" },
  { icon: "食", title: "找美食", desc: "本地吃法", prompt: "请结合目的地路线、用餐区域和预算，推荐当地美食、小吃、餐厅区域和避坑建议" },
  { icon: "天", title: "出行提醒", desc: "天气避坑", prompt: "请结合目的地天气、路线通勤、交通方式和预算，整理出行前提醒、穿搭、雨具和避坑事项" },
];

const promptGroups = [
  [
    "帮我规划三天成都行程，预算3000元，喜欢美食和轻松节奏",
    "预算3000元，三天两晚，推荐适合情侣出行的目的地",
    "杭州亲子周末怎么安排，请结合天气、路线和餐饮",
    "推荐广州早茶路线，并顺路安排半日游",
  ],
  [
    "北京亲子三日游怎么安排，预算5000元",
    "上海两天一晚住哪里方便，想少换酒店",
    "西安三天怎么吃怎么玩，请把门票和交通预算拆开",
    "三亚四天海岛游怎么安排，预算8000元",
  ],
  [
    "重庆两天美食路线怎么安排，别太赶",
    "大理五天慢旅行怎么规划，预算6000元",
    "苏州周末园林路线怎么走，顺便推荐住宿区域",
    "厦门三天亲子游，请结合天气和交通提醒",
  ],
];

const prompts = ref(promptGroups[0]);
const visiblePrompts = computed(() => promptExpanded.value ? prompts.value : prompts.value.slice(0, 2));

onMounted(() => {
  if (uni.hideTabBar) uni.hideTabBar();
  restoreSession();
  sendPendingQuestion();
});

onLoad((options) => {
  if (options && options.autoSend === "1" && options.question) {
    pendingAutoQuestion.value = decodeURIComponent(options.question);
  }
});

watch([message, messages, pendingTripPlan, promptGroupIndex, promptExpanded], () => {
  persistSession();
}, { deep: true });

function restoreSession() {
  const session = pageState.getSession();
  if (!session || typeof session !== "object") {
    return;
  }

  message.value = typeof session.draftMessage === "string" ? session.draftMessage : "";
  promptGroupIndex.value = normalizePromptGroupIndex(session.promptGroupIndex);
  prompts.value = promptGroups[promptGroupIndex.value];
  promptExpanded.value = !!session.promptExpanded;
  pendingTripPlan.value = normalizePendingTripPlan(session.pendingTripPlan);

  if (Array.isArray(session.messages) && session.messages.length) {
    messages.value = session.messages;
    lastId.value = session.messages[session.messages.length - 1]?.id || "msg-1";
  }
}

function persistSession() {
  pageState.saveSession({
    draftMessage: message.value,
    promptGroupIndex: promptGroupIndex.value,
    promptExpanded: promptExpanded.value,
    pendingTripPlan: pendingTripPlan.value,
    messages: messages.value,
  });
}

function normalizePromptGroupIndex(value) {
  const index = Number(value);
  if (!Number.isInteger(index) || index < 0 || index >= promptGroups.length) {
    return 0;
  }
  return index;
}

function normalizePendingTripPlan(value) {
  if (!value || typeof value !== "object") {
    return null;
  }
  if (typeof value.userText !== "string" || typeof value.replyText !== "string") {
    return null;
  }
  return {
    userText: value.userText,
    replyText: value.replyText,
  };
}

async function sendPrompt(text) {
  if (isSending.value) return;
  message.value = text;
  await nextTick();
  await send();
}

function shufflePrompts() {
  promptGroupIndex.value = (promptGroupIndex.value + 1) % promptGroups.length;
  prompts.value = promptGroups[promptGroupIndex.value];
  promptExpanded.value = false;
}

function togglePrompts() {
  promptExpanded.value = !promptExpanded.value;
}

async function sendPendingQuestion() {
  if (hasAutoSent.value || !pendingAutoQuestion.value.trim()) return;
  hasAutoSent.value = true;
  message.value = pendingAutoQuestion.value.trim();
  await nextTick();
  await send();
}

async function send() {
  if (isSending.value) return;
  const content = message.value.trim();
  if (!content) {
    toast("请输入内容");
    return;
  }
  if (await handleTripConfirmation(content)) return;

  const token = uni.getStorageSync("travel_token");
  if (!token) {
    toast("请先登录后再使用旅行管家");
    uni.navigateTo({ url: "/pages/login/index" });
    return;
  }

  const userId = `msg-${Date.now()}`;
  const replyId = `${userId}-reply`;
  messages.value.push({ id: userId, role: "user", text: content, time: "现在" });
  messages.value.push({ id: replyId, role: "ai", text: "", time: "现在" });
  message.value = "";
  isSending.value = true;
  await nextTick();
  lastId.value = replyId;

  try {
    const response = await fetch(chatStreamUrl, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify({ message: content }),
    });

    const reply = messages.value.find((item) => item.id === replyId);
    if (!reply) return;

    if (!response.ok) {
      reply.text = "旅行管家服务暂不可用，请确认网关已启动并检查模型配置。";
      toast("旅行管家服务暂不可用");
      return;
    }

    if (!response.body) {
      reply.text = await response.text();
      rememberTripPlanIfNeeded(content, reply);
      return;
    }

    const reader = response.body.getReader();
    const decoder = new TextDecoder("utf-8");
    while (true) {
      const { done, value } = await reader.read();
      if (done) break;
      reply.text += decoder.decode(value, { stream: true });
      await nextTick();
      lastId.value = replyId;
    }
    reply.text += decoder.decode();
    rememberTripPlanIfNeeded(content, reply);
  } catch (error) {
    const reply = messages.value.find((item) => item.id === replyId);
    if (reply) reply.text = "旅行管家服务暂不可用，请确认网关已启动后再试。";
    toast("旅行管家服务暂不可用");
  } finally {
    isSending.value = false;
    await nextTick();
    lastId.value = replyId;
  }
}

function rememberTripPlanIfNeeded(userText, reply) {
  if (!reply || !shouldOfferTripCreation(userText, reply.text)) return;
  const cleanReply = stripCreateTripQuestion(reply.text);
  pendingTripPlan.value = { userText, replyText: cleanReply };
  reply.text = `${cleanReply}\n\n是否根据此计划生成行程？你可以回复“可以”“行”“好的”或“生成”。`;
}

async function handleTripConfirmation(content) {
  if (!pendingTripPlan.value || !isConfirmCreateTrip(content)) return false;
  const token = uni.getStorageSync("travel_token");
  if (!token) {
    toast("请先登录后再生成行程");
    uni.navigateTo({ url: "/pages/login/index" });
    return true;
  }

  const userId = `msg-${Date.now()}`;
  const replyId = `${userId}-reply`;
  messages.value.push({ id: userId, role: "user", text: content, time: "现在" });

  const tripPayload = buildTripDraft(pendingTripPlan.value.userText, pendingTripPlan.value.replyText);
  const savedTrip = await createTripOnServer(tripPayload, token);
  if (savedTrip) {
    syncSavedTrip(savedTrip);
    pendingTripPlan.value = null;
  }

  messages.value.push({
    id: replyId,
    role: "ai",
    text: savedTrip
      ? "已根据上面的计划生成行程安排，你可以去行程页继续查看。"
      : "这次还没有生成成功，请根据上面的提示调整后再试一次。",
    time: "现在",
  });
  await nextTick();
  lastId.value = replyId;
  return true;
}

function syncSavedTrip(savedTrip) {
  if (!savedTrip || typeof savedTrip !== "object") {
    return;
  }
  const currentTrips = tripPageState.getTrips();
  const nextTrips = [savedTrip, ...currentTrips.filter((item) => item?.id !== savedTrip.id)];
  tripPageState.saveTrips(nextTrips);
  tripPageState.selectTrip(savedTrip.id || "");
  tripPageState.saveSelectedTrip(savedTrip);
}

function isConfirmCreateTrip(content) {
  return /^(可以|行|好的|好|生成|确认|没问题|安排|行的|好呀|好啊)[。！! ]*$/.test(content.trim());
}

function stripCreateTripQuestion(text) {
  return text.replace(/\n*是否根据此计划生成行程？.*$/s, "").trim();
}

function buildTripDraft(userText, replyText) {
  const summary = extractTripSummary(userText);
  return {
    id: `trip-${Date.now()}`,
    title: buildTripTitle(userText, replyText),
    subtitle: buildTripSubtitle(replyText),
    destination: extractDestination(userText),
    dateRange: summary.date,
    dayCount: Number(String(summary.days || "").replace(/[^\d]/g, "")) || 0,
    people: summary.people,
    budgetLevel: summary.budget,
    sourceQuestion: userText,
    sourceReply: replyText,
    days: extractTripDays(replyText),
    summary,
    budget: extractBudgetItems(replyText),
    alerts: extractAlerts(replyText),
    createdAt: new Date().toISOString(),
  };
}

async function createTripOnServer(trip, token) {
  try {
    const response = await fetch(createTripUrl, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify({
        title: trip.title,
        subtitle: trip.subtitle,
        destination: trip.destination,
        date_range: trip.dateRange,
        day_count: trip.dayCount,
        people: trip.people,
        budget_level: trip.budgetLevel,
        source_question: trip.sourceQuestion,
        source_reply: trip.sourceReply,
        summary: trip.summary,
        days: trip.days,
        budget: trip.budget,
        alerts: trip.alerts,
      }),
    });
    const result = await response.json();
    if (!response.ok || result.code !== 0) {
      toast(result.msg || "行程保存失败");
      return null;
    }
    return result.data;
  } catch (error) {
    toast("行程保存失败");
    return null;
  }
}

function extractDestination(userText) {
  const destinationMatch = userText.match(/目的地\s*([^，。；\s]+)/) || userText.match(/去([\u4e00-\u9fa5]{2,8})/);
  return destinationMatch ? destinationMatch[1] : "";
}

function buildTripTitle(userText, replyText) {
  const destination = extractDestination(userText) || "生成行程";
  const dayMatch = `${userText}\n${replyText}`.match(/共\s*(\d+)\s*天|(\d+)\s*天/);
  const days = dayMatch ? `${dayMatch[1] || dayMatch[2]}天` : "旅行";
  return `${destination} / ${days}`;
}

function buildTripSubtitle(replyText) {
  const firstLine = replyText.split(/\n+/).map((item) => item.trim()).find(Boolean);
  return firstLine ? firstLine.slice(0, 32) : "由旅行管家生成";
}

function extractTripSummary(userText) {
  const dateMatch = userText.match(/出行日期([^，。\n]+)/);
  const peopleMatch = userText.match(/(\d+人)同行/);
  const budgetMatch = userText.match(/预算级别([^，。\n]+)/);
  const dayMatch = userText.match(/共\s*(\d+)\s*天|(\d+)\s*天/);
  return {
    date: dateMatch ? dateMatch[1] : "待定",
    days: dayMatch ? `${dayMatch[1] || dayMatch[2]}天` : "待定",
    people: peopleMatch ? peopleMatch[1] : "待定",
    budget: budgetMatch ? budgetMatch[1] : "待定",
  };
}

function extractTripDays(replyText) {
  const normalized = replyText.replace(/\r/g, "");
  const matches = [...normalized.matchAll(/第\s*([一二三四五六七八九十\d]+)\s*天[：:、\-\s]*(.*)/g)];
  if (!matches.length) {
    return [{
      day: 1,
      title: "生成行程",
      route: normalized.slice(0, 260),
      tips: [{ icon: "路", title: "行程重点", text: pickPlanSentence(normalized, /路线|景点|住宿|交通|预算/) }],
      weather: "出行前请以实时天气和平台信息为准。",
    }];
  }
  return matches.map((match, index) => {
    const start = match.index || 0;
    const end = index + 1 < matches.length ? matches[index + 1].index || normalized.length : normalized.length;
    const block = normalized.slice(start, end).trim();
    return {
      day: chineseDayToNumber(match[1]) || index + 1,
      title: (match[2] || "当日安排").trim().slice(0, 24),
      route: block.replace(/\n+/g, " ").slice(0, 180),
      tips: [
        { icon: "行", title: "交通餐饮", text: pickPlanSentence(block, /交通|地铁|打车|餐饮|美食|午餐|晚餐/) },
        { icon: "住", title: "住宿提醒", text: pickPlanSentence(block, /住宿|酒店|区域|预订|门票/) },
      ],
      weather: "出行前请以实时天气和平台信息为准。",
    };
  });
}

function pickPlanSentence(text, pattern) {
  const sentence = text.split(/[。；;\n]/).map((item) => item.trim()).find((item) => pattern.test(item));
  return sentence ? sentence.slice(0, 52) : "根据旅行计划灵活安排，并预留交通缓冲时间。";
}

function chineseDayToNumber(value) {
  if (/^\d+$/.test(value)) return Number(value);
  const map = { 一: 1, 二: 2, 三: 3, 四: 4, 五: 5, 六: 6, 七: 7, 八: 8, 九: 9, 十: 10 };
  return map[value] || 0;
}

function extractBudgetItems(replyText) {
  const moneyMatches = [...replyText.matchAll(/([^，。；;\n]{0,8})(?:约|预计)?\s*[¥￥]\s*([\d,]+)/g)];
  if (moneyMatches.length < 2) {
    return [];
  }
  return moneyMatches.slice(0, 3).map((match) => ({ label: normalizeBudgetLabel(match[1]), amount: `¥${match[2]}` }));
}

function normalizeBudgetLabel(text) {
  const label = String(text || "").replace(/[：:\s]/g, "");
  if (/酒店|住宿/.test(label)) return "酒店住宿";
  if (/交通|车票|机票|高铁/.test(label)) return "交通出行";
  if (/餐饮|美食|门票|景点/.test(label)) return "餐饮门票";
  return label || "预算项目";
}

function extractAlerts(replyText) {
  const lines = replyText.split(/\n|。|；/).map((item) => item.trim()).filter(Boolean);
  const alerts = lines.filter((item) => /注意|提醒|建议|避开|提前|预约|天气|防晒|下雨|缓冲/.test(item)).slice(0, 3);
  return alerts.length ? alerts : ["出行前请再次确认天气、营业时间和交通班次。", "热门景点和餐厅建议提前预约，并预留路上缓冲时间。"];
}

function formatAiBlocks(text) {
  const lines = String(text || "").split(/\n+/).map((item) => item.trim()).filter(Boolean);
  if (!lines.length) {
    return [{ type: "text", text: "" }];
  }
  return lines.map((line) => ({
    type: isSectionLine(line) ? "title" : "text",
    text: line,
  }));
}

function isSectionLine(text) {
  return /^第\s*[一二三四五六七八九十\d]+天/.test(text) || /^(路线|交通|住宿|预算|餐饮|美食|注意事项|提醒|建议)[：:]/.test(text);
}

function toast(title) {
  uni.showToast({ title, icon: "none" });
}
</script>

<style scoped>
.page {
  min-height: 100vh;
  background: #dfe9e6;
}

.phone-shell {
  width: 100%;
  max-width: 430px;
  min-height: 100vh;
  margin: 0 auto;
  padding: calc(env(safe-area-inset-top) + 12px) 18px calc(env(safe-area-inset-bottom) + 94px);
  box-sizing: border-box;
  background: linear-gradient(180deg, #fbfdfb 0%, #eef6f3 46%, #f8faf7 100%);
}

.chat-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 14px;
}

.header-main {
  display: flex;
  align-items: center;
  gap: 12px;
  flex: 1;
  min-width: 0;
}

.header-copy {
  min-width: 0;
}

.bot-avatar,
.msg-ai-icon {
  width: 40px;
  height: 40px;
  border-radius: 16px;
  background: #102033;
  color: #f8fbf9;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 15px;
  font-weight: 900;
  flex: 0 0 40px;
}

.bot-name,
.panel-title,
.tool-title,
.welcome-title,
.plan-title,
.ai-block-title {
  display: block;
  color: #102033;
  font-weight: 900;
}

.bot-name {
  font-size: 19px;
}

.bot-status,
.tool-desc,
.welcome-desc,
.plan-desc,
.panel-subtitle,
.msg-time,
.input-caption,
.ai-block-text {
  display: block;
  color: #667a78;
  font-size: 12px;
  line-height: 1.6;
}

.panel-more,
.panel-more,
.prompt-toggle {
  flex: 0 0 auto;
  min-height: 34px;
  padding: 0 12px;
  border-radius: 17px;
  background: rgba(255,255,252,.98);
  border: 1px solid rgba(224,232,229,.95);
  color: #0f5f75;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: 900;
}

.tools-scroll {
  margin-top: 18px;
  white-space: nowrap;
}

.quick-tools {
  display: inline-flex;
  gap: 10px;
  padding-right: 8px;
}

.tool-chip,
.prompt-panel,
.welcome-card,
.bubble.ai,
.plan-card-compact,
.input-bar {
  background: rgba(255,255,252,.98);
  border: 1px solid rgba(224,232,229,.92);
  box-shadow: 0 12px 28px rgba(16,32,51,.07);
}

.tool-chip {
  display: inline-flex;
  align-items: center;
  gap: 10px;
  width: 164px;
  min-height: 68px;
  padding: 12px;
  border-radius: 20px;
  box-sizing: border-box;
}

.tool-icon {
  width: 32px;
  height: 32px;
  border-radius: 12px;
  background: #edf7f3;
  color: #0f5f75;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: 900;
  flex: 0 0 32px;
}

.tool-copy {
  min-width: 0;
}

.tool-title {
  font-size: 14px;
}

.prompt-panel {
  margin-top: 16px;
  padding: 16px;
  border-radius: 24px;
}

.panel-head {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 12px;
}

.panel-title {
  font-size: 18px;
}

.prompt-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
  margin-top: 14px;
}

.prompt-list button {
  width: 100%;
  min-height: 42px;
  padding: 10px 14px;
  border-radius: 18px;
  background: #edf7f3;
  color: #0f5f75;
  font-size: 13px;
  line-height: 1.5;
  text-align: left;
}

.prompt-toggle {
  width: fit-content;
  margin-top: 12px;
}

.chat-scroll {
  height: calc(100vh - 402px);
  margin-top: 22px;
  padding-bottom: 14px;
}

.msg-wrap {
  display: flex;
  gap: 12px;
  margin-bottom: 22px;
}

.msg-wrap.user {
  justify-content: flex-end;
}

.msg-side {
  padding-top: 10px;
}

.msg-content {
  max-width: 84%;
}

.welcome-card {
  margin-bottom: 14px;
  padding: 18px;
  border-radius: 24px;
  background: linear-gradient(135deg, rgba(16,32,51,.98) 0%, rgba(26,56,76,.96) 100%);
  border-color: rgba(16,32,51,.08);
}

.welcome-kicker {
  display: inline-flex;
  width: fit-content;
  min-height: 24px;
  padding: 0 10px;
  border-radius: 999px;
  background: rgba(185,234,217,.18);
  color: #b9ead9;
  align-items: center;
  font-size: 11px;
  font-weight: 900;
}

.welcome-title {
  margin-top: 12px;
  color: #f8fbf9;
  font-size: 19px;
}

.welcome-desc {
  margin-top: 8px;
  color: rgba(248,251,249,.78);
}

.welcome-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-top: 14px;
}

.welcome-tags text,
.plan-tags text {
  padding: 5px 9px;
  border-radius: 999px;
  background: #edf7f3;
  color: #0f5f75;
  font-size: 11px;
  font-weight: 900;
}

.welcome-tags text {
  background: rgba(255,255,255,.08);
  color: #f8fbf9;
}

.bubble {
  padding: 16px 18px;
  border-radius: 22px;
  white-space: pre-wrap;
  line-height: 1.8;
  font-size: 14px;
}

.bubble.ai {
  color: #102033;
  border-top-left-radius: 10px;
}

.bubble.user {
  background: #102033;
  color: #f8fbf9;
  border-top-right-radius: 10px;
}

.ai-block + .ai-block {
  margin-top: 12px;
}

.ai-block.title {
  padding-top: 2px;
}

.ai-block-title {
  color: #0f5f75;
  font-size: 14px;
}

.ai-block-text {
  color: #344054;
  font-size: 14px;
  line-height: 1.8;
}

.msg-time {
  margin-top: 8px;
  padding: 0 6px;
}

.msg-time.ai {
  text-align: left;
}

.msg-time.user {
  text-align: right;
}

.plan-card-compact {
  margin-top: 12px;
  overflow: hidden;
  border-radius: 22px;
}

.plan-img {
  width: 100%;
  height: 128px;
}

.plan-body {
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding: 14px;
}

.plan-rating {
  color: #0f5f75;
  font-size: 12px;
  font-weight: 900;
}

.plan-title {
  font-size: 16px;
}

.plan-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}

.input-area {
  position: sticky;
  bottom: calc(env(safe-area-inset-bottom) + 8px);
  margin-top: 10px;
  padding-top: 12px;
  background: linear-gradient(180deg, rgba(248,250,247,0), rgba(248,250,247,.92) 36%, rgba(248,250,247,.98) 100%);
}

.input-caption {
  margin: 0 4px 8px;
}

.input-bar {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  border-radius: 24px;
}

.input-bar input {
  flex: 1;
  min-width: 0;
  min-height: 42px;
  color: #102033;
  font-size: 14px;
}

.send-btn {
  min-width: 78px;
  height: 42px;
  border-radius: 21px;
  background: #102033;
  color: #f8fbf9;
  font-size: 13px;
  font-weight: 900;
}
</style>





