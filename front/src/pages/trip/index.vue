<template>
  <view class="page">
    <view class="phone-shell with-nav">
      <view class="topbar">
        <text class="title">我的行程</text>
        <view class="top-actions">
          <image class="avatar" :src="avatar" mode="aspectFill" />
        </view>
      </view>

      <view v-if="tripList.length" class="trip-list-card">
        <view class="list-head">
          <text class="list-title">全部行程</text>
          <text class="list-meta">{{ tripList.length }}条</text>
        </view>
        <scroll-view class="trip-scroll" scroll-x show-scrollbar="false">
          <view class="trip-pill-row">
            <view
              v-for="item in tripList"
              :key="item.id"
              :class="['trip-pill', selectedTripId === item.id ? 'active' : '']"
              @tap="selectTrip(item.id)"
            >
              <view class="trip-pill-copy">
                <text class="trip-pill-title">{{ formatTripTitle(item) }}</text>
                <text v-if="formatTripMeta(item)" class="trip-pill-sub">{{ formatTripMeta(item) }}</text>
              </view>
              <text class="trip-pill-delete" @tap.stop="removeTrip(item.id)">{{ deletingTripId === item.id ? '...' : '删' }}</text>
            </view>
          </view>
        </scroll-view>
      </view>

      <view v-if="hasTrip" class="trip-content">
        <view class="dest-card">
          <image :src="destImage" mode="aspectFill" />
          <view class="dest-overlay"></view>
          <view class="dest-content">
            <text class="dest-name">{{ tripTitle }}</text>
            <text v-if="tripSubtitle" class="dest-sub">{{ tripSubtitle }}</text>
          </view>
        </view>

        <view v-if="summaryItems.length" class="summary-bar">
          <template v-for="(item, index) in summaryItems" :key="item.label">
            <view class="summary-item">
              <text class="summary-label">{{ item.label }}</text>
              <text class="summary-value">{{ item.value }}</text>
            </view>
            <view v-if="index < summaryItems.length - 1" class="summary-divider"></view>
          </template>
        </view>

        <text class="section-title">每日计划</text>
        <view v-if="displayDays.length" class="day-tabs">
          <text v-for="day in displayDays" :key="day.day" :class="['day-tab', activeDay === day.day ? 'active' : '']" @tap="activeDay = day.day">第{{ day.day }}天</text>
        </view>

        <view v-if="currentDay" class="day-card">
          <view class="day-header">
            <view class="day-num">{{ currentDay.day }}</view>
            <text class="day-title-text">{{ currentDayTitle }}</text>
          </view>
          <view class="day-body">
            <view v-if="currentDay.route" class="route-card">
              <text class="route-label">路线方案</text>
              <text class="route-text">{{ currentDay.route }}</text>
            </view>
            <view v-if="currentDay.tips.length" class="tips-grid">
              <view v-for="tip in currentDay.tips" :key="tip.title + tip.text" class="tip-card">
                <text class="tip-title">{{ tip.icon }} {{ tip.title }}</text>
                <text class="tip-text">{{ tip.text }}</text>
              </view>
            </view>
            <view v-if="currentDay.weather" class="weather-card">天气提醒 · {{ currentDay.weather }}</view>
          </view>
        </view>
        <view v-else class="inline-empty">暂无每日计划，保存完整行程后会自动展示。</view>

        <text class="section-title">预算构成</text>
        <view class="budget-card">
          <template v-if="budgetItems.length">
            <template v-for="(item, index) in budgetItems" :key="item.label || index">
              <view class="budget-row">
                <text>{{ item.label || '未命名项目' }}</text>
                <strong :class="index % 2 === 1 ? 'green' : ''">{{ item.amount }}</strong>
              </view>
              <view class="budget-bar">
                <view :class="['budget-fill', index % 2 === 1 ? 'green-bg' : 'blue']" :style="{ width: item.percent + '%' }"></view>
              </view>
            </template>
          </template>
          <view v-else class="budget-empty">当前行程还没有预算明细。</view>
          <view v-if="budgetTotal" class="budget-total">
            <text>预计总花费</text>
            <strong>{{ budgetTotal }}</strong>
          </view>
        </view>

        <template v-if="tripAlerts.length">
          <text class="section-title">出行提醒</text>
          <view class="alert-card">
            <text v-for="item in tripAlerts" :key="item">· {{ item }}</text>
          </view>
        </template>

        <view class="save-btn" @tap="saveTrip">{{ selectedTrip ? '当前行程已保存' : '保存行程' }}</view>
      </view>

      <view v-else class="empty-state">
        <view class="empty-visual">行</view>
        <text class="empty-title">还没有保存行程</text>
        <text class="empty-desc">向旅行管家说明目的地、日期、同行人数和预算，确认生成后即可查看行程预览、每日计划和预算构成。</text>
        <view class="empty-btn" @tap="goToAI">去问旅行管家</view>
      </view>

      <BottomNav active="trip" />
    </view>
  </view>
</template>

<script setup>
import { computed, onMounted, ref, watch } from "vue";
import { onShow } from "@dcloudio/uni-app";
import BottomNav from "../../components/BottomNav.vue";
import { requestJSON, storedUserRef, syncStoredUser } from "../../utils/api.js";
import { normalizeTripView } from "../../utils/tripViewModel.js";
import { createTripPageState } from "../../utils/tripPageState.js";

const avatar = computed(() => storedUserRef.value?.avatarUrl || "/static/avatar-traveler.svg");
const destImage = "https://images.unsplash.com/photo-1519181245277-cffeb31da2e3?auto=format&fit=crop&w=760&q=80";
const pageState = createTripPageState();
const activeDay = ref(1);
const tripList = ref([]);
const selectedTripId = ref("");
const selectedTrip = ref(null);
const deletingTripId = ref("");
const tripView = computed(() => normalizeTripView(selectedTrip.value));
const hasTrip = computed(() => tripView.value.hasTrip);
const displayDays = computed(() => tripView.value.days);
const tripTitle = computed(() => tripView.value.title || "未命名行程");
const tripSubtitle = computed(() => tripView.value.subtitle);
const tripSummary = computed(() => tripView.value.summary);
const summaryItems = computed(() => [
  { label: "出发", value: tripSummary.value.date },
  { label: "天数", value: tripSummary.value.days },
  { label: "同行", value: tripSummary.value.people },
].filter((item) => item.value));
const currentDay = computed(() => displayDays.value.find((item) => item.day === activeDay.value) || displayDays.value[0] || null);
const currentDayTitle = computed(() => {
  if (!currentDay.value) {
    return "";
  }
  return currentDay.value.title ? `第 ${currentDay.value.day} 天：${currentDay.value.title}` : `第 ${currentDay.value.day} 天`;
});
const budgetItems = computed(() => tripView.value.budgetItems);
const budgetTotal = computed(() => tripView.value.budgetTotal);
const tripAlerts = computed(() => tripView.value.alerts);

onMounted(() => {
  uni.hideTabBar && uni.hideTabBar();
  restoreTripState();
});

onShow(() => {
  syncStoredUser();
  loadTripData();
});

watch(displayDays, (items) => {
  if (!items.some((item) => item.day === activeDay.value)) {
    activeDay.value = items[0]?.day || 1;
  }
});

function restoreTripState() {
  const cachedTrips = pageState.getTrips();
  const cachedSelectedTripId = pageState.getSelectedTripId();
  const cachedSelectedTrip = pageState.getSelectedTrip();
  tripList.value = Array.isArray(cachedTrips) ? cachedTrips : [];
  selectedTripId.value = cachedSelectedTripId || tripList.value[0]?.id || "";
  selectedTrip.value = cachedSelectedTrip && typeof cachedSelectedTrip === "object" ? cachedSelectedTrip : null;
}

async function loadTripData() {
  const token = uni.getStorageSync("travel_token");
  if (!token) {
    tripList.value = [];
    selectedTripId.value = "";
    selectedTrip.value = null;
    pageState.clearTripState();
    return;
  }

  try {
    const list = await requestJSON("/api/v1/trips");
    tripList.value = Array.isArray(list) ? list : [];
    pageState.saveTrips(tripList.value);

    if (!tripList.value.length) {
      selectedTripId.value = "";
      selectedTrip.value = null;
      pageState.clearTripState();
      return;
    }

    const matchedTrip = tripList.value.find((item) => item.id === selectedTripId.value);
    const nextTripId = matchedTrip?.id || tripList.value[0]?.id || "";
    await selectTrip(nextTripId, { silent: true });
  } catch (error) {
    console.warn("读取行程列表失败", error);
    toast(error.message || "读取行程失败");
  }
}

async function selectTrip(tripId, options = {}) {
  if (!tripId) {
    selectedTripId.value = "";
    selectedTrip.value = null;
    pageState.selectTrip("");
    pageState.saveSelectedTrip(null);
    return;
  }
  selectedTripId.value = tripId;
  pageState.selectTrip(tripId);

  const cachedSelectedTrip = pageState.getSelectedTrip();
  if (cachedSelectedTrip?.id === tripId) {
    selectedTrip.value = cachedSelectedTrip;
  }

  try {
    const detail = await requestJSON(`/api/v1/trips/${tripId}`);
    selectedTrip.value = detail;
    pageState.saveSelectedTrip(detail);
    if (!options.silent) {
      toast("已切换行程");
    }
  } catch (error) {
    console.warn("读取行程详情失败", error);
    if (!options.silent) {
      toast(error.message || "读取行程详情失败");
    }
  }
}

async function removeTrip(tripId) {
  if (!tripId || deletingTripId.value) {
    return;
  }
  deletingTripId.value = tripId;
  try {
    await requestJSON(`/api/v1/trips/${tripId}`, { method: "DELETE" });
    const nextTrips = tripList.value.filter((item) => item.id !== tripId);
    tripList.value = nextTrips;
    pageState.saveTrips(nextTrips);

    if (!nextTrips.length) {
      selectedTripId.value = "";
      selectedTrip.value = null;
      pageState.clearTripState();
      toast("已删除行程");
      return;
    }

    if (selectedTripId.value === tripId) {
      const nextTripId = nextTrips[0]?.id || "";
      await selectTrip(nextTripId, { silent: true });
    }
    toast("已删除行程");
  } catch (error) {
    console.warn("删除行程失败", error);
    toast(error.message || "删除行程失败");
  } finally {
    deletingTripId.value = "";
  }
}

function saveTrip() {
  toast(selectedTrip.value ? "当前行程已保存在列表中" : "请先生成并确认一份行程");
}

function formatTripTitle(item) {
  return item?.title || item?.destination || item?.date_range || item?.dateRange || "未命名行程";
}

function formatTripMeta(item) {
  const parts = [item?.date_range || item?.dateRange, item?.people, item?.budget_level || item?.budgetLevel].filter(Boolean);
  return parts.join(" / ");
}

function toast(title) { uni.showToast({ title, icon: "none" }); }
function goToAI() { uni.navigateTo({ url: "/pages/ai/index" }); }
</script>

<style scoped>
.page { min-height: 100vh; background: #eef2f7; }
.phone-shell { width: 100%; max-width: 430px; min-height: 100vh; margin: 0 auto; padding: calc(env(safe-area-inset-top) + 8px) 16px 0; box-sizing: border-box; background: #f6f8fc; }
.with-nav { padding-bottom: calc(env(safe-area-inset-bottom) + 132px); }
.topbar { display: flex; align-items: center; justify-content: space-between; min-height: 56px; }
.title { color: #101828; font-size: 22px; font-weight: 800; }
.top-actions { display: flex; align-items: center; }
.avatar { width: 34px; height: 34px; border-radius: 999px; border: 2px solid #fff; box-shadow: 0 2px 8px rgba(31,45,71,.1); }
.trip-list-card { margin-top: 8px; padding: 14px; border-radius: 22px; background: rgba(255,255,252,.98); border: 1px solid rgba(224,232,229,.9); box-shadow: 0 12px 28px rgba(16,32,51,.07); }
.list-head { display: flex; align-items: center; justify-content: space-between; gap: 12px; }
.list-title { color: #102033; font-size: 16px; font-weight: 900; }
.list-meta { color: #667a78; font-size: 12px; }
.trip-scroll { margin-top: 12px; white-space: nowrap; }
.trip-pill-row { display: inline-flex; gap: 10px; padding-right: 6px; }
.trip-pill { width: 210px; min-height: 74px; padding: 12px; border-radius: 18px; background: #f8fbf9; border: 1px solid rgba(224,232,229,.95); display: flex; align-items: flex-start; justify-content: space-between; gap: 10px; box-sizing: border-box; }
.trip-pill.active { background: #102033; border-color: rgba(16,32,51,.08); }
.trip-pill-copy { min-width: 0; flex: 1; }
.trip-pill-title { display: block; color: #102033; font-size: 14px; font-weight: 900; line-height: 1.45; }
.trip-pill-sub { display: block; margin-top: 6px; color: #667a78; font-size: 11px; line-height: 1.5; }
.trip-pill.active .trip-pill-title,
.trip-pill.active .trip-pill-sub { color: #f8fbf9; }
.trip-pill-delete { width: 26px; height: 26px; border-radius: 999px; background: #edf7f3; color: #0f5f75; display: flex; align-items: center; justify-content: center; font-size: 11px; font-weight: 900; flex: 0 0 26px; }
.trip-content { margin-top: 14px; }
.dest-card { position: relative; height: 168px; margin-top: 4px; border-radius: 24px; overflow: hidden; box-shadow: 18px 36px rgba(16,32,51,.16); }
.dest-card image { width: 100%; height: 100%; }
.dest-overlay { position: absolute; inset: 0; background: linear-gradient(180deg, rgba(8,23,35,.06), rgba(8,23,35,.72)); }
.dest-content { position: absolute; left: 18px; right: 18px; bottom: 18px; color: #fff; }
.dest-name { display: block; font-size: 27px; font-weight: 900; }
.dest-sub { display: block; margin-top: 4px; font-size: 13px; opacity: .9; }
.summary-bar { margin-top: 12px; padding: 14px 8px; border-radius: 18px; background: rgba(255,255,252,.98); display: flex; align-items: center; border: 1px solid rgba(224,232,229,.9); box-shadow: 0 12px 28px rgba(16,32,51,.07); }
.summary-item { flex: 1; text-align: center; }
.summary-label { display: block; color: #667a78; font-size: 12px; }
.summary-value { display: block; margin-top: 4px; color: #102033; font-size: 14px; font-weight: 700; }
.summary-divider { width: 1px; height: 28px; background: #dce8e4; }
.section-title { display: block; margin: 22px 0 12px; color: #102033; font-size: 19px; font-weight: 800; }
.day-tabs { display: flex; gap: 8px; margin-bottom: 14px; overflow-x: auto; padding-bottom: 2px; }
.day-tab { min-width: 78px; height: 38px; padding: 0 12px; border-radius: 19px; background: rgba(255,255,252,.96); color: #667a78; display: flex; align-items: center; justify-content: center; font-size: 13px; font-weight: 700; border: 1px solid rgba(224,232,229,.9); box-sizing: border-box; }
.day-tab.active { background: #102033; color: #f8fbf9; font-weight: 800; }
.day-card, .budget-card { background: rgba(255,255,252,.98); border: 1px solid rgba(224,232,229,.9); box-shadow: 0 12px 28px rgba(16,32,51,.07); }
.day-card { padding: 16px; border-radius: 20px; }
.day-header { display: flex; align-items: center; gap: 10px; margin-bottom: 16px; }
.day-num { width: 26px; height: 26px; border-radius: 999px; background: #102033; color: #fff; display: flex; align-items: center; justify-content: center; font-size: 13px; font-weight: 900; flex: 0 0 26px; }
.day-title-text { color: #102033; font-size: 17px; font-weight: 800; }
.day-body { display: flex; flex-direction: column; gap: 12px; }
.route-card { padding: 14px; border-radius: 14px; background: #edf7f3; border: 1px solid #dbece5; }
.route-label { display: block; color: #0f5f75; font-size: 14px; font-weight: 800; }
.route-text { display: block; margin-top: 6px; color: #344054; font-size: 14px; line-height: 1.55; }
.tips-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 10px; }
.tip-card { padding: 12px; border-radius: 14px; background: #fffdfa; border: 1px solid rgba(224,232,229,.92); }
.tip-title { display: block; color: #0f5f75; font-size: 13px; font-weight: 800; }
.tip-text { display: block; margin-top: 6px; color: #344054; font-size: 12px; line-height: 1.45; }
.weather-card { padding: 14px; border-radius: 14px; background: #f7f0df; color: #6d4b13; border: 1px solid #ead7a8; font-size: 13px; line-height: 1.5; }
.budget-card { padding: 16px; border-radius: 20px; }
.budget-row { display: flex; align-items: center; justify-content: space-between; color: #344054; font-size: 14px; }
.budget-row strong { color: #0f5f75; font-size: 16px; }
.budget-row .green { color: #0f5f75; }
.budget-bar { height: 8px; margin: 10px 0 16px; border-radius: 999px; background: #e1ebe7; }
.budget-fill { height: 100%; border-radius: 999px; }
.blue { background: #0f5f75; }
.green-bg { background: #caa459; }
.budget-total { display: flex; align-items: center; justify-content: space-between; padding-top: 16px; border-top: 1px solid #eaecf0; font-size: 16px; }
.budget-total strong { color: #0f5f75; font-size: 24px; font-weight: 800; }
.alert-card { padding: 16px; border-radius: 18px; background: #fffdfa; border: 1px solid #ead7a8; color: #6d4b13; font-size: 13px; line-height: 1.7; display: flex; flex-direction: column; gap: 8px; }
.inline-empty, .budget-empty { padding: 16px; border-radius: 18px; background: rgba(255,255,252,.98); border: 1px solid rgba(224,232,229,.92); color: #667a78; font-size: 13px; line-height: 1.6; }
.budget-empty { margin-bottom: 14px; }
.empty-state { min-height: calc(100vh - 190px); display: flex; flex-direction: column; align-items: center; justify-content: center; padding: 28px 8px 80px; text-align: center; box-sizing: border-box; }
.empty-visual { width: 72px; height: 72px; border-radius: 28px; background: #102033; color: #b9ead9; display: flex; align-items: center; justify-content: center; font-size: 24px; font-weight: 900; box-shadow: 0 16px 30px rgba(16,32,51,.16); }
.empty-title { display: block; margin-top: 18px; color: #102033; font-size: 22px; font-weight: 900; }
.empty-desc { display: block; max-width: 310px; margin-top: 10px; color: #667a78; font-size: 14px; line-height: 1.65; }
.empty-btn { height: 44px; margin-top: 22px; padding: 0 22px; border-radius: 22px; background: #102033; color: #f8fbf9; display: flex; align-items: center; justify-content: center; font-size: 14px; font-weight: 800; }
.save-btn { position: fixed; left: 50%; bottom: calc(env(safe-area-inset-bottom) + 78px); z-index: 25; width: calc(100% - 32px); max-width: 398px; height: 48px; transform: translateX(-50%); border-radius: 24px; background: #102033; color: #fff; display: flex; align-items: center; justify-content: center; font-size: 15px; font-weight: 800; box-shadow: 0 14px 28px rgba(16,32,51,.18); }
</style>
