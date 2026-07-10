<template>
  <view class="bottom-nav">
    <view :class="['nav-item', active === 'explore' ? 'active' : '']" @tap="go('/pages/index/index')">
      <text class="nav-label">探索</text>
    </view>
    <view :class="['nav-item', active === 'trip' ? 'active' : '']" @tap="go('/pages/trip/index')">
      <text class="nav-label">行程</text>
    </view>
    <view :class="['nav-item', active === 'ai' ? 'active' : '']" @tap="goAi">
      <text class="nav-label">助手</text>
    </view>
    <view :class="['nav-item', active === 'profile' ? 'active' : '']" @tap="go('/pages/profile/index')">
      <text class="nav-label">我的</text>
    </view>
  </view>
</template>

<script setup>
const props = defineProps({ active: { type: String, default: "explore" } });
function go(url) { uni.switchTab({ url }); }
function goAi() {
  if (props.active === "ai") {
    return;
  }
  uni.navigateTo({ url: "/pages/ai/index" });
}
</script>

<style scoped>
.bottom-nav {
  position: fixed; left: 50%; bottom: 0; z-index: 100;
  width: 100%; max-width: 430px; min-height: 66px;
  transform: translateX(-50%);
  display: grid; grid-template-columns: repeat(4, minmax(0, 1fr));
  column-gap: 6px;
  align-items: center;
  background: rgba(255,255,253,.97);
  border-top: 1px solid rgba(202,214,214,.86);
  box-shadow: 0 -8px 24px rgba(16,32,51,.07);
  box-sizing: border-box;
  padding: 8px 12px calc(env(safe-area-inset-bottom) + 8px);
  backdrop-filter: blur(12px);
}
.nav-item {
  display: flex; flex-direction: column; align-items: center; justify-content: center;
  min-width: 0;
  min-height: 46px;
  border-radius: 18px;
  transition: background-color .18s ease, color .18s ease, transform .18s ease;
}
.nav-item.active {
  background: #e9f6f4;
  box-shadow: inset 0 0 0 1px rgba(210,236,234,.92);
}
.nav-label {
  font-size: 15px;
  line-height: 1;
  font-weight: 900;
  color: #7f8c99;
}
.nav-item.active .nav-label {
  color: #0c5f73;
  font-size: 16px;
}

</style>
