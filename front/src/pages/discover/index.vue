<template>
  <view class="page">
    <view class="phone-shell with-nav">
      <view class="topbar">
        <text class="title">探索</text>
        <text class="head-placeholder"></text>
      </view>

      <view class="hero-card">
        <text class="hero-kicker">下一次旅行，从这里开始</text>
        <text class="hero-title">找一个适合现在出发的方向</text>
        <text class="hero-desc">输入目的地、假期天数或预算，旅行助手会帮你整理玩法、路线和行前建议。</text>
      </view>

      <view class="search-bar">
        <text class="search-icon">寻</text>
        <input v-model="keyword" placeholder="搜目的地、玩法或出行灵感" @confirm="askAI" />
      </view>

      <view class="action-card">
        <view class="action-item" @tap="sendIdea('请推荐适合周末短途放松的旅行目的地')">
          <text class="action-name">周末短途</text>
          <text class="action-note">轻松出发</text>
        </view>
        <view class="action-item" @tap="sendIdea('请推荐适合情侣出行的旅行目的地，并控制预算')">
          <text class="action-name">情侣旅行</text>
          <text class="action-note">氛围舒适</text>
        </view>
        <view class="action-item" @tap="sendIdea('请推荐适合亲子出行的旅行目的地和玩法')">
          <text class="action-name">亲子出行</text>
          <text class="action-note">省心好玩</text>
        </view>
      </view>

      <view class="empty-card">
        <text class="empty-title">这里会慢慢装满你的旅行灵感</text>
        <text class="empty-desc">先告诉我你想去哪、玩几天，或者让助手先给你几个热门推荐。</text>
        <view class="primary-btn" @tap="askAI">去问旅行助手</view>
      </view>

      <BottomNav active="explore" />
    </view>
  </view>
</template>

<script setup>
import { ref, onMounted } from "vue";
import BottomNav from "../../components/BottomNav.vue";

const keyword = ref("");

onMounted(() => uni.hideTabBar && uni.hideTabBar());

function askAI() {
  const question = keyword.value.trim() || "请根据我的需求推荐合适的旅行目的地和玩法";
  uni.navigateTo({ url: `/pages/ai/index?autoSend=1&question=${encodeURIComponent(question)}` });
}

function sendIdea(text) {
  keyword.value = text;
  askAI();
}
</script>

<style scoped>
.page { min-height: 100vh; background: #dfe9e6; }
.phone-shell { width: 100%; max-width: 430px; min-height: 100vh; margin: 0 auto; padding: calc(env(safe-area-inset-top) + 10px) 18px 0; box-sizing: border-box; background: linear-gradient(180deg, #fbfdfb 0%, #eef6f3 46%, #f8faf7 100%); }
.with-nav { padding-bottom: calc(env(safe-area-inset-bottom) + 92px); }
.topbar { min-height: 56px; display: flex; align-items: center; justify-content: space-between; }
.title { color: #102033; font-size: 24px; font-weight: 900; }
.head-placeholder { width: 40px; }
.hero-card, .action-card, .empty-card, .search-bar { margin-top: 16px; border-radius: 24px; background: rgba(255,255,252,.98); border: 1px solid rgba(224,232,229,.9); box-shadow: 0 10px 24px rgba(16,32,51,.07); }
.hero-card { margin-top: 10px; padding: 22px 18px; background: linear-gradient(135deg, #102033 0%, #17354a 100%); border-color: rgba(16,32,51,.06); }
.hero-kicker { display: block; color: rgba(185,234,217,.92); font-size: 12px; font-weight: 900; }
.hero-title { display: block; margin-top: 10px; color: #f8fbf9; font-size: 26px; font-weight: 900; line-height: 1.35; }
.hero-desc { display: block; margin-top: 10px; color: rgba(248,251,249,.78); font-size: 13px; line-height: 1.7; }
.search-bar { min-height: 52px; padding: 0 14px; display: flex; align-items: center; gap: 10px; }
.search-icon { width: 28px; height: 28px; border-radius: 10px; background: #edf7f3; color: #0f5f75; display: flex; align-items: center; justify-content: center; font-size: 12px; font-weight: 900; }
.search-bar input { flex: 1; min-width: 0; color: #102033; font-size: 14px; }
.action-card { padding: 12px; display: grid; grid-template-columns: 1fr; gap: 10px; }
.action-item { padding: 16px; border-radius: 18px; background: #f7faf8; }
.action-name { display: block; color: #102033; font-size: 15px; font-weight: 900; }
.action-note { display: block; margin-top: 4px; color: #667a78; font-size: 12px; }
.empty-card { padding: 18px; }
.empty-title { display: block; color: #102033; font-size: 20px; font-weight: 900; line-height: 1.45; }
.empty-desc { display: block; margin-top: 8px; color: #667a78; font-size: 13px; line-height: 1.7; }
.primary-btn { width: fit-content; min-width: 124px; height: 42px; margin-top: 16px; padding: 0 16px; border-radius: 21px; background: #102033; color: #f8fbf9; display: flex; align-items: center; justify-content: center; font-size: 13px; font-weight: 900; }
</style>

