<template>
  <view class="page">
    <view class="phone-shell">
      <view class="topbar">
        <text class="back" @tap="goBack">‹</text>
        <text class="top-title">优惠券</text>
        <text class="top-action">{{ couponCount }}张</text>
      </view>

      <view class="hero-card">
        <text class="hero-title">有可用优惠时，这里会提醒你</text>
        <text class="hero-desc">领到的出行优惠、酒店券和活动福利都会放在这里，使用前也能更方便核对。</text>
      </view>

      <view class="empty-card">
        <text class="empty-title">还没有可用优惠</text>
        <text class="empty-desc">先去逛逛目的地和行程推荐，合适的活动和优惠出现后，这里会自动更新。</text>
        <view class="empty-meta">当前可用 {{ couponCount }} 张</view>
      </view>
    </view>
  </view>
</template>

<script setup>
import { onMounted, ref } from "vue";
import { requestJSON } from "../../utils/api.js";

const couponCount = ref(0);

onMounted(() => loadDashboard());

async function loadDashboard() {
  try {
    const data = await requestJSON("/api/v1/user/dashboard");
    couponCount.value = Number(data?.stats?.couponCount || 0);
  } catch (error) {
    console.warn("读取优惠券统计失败", error);
  }
}

function goBack() { uni.navigateBack(); }
</script>

<style scoped>
.page { min-height: 100vh; background: #dfe9e6; }
.phone-shell { width: 100%; max-width: 430px; min-height: 100vh; margin: 0 auto; padding: calc(env(safe-area-inset-top) + 10px) 18px 28px; box-sizing: border-box; background: linear-gradient(180deg,#fbfdfb 0%,#eef6f3 48%,#f8faf7 100%); }
.topbar { height: 48px; display: flex; align-items: center; justify-content: space-between; }
.back { width: 38px; height: 38px; border-radius: 14px; background: #fffdfa; color: #0f5f75; display: flex; align-items: center; justify-content: center; font-size: 24px; box-shadow: 0 8px 20px rgba(16,32,51,.08); }
.top-title { color: #102033; font-size: 17px; font-weight: 900; }
.top-action { color: #0f5f75; font-size: 13px; font-weight: 900; }
.hero-card, .empty-card { margin-top: 16px; padding: 18px; border-radius: 24px; background: rgba(255,255,252,.98); border: 1px solid rgba(224,232,229,.9); box-shadow: 0 10px 24px rgba(16,32,51,.07); }
.hero-title, .empty-title { display: block; color: #102033; font-size: 20px; font-weight: 900; line-height: 1.45; }
.hero-desc, .empty-desc { display: block; margin-top: 8px; color: #667a78; font-size: 13px; line-height: 1.7; }
.empty-meta { margin-top: 12px; color: #0f5f75; font-size: 12px; font-weight: 900; }
</style>
