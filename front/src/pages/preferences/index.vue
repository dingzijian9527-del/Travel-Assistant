<template>
  <view class="page">
    <view class="phone-shell">
      <view class="topbar">
        <text class="back" @tap="goBack">‹</text>
        <text class="top-title">旅行偏好</text>
        <text class="top-action" @tap="save">保存</text>
      </view>
      <view class="profile-card">
        <text class="profile-title">告诉我你更喜欢怎样旅行</text>
        <text class="profile-desc">选好你的出行偏好后，后续推荐和行程规划会更贴近你的习惯。</text>
      </view>
      <view v-for="group in groups" :key="group.title" class="section">
        <text class="section-title">{{ group.title }}</text>
        <view class="chip-grid">
          <text v-for="item in group.items" :key="item" :class="['chip', selected.includes(item) ? 'active' : '']" @tap="toggle(item)">{{ item }}</text>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup>
import { onMounted, ref } from "vue";
import { requestJSON } from "../../utils/api.js";

const selected = ref([]);
const groups = [
  { title: "兴趣", items: ["美食", "海岛", "博物馆", "茶馆", "骑行", "夜游", "亲子", "徒步"] },
  { title: "节奏", items: ["不赶路", "深度游", "早出发", "睡到自然醒", "少换酒店", "拍照优先"] },
  { title: "住宿", items: ["舒适酒店", "地铁附近", "景区内", "亲子房", "高性价比", "安静街区"] },
];

onMounted(() => {
  uni.hideTabBar && uni.hideTabBar();
  loadPreferences();
});

function goBack() { uni.navigateBack(); }
function toggle(item) {
  selected.value = selected.value.includes(item)
    ? selected.value.filter((value) => value !== item)
    : [...selected.value, item];
}
async function loadPreferences() {
  try {
    const data = await requestJSON("/api/v1/user/preferences");
    selected.value = Array.isArray(data?.items) ? data.items : [];
  } catch (error) {
    console.warn("读取旅行偏好失败", error);
    selected.value = [];
  }
}
async function save() {
  try {
    await requestJSON("/api/v1/user/preferences", {
      method: "POST",
      data: { items: selected.value },
    });
    uni.showToast({ title: "偏好已保存", icon: "none" });
  } catch (error) {
    uni.showToast({ title: error.message || "保存失败", icon: "none" });
  }
}
</script>

<style scoped>
.page { min-height: 100vh; background: #dfe9e6; }
.phone-shell { width: 100%; max-width: 430px; min-height: 100vh; margin: 0 auto; padding: calc(env(safe-area-inset-top) + 10px) 18px 28px; box-sizing: border-box; background: linear-gradient(180deg,#fbfdfb 0%,#eef6f3 48%,#f8faf7 100%); }
.topbar { height: 48px; display: flex; align-items: center; justify-content: space-between; }
.back { width: 38px; height: 38px; border-radius: 14px; background: #fffdfa; color: #0f5f75; display: flex; align-items: center; justify-content: center; font-size: 24px; box-shadow: 0 8px 20px rgba(16,32,51,.08); }
.top-title { color: #102033; font-size: 17px; font-weight: 900; }
.top-action { color: #0f5f75; font-size: 13px; font-weight: 900; }
.profile-card { margin-top: 14px; padding: 18px; border-radius: 24px; background: #102033; box-shadow: 0 16px 32px rgba(16,32,51,.16); }
.profile-title { display: block; color: #f8fbf9; font-size: 22px; font-weight: 900; }
.profile-desc { display: block; margin-top: 8px; color: rgba(248,251,249,.72); font-size: 13px; line-height: 1.6; }
.section { margin-top: 22px; }
.section-title { display: block; color: #102033; font-size: 18px; font-weight: 900; margin-bottom: 12px; }
.chip-grid { display: flex; flex-wrap: wrap; gap: 10px; padding: 16px; border-radius: 22px; background: rgba(255,255,252,.98); border: 1px solid rgba(224,232,229,.9); box-shadow: 0 10px 24px rgba(16,32,51,.07); }
.chip { min-height: 34px; padding: 0 13px; border-radius: 999px; background: #f1f5f2; color: #667a78; display: flex; align-items: center; font-size: 13px; font-weight: 900; }
.chip.active { background: #0f5f75; color: #f8fbf9; }
</style>

