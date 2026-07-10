<template>
  <view class="page">
    <view class="phone-shell">
      <view class="topbar">
        <text class="back" @tap="goBack">‹</text>
        <text class="top-title">我的收藏</text>
        <text class="top-action">{{ favoriteCount }}条</text>
      </view>

      <view class="hero-card">
        <text class="hero-mark">藏</text>
        <text class="hero-title">喜欢的地方，会留在这里</text>
        <text class="hero-desc">看到心动的目的地、攻略或酒店后，收藏一下，下次回来会更快找到。</text>
      </view>

      <view class="empty-card">
        <text class="empty-title">还没有收藏内容</text>
        <text class="empty-desc">先去逛逛探索页，找到喜欢的旅行灵感后，这里就会慢慢丰富起来。</text>
        <view class="empty-meta">当前已收藏 {{ favoriteCount }} 条</view>
        <view class="empty-btn" @tap="goDiscover">去探索</view>
      </view>
    </view>
  </view>
</template>

<script setup>
import { onMounted, ref } from "vue";
import { requestJSON } from "../../utils/api.js";

const favoriteCount = ref(0);

onMounted(() => loadDashboard());

async function loadDashboard() {
  try {
    const data = await requestJSON("/api/v1/user/dashboard");
    favoriteCount.value = Number(data?.stats?.favoriteCount || 0);
  } catch (error) {
    console.warn("读取收藏统计失败", error);
  }
}

function goBack() { uni.navigateBack(); }
function goDiscover() { uni.switchTab({ url: "/pages/discover/index" }); }
</script>

<style scoped>
.page { min-height: 100vh; background: #dfe9e6; }
.phone-shell { width: 100%; max-width: 430px; min-height: 100vh; margin: 0 auto; padding: calc(env(safe-area-inset-top) + 10px) 18px 28px; box-sizing: border-box; background: linear-gradient(180deg,#fbfdfb 0%,#eef6f3 48%,#f8faf7 100%); }
.topbar { height: 48px; display: flex; align-items: center; justify-content: space-between; }
.back { width: 38px; height: 38px; border-radius: 14px; background: #fffdfa; color: #0f5f75; display: flex; align-items: center; justify-content: center; font-size: 24px; box-shadow: 0 8px 20px rgba(16,32,51,.08); }
.top-title { color: #102033; font-size: 17px; font-weight: 900; }
.top-action { color: #0f5f75; font-size: 13px; font-weight: 900; }
.hero-card, .empty-card { margin-top: 16px; padding: 18px; border-radius: 24px; background: rgba(255,255,252,.98); border: 1px solid rgba(224,232,229,.9); box-shadow: 0 10px 24px rgba(16,32,51,.07); }
.hero-mark { width: 42px; height: 42px; border-radius: 16px; background: #edf7f3; color: #0f5f75; display: flex; align-items: center; justify-content: center; font-size: 16px; font-weight: 900; }
.hero-title, .empty-title { display: block; margin-top: 14px; color: #102033; font-size: 20px; font-weight: 900; line-height: 1.45; }
.hero-desc, .empty-desc { display: block; margin-top: 8px; color: #667a78; font-size: 13px; line-height: 1.7; }
.empty-meta { margin-top: 12px; color: #0f5f75; font-size: 12px; font-weight: 900; }
.empty-btn { width: fit-content; min-width: 120px; height: 40px; margin-top: 14px; padding: 0 16px; border-radius: 20px; background: #102033; color: #f8fbf9; display: flex; align-items: center; justify-content: center; font-size: 13px; font-weight: 900; }
</style>
