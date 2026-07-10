<template>
  <view class="page">
    <view class="phone-shell">
      <view class="topbar">
        <text class="back" @tap="goBack">‹</text>
        <text class="top-title">会员码</text>
        <text class="top-action" @tap="loadDashboard">刷新</text>
      </view>

      <view v-if="user" class="member-card">
        <image class="avatar" :src="user.avatarUrl || avatar" mode="aspectFill" />
        <text class="member-name">{{ user.nickname || '未设置昵称' }}</text>
        <text class="member-level">{{ memberMetaText }}</text>
        <view class="code-box">
          <view v-for="row in 7" :key="row" class="code-row">
            <view v-for="col in 7" :key="col" :class="['code-dot', isDark(row, col) ? 'dark' : '']"></view>
          </view>
        </view>
        <text class="code-number">{{ codeText }}</text>
      </view>

      <view v-else class="empty-card">
        <text class="empty-title">登录后即可查看会员码</text>
        <text class="empty-desc">会员码可用于快速核验身份，出行时查看会更方便。</text>
      </view>

      <view class="info-card">
        <text class="info-title">账户信息</text>
        <view class="info-row">
          <text>手机号</text>
          <text>{{ maskedPhone }}</text>
        </view>
        <view class="info-row">
          <text>账户状态</text>
          <text>{{ user?.accountStatus || '待完善' }}</text>
        </view>
        <view class="info-row">
          <text>常驻城市</text>
          <text>{{ cityText }}</text>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup>
import { computed, onMounted, ref } from "vue";
import { onShow } from "@dcloudio/uni-app";
import { requestJSON, setStoredUser, storedUserRef, syncStoredUser } from "../../utils/api.js";

const avatar = "/static/avatar-traveler.svg";
const user = ref(storedUserRef.value);

const cityText = computed(() => user.value?.currentCity || user.value?.homeCity || "未设置常驻城市");
const maskedPhone = computed(() => maskPhone(user.value?.phone || ""));
const memberMetaText = computed(() => [user.value?.memberLevel, cityText.value].filter(Boolean).join(" · "));
const codeText = computed(() => {
  const source = String(user.value?.id || user.value?.phone || "登录后生成会员识别码");
  return source === "登录后生成会员识别码" ? source : `会员识别码 ${source}`;
});

onMounted(loadDashboard);
onShow(() => {
  user.value = syncStoredUser();
  loadDashboard();
});

async function loadDashboard() {
  try {
    const data = await requestJSON("/api/v1/user/dashboard");
    user.value = data?.user || null;
    if (data?.user) {
      setStoredUser(data.user);
    }
  } catch (error) {
    console.warn("读取会员信息失败", error);
    user.value = null;
  }
}

function isDark(row, col) {
  const sourceLength = String(user.value?.id || user.value?.phone || "0").length || 1;
  return (row * col + sourceLength) % 3 === 0;
}

function maskPhone(phone) {
  const text = String(phone || "").trim();
  if (text.length !== 11) {
    return "未绑定";
  }
  return `${text.slice(0, 3)}****${text.slice(7)}`;
}

function goBack() {
  uni.navigateBack();
}
</script>

<style scoped>
.page { min-height: 100vh; background: #dfe9e6; }
.phone-shell { width: 100%; max-width: 430px; min-height: 100vh; margin: 0 auto; padding: calc(env(safe-area-inset-top) + 10px) 18px 28px; box-sizing: border-box; background: linear-gradient(180deg,#fbfdfb 0%,#eef6f3 48%,#f8faf7 100%); }
.topbar { height: 48px; display: flex; align-items: center; justify-content: space-between; }
.back { width: 38px; height: 38px; border-radius: 14px; background: #fffdfa; color: #0f5f75; display: flex; align-items: center; justify-content: center; font-size: 24px; box-shadow: 0 8px 20px rgba(16,32,51,.08); }
.top-title { color: #102033; font-size: 17px; font-weight: 900; }
.top-action { color: #0f5f75; font-size: 13px; font-weight: 900; }
.member-card { margin-top: 16px; padding: 26px 18px; border-radius: 28px; background: #102033; display: flex; flex-direction: column; align-items: center; box-shadow: 0 18px 36px rgba(16,32,51,.18); }
.avatar { width: 74px; height: 74px; border-radius: 24px; border: 3px solid rgba(255,255,255,.22); }
.member-name { margin-top: 12px; color: #f8fbf9; font-size: 21px; font-weight: 900; }
.member-level { margin-top: 4px; color: rgba(248,251,249,.72); font-size: 12px; }
.code-box { margin-top: 24px; padding: 14px; border-radius: 20px; background: #fffdfa; display: flex; flex-direction: column; gap: 5px; }
.code-row { display: flex; gap: 5px; }
.code-dot { width: 16px; height: 16px; border-radius: 4px; background: #e6eee9; }
.code-dot.dark { background: #102033; }
.code-number { margin-top: 16px; color: #b9ead9; font-size: 13px; font-weight: 900; }
.empty-card, .info-card { margin-top: 16px; padding: 18px; border-radius: 24px; background: rgba(255,255,252,.98); border: 1px solid rgba(224,232,229,.9); box-shadow: 0 10px 24px rgba(16,32,51,.07); }
.empty-title, .info-title { display: block; color: #102033; font-size: 20px; font-weight: 900; }
.empty-desc { display: block; margin-top: 8px; color: #667a78; font-size: 13px; line-height: 1.7; }
.info-row { min-height: 42px; display: flex; align-items: center; justify-content: space-between; border-bottom: 1px solid #edf2ef; color: #3d4f5d; font-size: 13px; }
.info-row:last-child { border-bottom: 0; }
</style>
