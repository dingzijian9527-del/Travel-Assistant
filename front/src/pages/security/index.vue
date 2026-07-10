<template>
  <view class="page">
    <view class="phone-shell">
      <view class="topbar">
        <text class="back" @tap="goBack">‹</text>
        <text class="top-title">账号与安全</text>
        <text class="top-action"></text>
      </view>

      <view class="safe-card">
        <text class="safe-score">{{ safeTitle }}</text>
        <text class="safe-desc">{{ safeDesc }}</text>
      </view>

      <view class="menu-list">
        <view class="menu-item">
          <view>
            <text class="menu-title">绑定手机号</text>
            <text class="menu-desc">{{ maskedPhone }}</text>
          </view>
        </view>
        <view class="menu-item">
          <view>
            <text class="menu-title">账户状态</text>
            <text class="menu-desc">{{ user?.accountStatus || '待完善' }}</text>
          </view>
        </view>
        <view class="menu-item">
          <view>
            <text class="menu-title">个性化推荐</text>
            <text class="menu-desc">{{ recommendStatus }}</text>
          </view>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup>
import { computed, onMounted, ref } from "vue";
import { requestJSON } from "../../utils/api.js";

const user = ref(null);
const settings = ref(null);

const maskedPhone = computed(() => {
  const phone = String(user.value?.phone || "").trim();
  return phone.length === 11 ? `${phone.slice(0, 3)}****${phone.slice(7)}` : "未绑定";
});
const recommendStatus = computed(() => settings.value?.personalizedRecommendEnabled ? "已开启" : "未开启");
const safeTitle = computed(() => user.value?.phone ? "账号状态良好" : "登录后查看账号信息");
const safeDesc = computed(() => user.value?.phone ? "你可以查看常用的账号信息和推荐设置，后续也会逐步补充更多安全能力。" : "登录后可查看手机号、账户状态和推荐设置。" );

onMounted(() => loadDashboard());

async function loadDashboard() {
  try {
    const data = await requestJSON("/api/v1/user/dashboard");
    user.value = data?.user || null;
    settings.value = data?.settings || null;
  } catch (error) {
    console.warn("读取账号安全信息失败", error);
    user.value = null;
    settings.value = null;
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
.top-action { width: 38px; }
.safe-card { margin-top: 14px; padding: 18px; border-radius: 24px; background: #102033; box-shadow: 0 16px 32px rgba(16,32,51,.16); }
.safe-score { display: block; color: #b9ead9; font-size: 22px; font-weight: 900; }
.safe-desc { display: block; margin-top: 8px; color: rgba(248,251,249,.72); font-size: 13px; line-height: 1.6; }
.menu-list { margin-top: 16px; border-radius: 22px; overflow: hidden; background: rgba(255,255,252,.98); border: 1px solid rgba(224,232,229,.9); box-shadow: 0 10px 24px rgba(16,32,51,.07); }
.menu-item { min-height: 66px; padding: 0 16px; border-bottom: 1px solid #edf2ef; display: flex; align-items: center; justify-content: space-between; }
.menu-item:last-child { border-bottom: 0; }
.menu-title { display: block; color: #102033; font-size: 15px; font-weight: 900; }
.menu-desc { display: block; margin-top: 4px; color: #667a78; font-size: 12px; }
</style>
