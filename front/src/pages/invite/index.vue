<template>
  <view class="page">
    <view class="phone-shell">
      <view class="topbar">
        <text class="back" @tap="goBack">‹</text>
        <text class="top-title">邀请好友</text>
        <text class="top-action"></text>
      </view>

      <view class="hero-card">
        <text class="hero-title">把旅行灵感分享给同伴</text>
        <text class="hero-desc">邀请功能开放后，你可以查看分享进度和奖励信息。</text>
      </view>

      <view class="empty-card">
        <text class="empty-title">邀请功能即将开放</text>
        <text class="empty-desc">开放后你可以一键分享给朋友，一起更快定下下一次旅行。</text>
      </view>
    </view>
  </view>
</template>

<script setup>
function goBack() { uni.navigateBack(); }
</script>

<style scoped>
.page { min-height: 100vh; background: #dfe9e6; }
.phone-shell { width: 100%; max-width: 430px; min-height: 100vh; margin: 0 auto; padding: calc(env(safe-area-inset-top) + 10px) 18px 28px; box-sizing: border-box; background: linear-gradient(180deg,#fbfdfb 0%,#eef6f3 48%,#f8faf7 100%); }
.topbar { height: 48px; display: flex; align-items: center; justify-content: space-between; }
.back { width: 38px; height: 38px; border-radius: 14px; background: #fffdfa; color: #0f5f75; display: flex; align-items: center; justify-content: center; font-size: 24px; box-shadow: 0 8px 20px rgba(16,32,51,.08); }
.top-title { color: #102033; font-size: 17px; font-weight: 900; }
.top-action { width: 38px; }
.hero-card, .empty-card { margin-top: 16px; padding: 18px; border-radius: 24px; background: rgba(255,255,252,.98); border: 1px solid rgba(224,232,229,.9); box-shadow: 0 10px 24px rgba(16,32,51,.07); }
.hero-title, .empty-title { display: block; color: #102033; font-size: 20px; font-weight: 900; line-height: 1.45; }
.hero-desc, .empty-desc { display: block; margin-top: 8px; color: #667a78; font-size: 13px; line-height: 1.7; }
</style>
