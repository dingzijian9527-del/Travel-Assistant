<template>
  <view class="page">
    <view class="phone-shell with-nav">
      <view class="topbar">
        <text class="title">我的</text>
        <view class="top-actions">
          <text class="head-icon" @tap="goPage('/pages/messages/index')">通</text>
          <text class="head-icon" @tap="goPage('/pages/settings/index')">设</text>
        </view>
      </view>

      <view class="profile-card">
        <view class="profile-row">
          <view class="avatar-wrap">
            <image :src="avatarSrc" mode="aspectFill" />
            <view v-if="memberLevel" class="level-badge">{{ memberLevel }}</view>
          </view>
          <view class="user-info">
            <text class="user-name">{{ userName }}</text>
            <view class="member-tag">
              <text v-if="memberLevel" class="vip-badge">{{ memberLevel }}</text>
              <text class="user-city">{{ cityText }}</text>
            </view>
            <text class="edit-link" @tap="goPage('/pages/profile/edit')">编辑资料</text>
          </view>
          <text class="qr-btn" @tap="goPage('/pages/member-code/index')">码</text>
        </view>
        <view class="profile-stats">
          <view class="stat-item" @tap="goTrip">
            <text class="stat-num blue">{{ tripCount }}</text>
            <text class="stat-label">已规划行程</text>
          </view>
          <view class="stat-divider"></view>
          <view class="stat-item" @tap="goPage('/pages/favorites/index')">
            <text class="stat-num green">{{ favoriteCount }}</text>
            <text class="stat-label">收藏目的地</text>
          </view>
        </view>
      </view>

      <view class="pref-card">
        <view class="pref-left">
          <text class="pref-title">旅行偏好</text>
          <text class="pref-desc">{{ preferenceSummary }}</text>
        </view>
        <button class="pref-btn" @tap="goPage('/pages/preferences/index')">编辑</button>
      </view>

      <view class="menu-list">
        <view class="menu-item" v-for="item in list" :key="item.name" @tap="openMenu(item.name)">
          <view class="menu-left">
            <view class="menu-icon" :style="{ background: item.bg }">
              <text :style="{ color: item.color }">{{ item.icon }}</text>
            </view>
            <text class="menu-name">{{ item.name }}</text>
          </view>
          <view class="menu-right">
            <text v-if="item.extra" class="menu-extra">{{ item.extra }}</text>
            <text v-if="item.dot" class="menu-dot"></text>
            <text class="menu-arrow">›</text>
          </view>
        </view>
      </view>

      <view class="invite-card">
        <view class="invite-left">
          <text class="invite-title">邀请好友得旅行基金</text>
          <text class="invite-desc">分享后可在邀请页查看奖励进度</text>
        </view>
        <button class="invite-btn" @tap="goPage('/pages/invite/index')">去邀请</button>
      </view>

      <view class="logout-btn" @tap="logout">
        <text>退出登录</text>
      </view>

      <BottomNav active="profile" />
    </view>
  </view>
</template>

<script setup>
import { computed, onMounted, ref } from "vue";
import { onShow } from "@dcloudio/uni-app";
import BottomNav from "../../components/BottomNav.vue";
import { clearStoredUser, requestJSON, setStoredUser, storedUserRef, syncStoredUser } from "../../utils/api.js";

const avatar = "/static/avatar-traveler.svg";
const dashboard = ref(null);
const userInfo = computed(() => dashboard.value?.user || storedUserRef.value || {});
const avatarSrc = computed(() => userInfo.value.avatarUrl || avatar);
const stats = computed(() => dashboard.value?.stats || {});
const userName = computed(() => userInfo.value.nickname || "未设置昵称");
const memberLevel = computed(() => userInfo.value.memberLevel || "");
const cityText = computed(() => userInfo.value.currentCity || userInfo.value.homeCity || "未设置常驻城市");
const tripCount = computed(() => stats.value.tripCount || 0);
const favoriteCount = computed(() => stats.value.favoriteCount || 0);
const unreadCount = computed(() => stats.value.unreadCount || 0);
const couponCount = computed(() => stats.value.couponCount || 0);
const preferenceSummary = computed(() => {
  const items = dashboard.value?.preferences?.items || [];
  return items.length ? items.join(" · ") : "设置偏好后会为你推荐更合适的行程";
});
const list = computed(() => [
  { name: "常用旅客", icon: "旅", bg: "#f0f2f8", color: "#4f5a6b" },
  { name: "消息通知", icon: "信", bg: "#fff5e6", color: "#e8830d", extra: unreadCount.value ? `${unreadCount.value}条未读` : "", dot: unreadCount.value > 0 },
  { name: "账号与安全", icon: "安", bg: "#edf7f3", color: "#0f5f75" },
  { name: "优惠券", icon: "券", bg: "#fff0e6", color: "#e8600d", extra: `${couponCount.value || 0}张` },
  { name: "设置", icon: "设", bg: "#f0f2f8", color: "#4f5a6b" },
]);

onMounted(() => {
  uni.hideTabBar && uni.hideTabBar();
  loadDashboard();
});

onShow(() => {
  syncStoredUser();
  loadDashboard();
});

async function loadDashboard() {
  try {
    const data = await requestJSON("/api/v1/user/dashboard");
    dashboard.value = data;
    if (data?.user) {
      setStoredUser(data.user);
    }
  } catch (error) {
    console.warn("读取个人中心数据失败", error);
  }
}

function logout() {
  uni.removeStorageSync("travel_token");
  clearStoredUser();
  uni.reLaunch({ url: "/pages/login/index" });
}

function toast(title) {
  uni.showToast({ title, icon: "none" });
}

function goTrip() {
  uni.switchTab({ url: "/pages/trip/index" });
}

function goPage(url) {
  uni.navigateTo({ url });
}

function openMenu(name) {
  const routeMap = {
    常用旅客: "/pages/passengers/index",
    消息通知: "/pages/messages/index",
    账号与安全: "/pages/security/index",
    优惠券: "/pages/coupons/index",
    设置: "/pages/settings/index",
  };
  if (routeMap[name]) {
    goPage(routeMap[name]);
    return;
  }
  toast(name);
}
</script>

<style scoped>
.page { min-height: 100vh; background: #dfe9e6; }
.phone-shell { width: 100%; max-width: 430px; min-height: 100vh; margin: 0 auto; padding: calc(env(safe-area-inset-top) + 8px) 18px 0; box-sizing: border-box; background: linear-gradient(180deg, #fbfdfb 0%, #eef6f3 46%, #f8faf7 100%); }
.with-nav { padding-bottom: calc(env(safe-area-inset-bottom) + 92px); }
.topbar { display: flex; align-items: center; justify-content: space-between; min-height: 56px; }
.title { color: #102033; font-size: 24px; font-weight: 800; }
.top-actions { display: flex; gap: 12px; }
.head-icon { width: 36px; height: 36px; border-radius: 14px; background: #fffdfa; color: #0f5f75; display: flex; align-items: center; justify-content: center; font-size: 13px; font-weight: 900; box-shadow: 0 8px 20px rgba(16,32,51,.08); }
.profile-card,
.menu-list,
.invite-card { background: rgba(255,255,252,.98); border: 1px solid rgba(224,232,229,.9); box-shadow: 0 12px 28px rgba(16,32,51,.07); border-radius: 22px; }
.profile-card { padding: 20px 18px; }
.profile-row { display: flex; align-items: center; }
.avatar-wrap { position: relative; width: 72px; height: 72px; flex: 0 0 72px; }
.avatar-wrap image { width: 72px; height: 72px; border-radius: 999px; border: 2px solid #fff; box-shadow: 0 10px 22px rgba(16,32,51,.14); }
.level-badge { position: absolute; right: 0; bottom: 0; padding: 2px 10px; border-radius: 999px; background: #caa459; color: #102033; font-size: 11px; font-weight: 900; border: 2px solid #fff; }
.user-info { flex: 1; min-width: 0; margin-left: 14px; }
.user-name { display: block; color: #102033; font-size: 22px; font-weight: 800; }
.member-tag { display: flex; align-items: center; gap: 8px; margin-top: 6px; }
.edit-link { display: inline-flex; margin-top: 10px; color: #0f5f75; font-size: 13px; font-weight: 800; }
.vip-badge { padding: 3px 10px; border-radius: 999px; background: #edf7f3; color: #0f5f75; font-size: 12px; font-weight: 800; }
.user-city { color: #667a78; font-size: 13px; }
.qr-btn { width: 40px; height: 40px; border-radius: 13px; background: #edf7f3; display: flex; align-items: center; justify-content: center; color: #0f5f75; font-size: 14px; font-weight: 900; }
.profile-stats { display: flex; align-items: center; margin-top: 18px; padding-top: 16px; border-top: 1px solid rgba(224,232,229,.92); }
.stat-item { flex: 1; display: flex; flex-direction: column; align-items: center; gap: 4px; }
.stat-divider { width: 1px; height: 32px; background: #e8ebf2; }
.stat-num { font-size: 28px; font-weight: 900; line-height: 1.2; }
.stat-num.blue,
.stat-num.green { color: #0f5f75; }
.stat-label { color: #667a78; font-size: 13px; }
.pref-card { margin-top: 16px; padding: 16px 18px; border-radius: 22px; background: #102033; display: flex; align-items: center; justify-content: space-between; box-shadow: 0 16px 32px rgba(16,32,51,.16); }
.pref-left { flex: 1; min-width: 0; }
.pref-title { display: block; color: #f8fbf9; font-size: 18px; font-weight: 800; }
.pref-desc { display: block; margin-top: 6px; color: rgba(248,251,249,.72); font-size: 12px; }
.pref-btn { width: 64px; height: 32px; border-radius: 16px; background: #b9ead9; color: #102033; font-size: 13px; font-weight: 800; flex: 0 0 64px; }
.menu-list { margin-top: 16px; overflow: hidden; }
.menu-item { display: flex; align-items: center; justify-content: space-between; min-height: 58px; padding: 0 18px; border-bottom: 1px solid rgba(224,232,229,.92); }
.menu-item:last-child { border-bottom: 0; }
.menu-left { display: flex; align-items: center; gap: 14px; }
.menu-icon { width: 36px; height: 36px; border-radius: 13px; display: flex; align-items: center; justify-content: center; font-size: 17px; font-weight: 900; flex: 0 0 36px; }
.menu-name { color: #102033; font-size: 16px; font-weight: 600; }
.menu-right { display: flex; align-items: center; gap: 8px; }
.menu-extra { color: #667a78; font-size: 13px; }
.menu-dot { width: 8px; height: 8px; border-radius: 999px; background: #e53e5c; }
.menu-arrow { color: #98a1b3; font-size: 24px; line-height: 1; }
.invite-card { margin-top: 16px; padding: 16px 18px; background: #fffdfa; border: 1px solid #ead7a8; display: flex; align-items: center; justify-content: space-between; }
.invite-left { flex: 1; min-width: 0; }
.invite-title { display: block; color: #102033; font-size: 17px; font-weight: 800; }
.invite-desc { display: block; margin-top: 4px; color: #667a78; font-size: 12px; }
.invite-btn { width: 82px; height: 34px; border-radius: 17px; background: #0f5f75; color: #f8fbf9; font-size: 13px; font-weight: 800; flex: 0 0 82px; }
.logout-btn { margin-top: 26px; margin-bottom: 16px; height: 52px; border-radius: 18px; background: rgba(255,255,252,.98); border: 1px solid rgba(224,232,229,.92); display: flex; align-items: center; justify-content: center; color: #b6403f; font-size: 16px; font-weight: 700; }
</style>
