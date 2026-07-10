<template>
  <view class="page">
    <view class="phone-shell with-nav">
      <view class="topbar">
        <view>
          <text class="greeting">你好</text>
          <text class="greeting-name">{{ greetingName || "欢迎回来" }}</text>
        </view>
        <image class="avatar" :src="avatarSrc" mode="aspectFill" @tap="goProfile" />
      </view>

      <view class="planner-card">
        <view class="planner-head">
          <view>
            <text class="planner-kicker">行程定制</text>
            <text class="planner-title">从你的想法开始规划</text>
          </view>
          <text class="planner-sub">填完信息后会自动带去助手页</text>
        </view>

        <view class="field-grid">
          <view class="field full">
            <text class="field-label">目的地</text>
            <input class="field-input" v-model="destinationName" placeholder="请输入想去的城市或地区" />
          </view>

          <view class="field full">
            <view class="field-label-row">
              <text class="field-label">出行日期</text>
              <text class="field-tip">{{ tripDaysText }}</text>
            </view>
            <view class="date-range-box">
              <view class="date-field">
                <input class="date-input" v-model="startDateText" placeholder="7月1日" @blur="normalizeStartDate" />
                <picker mode="date" :value="pickerStartValue" fields="day" @change="onStartDatePick">
                  <text class="date-picker-text">选</text>
                </picker>
              </view>
              <text class="date-sep">到</text>
              <view class="date-field">
                <input class="date-input" v-model="endDateText" placeholder="7月3日" @blur="normalizeEndDate" />
                <picker mode="date" :value="pickerEndValue" fields="day" @change="onEndDatePick">
                  <text class="date-picker-text">选</text>
                </picker>
              </view>
            </view>
          </view>

          <view class="field">
            <text class="field-label">同行人数</text>
            <view class="stepper-row">
              <text class="stepper-btn" @tap="setPeople(people - 1)">-</text>
              <text class="field-value">{{ peopleText }}</text>
              <text class="stepper-btn" @tap="setPeople(people + 1)">+</text>
            </view>
          </view>

          <view class="field">
            <text class="field-label">预算级别</text>
            <view class="chip-row compact">
              <text v-for="item in budgetOptions" :key="item" :class="['chip', budget === item ? 'active' : '']" @tap="budget = item">{{ item }}</text>
            </view>
          </view>

          <view class="field full">
            <text class="field-label">兴趣偏好</text>
            <view class="chip-row">
              <text v-for="item in interestOptions" :key="item" :class="['chip', selectedInterests.includes(item) ? 'active' : '']" @tap="toggleInterest(item)">{{ item }}</text>
            </view>
          </view>
        </view>

        <button class="generate-btn" @tap="goAI">生成我的专属行程</button>
      </view>

      <view class="section-row">
        <view>
          <text class="section-eyebrow">最近安排</text>
          <text class="section-title">最近保存的行程</text>
        </view>
        <text class="section-more" @tap="goTrip">查看全部</text>
      </view>

      <view v-if="tripView.hasTrip" class="trip-card" @tap="goTrip">
        <text class="trip-title">{{ tripView.title }}</text>
        <text class="trip-subtitle">{{ tripView.subtitle }}</text>
        <view class="summary-grid">
          <view class="summary-item">
            <text class="summary-label">出发时间</text>
            <text class="summary-value">{{ tripView.summary.date }}</text>
          </view>
          <view class="summary-item">
            <text class="summary-label">游玩天数</text>
            <text class="summary-value">{{ tripView.summary.days }}</text>
          </view>
          <view class="summary-item">
            <text class="summary-label">同行人数</text>
            <text class="summary-value">{{ tripView.summary.people }}</text>
          </view>
        </view>
        <view v-if="tripView.days.length" class="day-preview">
          <text class="day-preview-title">第{{ tripView.days[0].day }}天</text>
          <text class="day-preview-text">{{ tripView.days[0].route }}</text>
        </view>
      </view>

      <view v-else class="empty-card">
        <text class="empty-title">还没有你的行程安排</text>
        <text class="empty-desc">先填好目的地、日期、人数和预算，再让旅行助手帮你整理一份计划。</text>
      </view>

      <view class="section-row">
        <view>
          <text class="section-eyebrow">出行习惯</text>
          <text class="section-title">旅行偏好</text>
        </view>
      </view>

      <view v-if="selectedInterests.length" class="pref-card">
        <text v-for="item in selectedInterests" :key="item" class="pref-chip">{{ item }}</text>
      </view>
      <view v-else class="empty-inline">当前还没有保存旅行偏好。</view>

      <BottomNav active="explore" />
    </view>
  </view>
</template>

<script setup>
import { computed, onMounted, ref, watch } from "vue";
import { onShow } from "@dcloudio/uni-app";
import BottomNav from "../../components/BottomNav.vue";
import { apiBase, getToken, requestJSON, setStoredUser, storedUserRef, syncStoredUser } from "../../utils/api.js";
import { createTripPageState } from "../../utils/tripPageState.js";
import { normalizeTripView } from "../../utils/tripViewModel.js";

const avatar = "/static/avatar-traveler.svg";
const plannerState = createTripPageState();
const greetingName = computed(() => storedUserRef.value?.nickname || "");
const destinationName = ref("");
const startDateText = ref("");
const endDateText = ref("");
const people = ref(2);
const budget = ref("舒适型");
const budgetOptions = ["经济型", "舒适型", "轻奢型"];
const interestOptions = ["美食", "海岛", "城市漫步", "博物馆", "亲子", "徒步", "摄影", "温泉"];
const selectedInterests = ref([]);
const latestTrip = ref(null);
const avatarSrc = computed(() => storedUserRef.value?.avatarUrl || avatar);
const tripView = computed(() => normalizeTripView(latestTrip.value));
const startDate = computed(() => parseMonthDay(startDateText.value));
const endDate = computed(() => parseMonthDay(endDateText.value));
const pickerStartValue = computed(() => toPickerValue(startDate.value));
const pickerEndValue = computed(() => toPickerValue(endDate.value));
const dayCount = computed(() => calculateTripDays(startDate.value, endDate.value));
const tripDaysText = computed(() => `${dayCount.value}天`);
const peopleText = computed(() => `${people.value}人`);
const normalizedDateRange = computed(() => `${formatMonthDay(startDate.value)}到${formatMonthDay(endDate.value)}`);

onMounted(() => {
  uni.hideTabBar && uni.hideTabBar();
  restorePlannerDraft();
  loadDashboard();
  loadLatestTrip();
});

onShow(() => {
  syncStoredUser();
  loadDashboard();
  loadLatestTrip();
});

watch([destinationName, startDateText, endDateText, people, budget, selectedInterests], () => {
  plannerState.savePlannerDraft({
    destinationName: destinationName.value,
    startDateText: startDateText.value,
    endDateText: endDateText.value,
    people: people.value,
    budget: budget.value,
    selectedInterests: [...selectedInterests.value],
  });
}, { deep: true });

async function loadDashboard() {
  try {
    const data = await requestJSON("/api/v1/user/dashboard");
    if (data?.user) {
      greetingName.value = data.user.nickname || "";
      setStoredUser(data.user);
    }
    if (!selectedInterests.value.length) {
      selectedInterests.value = Array.isArray(data?.preferences?.items) ? data.preferences.items : [];
    }
  } catch (error) {
    console.warn("读取首页数据失败", error);
  }
}

async function loadLatestTrip() {
  const token = getToken();
  if (!token) {
    latestTrip.value = null;
    return;
  }

  const cachedTrips = plannerState.getTrips();
  const cachedSelectedTrip = plannerState.getSelectedTrip();
  if (cachedSelectedTrip && typeof cachedSelectedTrip === "object") {
    latestTrip.value = cachedSelectedTrip;
  } else if (Array.isArray(cachedTrips) && cachedTrips.length) {
    latestTrip.value = cachedTrips[0];
  }

  try {
    const response = await fetch(`${apiBase}/api/v1/trips/latest`, {
      headers: { Authorization: `Bearer ${token}` },
    });
    const result = await response.json();
    if (response.ok && result?.code === 0 && result.data) {
      latestTrip.value = result.data;
      plannerState.saveTrips([result.data]);
      plannerState.saveSelectedTrip(result.data);
    }
  } catch (error) {
    console.warn("读取最新行程失败", error);
  }
}

function restorePlannerDraft() {
  const draft = plannerState.getPlannerDraft();
  if (!draft || typeof draft !== "object") {
    return;
  }
  destinationName.value = draft.destinationName || "";
  startDateText.value = draft.startDateText || "";
  endDateText.value = draft.endDateText || "";
  people.value = Number(draft.people) || 2;
  budget.value = draft.budget || "舒适型";
  selectedInterests.value = Array.isArray(draft.selectedInterests) ? draft.selectedInterests : [];
}

function goAI() {
  const destination = destinationName.value.trim();
  if (!destination) {
    toast("请输入目的地");
    return;
  }
  if (!startDate.value || !endDate.value) {
    toast("请输入完整出行日期");
    return;
  }

  const interests = selectedInterests.value.length ? selectedInterests.value.join("、") : "无特别偏好";
  plannerState.savePlannerDraft({
    destinationName: destinationName.value,
    startDateText: startDateText.value,
    endDateText: endDateText.value,
    people: people.value,
    budget: budget.value,
    selectedInterests: [...selectedInterests.value],
  });

  const question = `请根据以下信息为我生成专属旅行行程：目的地${destination}，出行日期${normalizedDateRange.value}，共${dayCount.value}天，${peopleText.value}同行，预算级别${budget.value}，兴趣偏好${interests}。请给出每天路线、餐饮建议、住宿区域、交通方式、预算拆分和注意事项。`;
  uni.navigateTo({ url: `/pages/ai/index?autoSend=1&question=${encodeURIComponent(question)}` });
}

function goTrip() {
  uni.switchTab({ url: "/pages/trip/index" });
}

function goProfile() {
  uni.switchTab({ url: "/pages/profile/index" });
}

function setPeople(value) {
  people.value = Math.min(12, Math.max(1, value));
}

function toggleInterest(item) {
  selectedInterests.value = selectedInterests.value.includes(item)
    ? selectedInterests.value.filter((value) => value !== item)
    : [...selectedInterests.value, item];
}

function normalizeStartDate() {
  startDateText.value = formatMonthDay(startDate.value);
}

function normalizeEndDate() {
  endDateText.value = formatMonthDay(endDate.value);
}

function onStartDatePick(event) {
  startDateText.value = formatPickerDate(event.detail.value);
}

function onEndDatePick(event) {
  endDateText.value = formatPickerDate(event.detail.value);
}

function toast(title) {
  uni.showToast({ title, icon: "none" });
}

function parseMonthDay(text) {
  const match = String(text || "").match(/(\d{1,2})\s*月\s*(\d{1,2})\s*(?:日|号)?/);
  if (!match) {
    return null;
  }
  const month = Number(match[1]);
  const day = Number(match[2]);
  const date = new Date(new Date().getFullYear(), month - 1, day);
  if (date.getMonth() !== month - 1 || date.getDate() !== day) {
    return null;
  }
  return date;
}

function calculateTripDays(start, end) {
  if (!start || !end) {
    return 1;
  }
  const adjustedEnd = new Date(end);
  if (adjustedEnd < start) {
    adjustedEnd.setFullYear(adjustedEnd.getFullYear() + 1);
  }
  return Math.max(1, Math.round((adjustedEnd - start) / (24 * 60 * 60 * 1000)) + 1);
}

function formatMonthDay(date) {
  if (!date) {
    return "";
  }
  return `${date.getMonth() + 1}月${date.getDate()}日`;
}

function toPickerValue(date) {
  const value = date || new Date();
  return `${value.getFullYear()}-${pad2(value.getMonth() + 1)}-${pad2(value.getDate())}`;
}

function formatPickerDate(value) {
  const parts = String(value || "").split("-");
  if (parts.length !== 3) {
    return "";
  }
  return `${Number(parts[1])}月${Number(parts[2])}日`;
}

function pad2(value) {
  return String(value).padStart(2, "0");
}
</script>

<style scoped>
.page { min-height: 100vh; background: #dfe9e6; }
.phone-shell { width: 100%; max-width: 430px; min-height: 100vh; margin: 0 auto; padding: calc(env(safe-area-inset-top) + 10px) 18px 0; box-sizing: border-box; background: linear-gradient(180deg, #fbfdfb 0%, #eef6f3 46%, #f8faf7 100%); }
.with-nav { padding-bottom: calc(env(safe-area-inset-bottom) + 92px); }
.topbar { min-height: 58px; display: flex; align-items: center; justify-content: space-between; }
.greeting { display: block; color: #667a78; font-size: 12px; font-weight: 800; }
.greeting-name { display: block; margin-top: 4px; color: #102033; font-size: 26px; font-weight: 900; }
.avatar { width: 40px; height: 40px; border-radius: 15px; background: #fffdfa; box-shadow: 0 8px 22px rgba(16,32,51,.08); }
.planner-card { margin-top: 16px; padding: 16px; border-radius: 24px; background: #102033; box-shadow: 0 18px 36px rgba(16,32,51,.18); }
.planner-kicker { display: block; color: #b9ead9; font-size: 11px; font-weight: 900; }
.planner-title { display: block; margin-top: 4px; color: #f8fbf9; font-size: 22px; font-weight: 900; }
.planner-sub { color: rgba(248,251,249,.72); font-size: 11px; font-weight: 800; }
.planner-head { display: flex; align-items: flex-end; justify-content: space-between; gap: 12px; }
.field-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 10px; margin-top: 14px; }
.field { min-height: 74px; padding: 12px; border-radius: 18px; background: #f8fbf9; border: 1px solid rgba(224,232,229,.95); }
.field.full { grid-column: 1 / 3; }
.field-label { display: block; color: #6b7c8e; font-size: 12px; font-weight: 900; }
.field-label-row { display: flex; align-items: center; justify-content: space-between; gap: 8px; }
.field-tip { color: #0f5f75; font-size: 11px; font-weight: 900; }
.field-input { width: 100%; min-height: 34px; margin-top: 6px; color: #102033; font-size: 18px; font-weight: 900; }
.field-value { color: #102033; font-size: 17px; font-weight: 900; }
.date-range-box { display: grid; grid-template-columns: minmax(0, 1fr) auto minmax(0, 1fr); align-items: center; gap: 6px; margin-top: 8px; }
.date-field { min-width: 0; height: 36px; padding: 0 6px; border-radius: 13px; background: #fff; border: 1px solid #dfe8ef; display: flex; align-items: center; gap: 4px; }
.date-input { flex: 1; min-width: 0; height: 32px; color: #102033; font-size: 14px; font-weight: 900; }
.date-picker-text { width: 23px; height: 23px; border-radius: 999px; background: #edf7f3; color: #0f5f75; display: inline-flex; align-items: center; justify-content: center; font-size: 10px; font-weight: 900; }
.date-sep { color: #607386; font-size: 12px; font-weight: 800; }
.stepper-row { display: flex; align-items: center; justify-content: space-between; gap: 8px; margin-top: 8px; }
.stepper-btn { width: 28px; height: 28px; border-radius: 999px; background: #edf7f3; color: #0f5f75; display: flex; align-items: center; justify-content: center; font-size: 18px; font-weight: 900; }
.chip-row { display: flex; flex-wrap: wrap; gap: 6px; margin-top: 8px; }
.chip-row.compact { gap: 5px; }
.chip { min-height: 28px; padding: 0 10px; border-radius: 999px; background: #edf2f7; color: #607386; display: inline-flex; align-items: center; justify-content: center; font-size: 11px; font-weight: 900; }
.chip.active { background: #e0f4eb; color: #0f5f75; }
.generate-btn { height: 50px; margin-top: 14px; border-radius: 18px; background: #b9ead9; color: #102033; font-size: 15px; font-weight: 900; box-shadow: 0 14px 26px rgba(66,154,125,.2); }
.generate-btn::after { border: 0; }
.section-row { display: flex; align-items: flex-end; justify-content: space-between; margin: 24px 2px 12px; }
.section-eyebrow { display: block; color: #08766c; font-size: 11px; font-weight: 900; }
.section-title { display: block; margin-top: 4px; color: #102033; font-size: 20px; font-weight: 900; }
.section-more { color: #0f5f75; font-size: 13px; font-weight: 900; }
.trip-card, .pref-card, .empty-card, .empty-inline { background: rgba(255,255,252,.98); border: 1px solid rgba(224,232,229,.9); box-shadow: 0 10px 24px rgba(16,32,51,.07); }
.trip-card { padding: 16px; border-radius: 22px; }
.trip-title { display: block; color: #102033; font-size: 19px; font-weight: 900; }
.trip-subtitle { display: block; margin-top: 6px; color: #667a78; font-size: 13px; line-height: 1.6; }
.summary-grid { display: grid; grid-template-columns: repeat(3, 1fr); gap: 8px; margin-top: 14px; }
.summary-item { padding: 10px; border-radius: 14px; background: #f8fbf9; }
.summary-label { display: block; color: #8a9894; font-size: 11px; }
.summary-value { display: block; margin-top: 5px; color: #102033; font-size: 13px; font-weight: 900; }
.day-preview { margin-top: 14px; padding: 12px; border-radius: 16px; background: #edf7f3; }
.day-preview-title { display: block; color: #0f5f75; font-size: 13px; font-weight: 900; }
.day-preview-text { display: block; margin-top: 6px; color: #334155; font-size: 13px; line-height: 1.6; }
.empty-card { padding: 18px; border-radius: 22px; }
.empty-title { display: block; color: #102033; font-size: 17px; font-weight: 900; }
.empty-desc { display: block; margin-top: 8px; color: #667a78; font-size: 13px; line-height: 1.7; }
.pref-card { display: flex; flex-wrap: wrap; gap: 8px; padding: 14px; border-radius: 22px; }
.pref-chip { min-height: 30px; padding: 0 12px; border-radius: 999px; background: #edf7f3; color: #0f5f75; display: inline-flex; align-items: center; justify-content: center; font-size: 12px; font-weight: 900; }
.empty-inline { padding: 16px; border-radius: 18px; color: #667a78; font-size: 13px; line-height: 1.6; }
</style>
