<template>
  <view class="page">
    <view class="phone-shell">
      <view class="topbar">
        <text class="back" @tap="goBack">‹</text>
        <text class="top-title">消息</text>
        <text class="top-action">{{ unreadCount }}条未读</text>
      </view>

      <view class="hero-card">
        <text class="hero-title">重要提醒都会准时送达</text>
        <text class="hero-desc">订单进度、出行提醒和活动消息都会集中展示，方便你随时查看。</text>
      </view>

      <view class="empty-card">
        <text class="empty-title">暂时没有新消息</text>
        <text class="empty-desc">一旦有新的订单动态或出行提醒，我们会第一时间更新到这里。</text>
        <view class="empty-meta">当前未读 {{ unreadCount }} 条</view>
      </view>
    </view>
  </view>
</template>

<script setup>
import { onMounted, ref } from "vue";
import { requestJSON } from "../../utils/api.js";

const unreadCount = ref(0);

onMounted(() => loadDashboard());

async function loadDashboard() {
  try {
    const data = await requestJSON("/api/v1/user/dashboard");
    unreadCount.value = Number(data?.stats?.unreadCount || 0);
  } catch (error) {
    console.warn("读取消息统计失败", error);
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
