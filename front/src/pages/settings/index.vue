<template>
  <view class="page">
    <view class="phone-shell">
      <view class="topbar">
        <text class="back" @tap="goBack">‹</text>
        <text class="top-title">设置</text>
        <text class="top-action"></text>
      </view>
      <view class="section-card">
        <view v-for="item in settings" :key="item.title" class="setting-row" @tap="toggle(item)">
          <view>
            <text class="setting-title">{{ item.title }}</text>
            <text class="setting-desc">{{ item.desc }}</text>
          </view>
          <view :class="['switch', item.enabled ? 'on' : '']">
            <view></view>
          </view>
        </view>
      </view>
      <view class="section-card">
        <view class="plain-row" v-for="item in plainItems" :key="item" @tap="toast(item)">
          <text>{{ item }}</text>
          <text class="arrow">›</text>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup>
import { onMounted, ref } from "vue";
import { requestJSON } from "../../utils/api.js";

const settings = ref([
  { key: "tripReminderEnabled", title: "行程提醒", desc: "出发前、景点预约前及时提醒", enabled: true },
  { key: "priceReminderEnabled", title: "价格波动提醒", desc: "酒店和交通价格变动时通知你", enabled: true },
  { key: "personalizedRecommendEnabled", title: "个性化推荐", desc: "根据你的偏好推荐目的地和玩法", enabled: false },
]);
const plainItems = ["清理缓存", "关于旅行助手", "用户协议与隐私政策"];

function goBack() {
  uni.navigateBack();
}

onMounted(() => {
  uni.hideTabBar && uni.hideTabBar();
  loadSettings();
});

async function loadSettings() {
  try {
    const data = await requestJSON("/api/v1/user/settings");
    applySettings(data);
  } catch (error) {
    console.warn("读取设置失败", error);
  }
}

function applySettings(data) {
  const next = data || {};
  settings.value = settings.value.map((item) => ({
    ...item,
    enabled: !!next[item.key],
  }));
}

async function toggle(item) {
  const previous = item.enabled;
  item.enabled = !item.enabled;
  try {
    await requestJSON("/api/v1/user/settings", {
      method: "POST",
      data: {
        trip_reminder_enabled: settings.value.find((value) => value.key === "tripReminderEnabled")?.enabled || false,
        price_reminder_enabled: settings.value.find((value) => value.key === "priceReminderEnabled")?.enabled || false,
        personalized_recommend_enabled: settings.value.find((value) => value.key === "personalizedRecommendEnabled")?.enabled || false,
      },
    });
  } catch (error) {
    item.enabled = previous;
    uni.showToast({ title: error.message || "保存失败", icon: "none" });
  }
}

function toast(title) {
  uni.showToast({ title, icon: "none" });
}
</script>

<style scoped>
.page { min-height: 100vh; background: #dfe9e6; }
.phone-shell { width: 100%; max-width: 430px; min-height: 100vh; margin: 0 auto; padding: calc(env(safe-area-inset-top) + 10px) 18px 28px; box-sizing: border-box; background: linear-gradient(180deg,#fbfdfb 0%,#eef6f3 48%,#f8faf7 100%); }
.topbar { height: 48px; display: flex; align-items: center; justify-content: space-between; }
.back { width: 38px; height: 38px; border-radius: 14px; background: #fffdfa; color: #0f5f75; display: flex; align-items: center; justify-content: center; font-size: 32px; box-shadow: 0 8px 20px rgba(16,32,51,.08); }
.top-title { color: #102033; font-size: 17px; font-weight: 900; }
.top-action { width: 38px; }
.section-card { margin-top: 16px; border-radius: 22px; overflow: hidden; background: rgba(255,255,252,.98); border: 1px solid rgba(224,232,229,.9); box-shadow: 0 10px 24px rgba(16,32,51,.07); }
.setting-row,.plain-row { min-height: 66px; padding: 0 16px; border-bottom: 1px solid #edf2ef; display: flex; align-items: center; justify-content: space-between; }
.setting-row:last-child,.plain-row:last-child { border-bottom: 0; }
.setting-title,.plain-row text:first-child { display: block; color: #102033; font-size: 15px; font-weight: 900; }
.setting-desc { display: block; margin-top: 4px; color: #667a78; font-size: 12px; }
.switch { width: 48px; height: 28px; border-radius: 999px; background: #d8e2df; padding: 3px; box-sizing: border-box; transition: background-color .18s ease; }
.switch view { width: 22px; height: 22px; border-radius: 999px; background: #fff; transition: transform .18s ease; }
.switch.on { background: #0f5f75; }
.switch.on view { transform: translateX(20px); }
.arrow { color: #8a9894; font-size: 22px; }
</style>
