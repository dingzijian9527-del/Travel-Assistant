<template>
  <view class="page">
    <view class="phone-shell">
      <view class="topbar">
        <text class="back" @tap="goBack">‹</text>
        <text class="top-title">攻略</text>
        <text class="top-action"></text>
      </view>

      <view class="hero-card">
        <text class="hero-kicker">路线参考</text>
        <text class="hero-title">{{ titleText }}</text>
        <text class="hero-desc">如果你已经有大致想法，我可以继续帮你把时间、路线、花费和住宿安排细化出来。</text>
      </view>

      <view class="content-card">
        <text class="content-title">继续完善这次出行</text>
        <text class="content-desc">告诉我你更在意风景、美食、亲子体验或预算范围，我会继续把这份计划细化完整。</text>
        <view class="content-btn" @tap="askPlan">让助手继续整理</view>
      </view>
    </view>
  </view>
</template>

<script setup>
import { computed, ref } from "vue";
import { onLoad } from "@dcloudio/uni-app";

const rawTitle = ref("");
const titleText = computed(() => rawTitle.value || "还没选定攻略主题");

onLoad((options) => {
  rawTitle.value = options?.title ? decodeURIComponent(options.title) : "";
});

function goBack() { uni.navigateBack(); }
function askPlan() {
  const question = rawTitle.value ? `请根据“${rawTitle.value}”这个主题，为我生成一份可执行的旅行计划。` : "请根据我的需求为我生成一份旅行计划。";
  uni.navigateTo({ url: `/pages/ai/index?autoSend=1&question=${encodeURIComponent(question)}` });
}
</script>

<style scoped>
.page { min-height: 100vh; background: #dfe9e6; }
.phone-shell { width: 100%; max-width: 430px; min-height: 100vh; margin: 0 auto; padding: calc(env(safe-area-inset-top) + 10px) 18px 28px; box-sizing: border-box; background: linear-gradient(180deg,#fbfdfb 0%,#eef6f3 46%,#f8faf7 100%); }
.topbar { height: 48px; display: flex; align-items: center; justify-content: space-between; }
.back { width: 38px; height: 38px; border-radius: 14px; background: #fffdfa; color: #0f5f75; display: flex; align-items: center; justify-content: center; font-size: 24px; box-shadow: 0 8px 20px rgba(16,32,51,.08); }
.top-title { color: #102033; font-size: 17px; font-weight: 900; }
.top-action { width: 38px; }
.hero-card, .content-card { margin-top: 16px; padding: 18px; border-radius: 24px; background: rgba(255,255,252,.98); border: 1px solid rgba(224,232,229,.9); box-shadow: 0 10px 24px rgba(16,32,51,.07); }
.hero-kicker { display: block; color: #08766c; font-size: 11px; font-weight: 900; }
.hero-title, .content-title { display: block; margin-top: 6px; color: #102033; font-size: 20px; font-weight: 900; line-height: 1.45; }
.hero-desc, .content-desc { display: block; margin-top: 8px; color: #667a78; font-size: 13px; line-height: 1.7; }
.content-btn { width: fit-content; min-width: 132px; height: 40px; margin-top: 14px; padding: 0 16px; border-radius: 20px; background: #102033; color: #f8fbf9; display: flex; align-items: center; justify-content: center; font-size: 13px; font-weight: 900; }
</style>

