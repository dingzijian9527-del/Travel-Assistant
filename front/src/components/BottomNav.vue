<template>
  <view class="bottom-nav">
    <view :class="['nav-item', active === 'explore' ? 'active' : '']" @tap="go('/pages/index/index')"><text class="nav-icon">◉</text><text>探索</text></view>
    <view :class="['nav-item', active === 'trip' ? 'active' : '']" @tap="go('/pages/trip/index')"><text class="nav-icon">▤</text><text>行程</text></view>
    <view :class="['nav-item', 'ai', active === 'ai' ? 'active' : '']" @tap="goAi"><text class="nav-icon">▣</text><text>智能助手</text></view>
    <view :class="['nav-item', active === 'profile' ? 'active' : '']" @tap="go('/pages/profile/index')"><text class="nav-icon">♙</text><text>我的</text></view>
  </view>
</template>

<script setup>
defineProps({ active: { type: String, default: "explore" } });
function go(url) { uni.switchTab({ url }); }
function goAi() { uni.navigateTo({ url: "/pages/ai/index" }); }
</script>

<style scoped>
.bottom-nav { position: fixed; left: 50%; bottom: 0; z-index: 20; width: 390px; max-width: 100vw; transform: translateX(-50%); height: 82px; padding: 8px 12px 10px; display: grid; grid-template-columns: repeat(4, 1fr); align-items: center; background: rgba(255,255,255,.96); border-top: 1px solid #e2e6f0; box-shadow: 0 -8px 28px rgba(40,55,85,.08); box-sizing: border-box; }
.nav-item { height: 58px; border-radius: 30px; color: #27303f; display: flex; flex-direction: column; align-items: center; justify-content: center; gap: 4px; font-size: 12px; font-weight: 800; }
.nav-icon { font-size: 22px; line-height: 1; }
.nav-item.active { color: #fff; background: #0b5be7; box-shadow: 0 10px 22px rgba(11,91,231,.28); }
.nav-item.ai { transform: translateY(-12px); }
.nav-item.ai .nav-icon { width: 42px; height: 42px; border-radius: 999px; background: #0b5be7; color: #fff; display: flex; align-items: center; justify-content: center; box-shadow: 0 8px 18px rgba(11,91,231,.28); }
.nav-item.ai.active .nav-icon { background: transparent; box-shadow: none; }
.nav-item.ai.active { transform: translateY(-12px); }
</style>
