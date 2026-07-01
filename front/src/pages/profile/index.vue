<template>
  <view class="page">
    <view class="phone-shell with-nav">
      <view class="topbar"><text class="title">我的</text><view><text class="head-icon" @tap="toast('暂无新通知')">♧</text><text class="head-icon" @tap="toast('设置入口')">⚙</text></view></view>
      <view class="profile-row"><view class="avatar-wrap"><image :src="avatar" mode="aspectFill" /><text>金卡</text></view><view class="user"><text class="name">张先生</text><view><text class="vip">黄金会员</text><text class="uid">常住成都</text></view></view><text class="qr" @tap="toast('会员码已打开')">▦</text></view>
      <view class="stats"><view @tap="goTrip"><text>已规划行程</text><strong>12<small>个</small></strong></view><view @tap="toast('打开收藏目的地')"><text>收藏目的地</text><strong class="green">34<small>处</small></strong></view></view>
      <view class="shortcut"><view v-for="item in shortcuts" :key="item.name" @tap="handleShortcut(item.name)"><text class="sicon">{{ item.icon }}</text><text>{{ item.name }}</text></view></view>
      <view class="member-card"><view><text>旅行偏好</text><small>川菜、海岛、城市漫步、舒适酒店</small></view><button @tap="toast('偏好已更新')">编辑</button></view>
      <view class="list"><view v-for="item in list" :key="item.name" class="row" @tap="toast(item.name)"><view><text class="licon">{{ item.icon }}</text><text>{{ item.name }}</text></view><view><text v-if="item.extra" class="extra">{{ item.extra }}</text><text class="arrow">›</text></view><text v-if="item.dot" class="dot"></text></view></view>
      <view class="invite"><view><text>邀请好友，得旅行基金</text><small>最高可获 ¥500 国内出行补贴</small></view><button @tap="toast('邀请链接已生成')">去邀请</button></view>
      <BottomNav active="profile" />
    </view>
  </view>
</template>

<script setup>
import { onMounted } from "vue";
import BottomNav from "../../components/BottomNav.vue";
const avatar = "https://images.unsplash.com/photo-1500648767791-00dcc994a43e?auto=format&fit=crop&w=160&q=80";
const shortcuts = [{ name: "我的行程", icon: "▤" }, { name: "我的收藏", icon: "♡" }, { name: "旅行偏好", icon: "☷" }];
const list = [{ name: "常用旅客", icon: "♙" }, { name: "消息通知", icon: "✉", extra: "3条未读", dot: true }, { name: "账号与安全", icon: "♜" }, { name: "优惠券", icon: "□", extra: "5张" }, { name: "设置", icon: "⚙" }];
onMounted(() => uni.hideTabBar && uni.hideTabBar());
function toast(title) { uni.showToast({ title, icon: "none" }); }
function goTrip() { uni.switchTab({ url: "/pages/trip/index" }); }
function handleShortcut(name) { if (name === "我的行程") goTrip(); else toast(name); }
</script>

<style scoped>
.page { min-height: 100vh; background: #f5f7fb; }
.phone-shell { width: 390px; max-width: 100vw; min-height: 100vh; margin: 0 auto; padding: 28px 18px 0; box-sizing: border-box; background: #f5f7fb; }
.with-nav { padding-bottom: 108px; }
.topbar { display: flex; align-items: center; justify-content: space-between; }
.title { color: #0b5be7; font-size: 26px; font-weight: 900; }
.head-icon { margin-left: 18px; color: #303746; font-size: 27px; }
.profile-row { display: flex; align-items: center; margin-top: 42px; }
.avatar-wrap { position: relative; width: 94px; height: 94px; border-radius: 999px; border: 3px solid #cbdcff; display: flex; align-items: center; justify-content: center; background: #fff; }
.avatar-wrap image { width: 72px; height: 72px; border-radius: 12px; }
.avatar-wrap text { position: absolute; right: -8px; bottom: 8px; padding: 3px 11px; border-radius: 999px; background: #07843b; color: #fff; font-size: 12px; font-weight: 900; }
.user { flex: 1; min-width: 0; margin-left: 24px; }
.name { display: block; color: #151a24; font-size: 28px; font-weight: 900; }
.vip { display: inline-flex; margin-top: 8px; padding: 4px 10px; border-radius: 999px; background: #dfe9ff; color: #0b5be7; font-size: 13px; font-weight: 900; }
.uid { margin-left: 10px; color: #2d3444; font-size: 14px; }
.qr { color: #0b5be7; font-size: 34px; }
.stats { display: grid; grid-template-columns: 1fr 1fr; gap: 16px; margin-top: 36px; }
.stats view { height: 112px; border-radius: 22px; background: #fff; padding: 22px 24px; box-sizing: border-box; box-shadow: 0 6px 18px rgba(31,45,71,.07); }
.stats text { display: block; color: #485063; font-size: 14px; }
.stats strong { display: block; margin-top: 22px; color: #0b5be7; font-size: 34px; line-height: 1; }
.stats .green { color: #07843b; }
.stats small { color: #171b25; font-size: 14px; margin-left: 3px; }
.shortcut { display: grid; grid-template-columns: repeat(3, 1fr); gap: 8px; margin-top: 22px; padding: 22px 18px; border-radius: 24px; background: #fff; box-shadow: 0 6px 18px rgba(31,45,71,.07); }
.shortcut view { display: flex; flex-direction: column; align-items: center; gap: 10px; color: #171b25; font-size: 14px; }
.sicon { color: #0b5be7; font-size: 28px; }
.member-card { margin-top: 20px; padding: 18px; border-radius: 20px; background: linear-gradient(135deg, #0b5be7, #44c3ff); display: flex; align-items: center; justify-content: space-between; color: #fff; box-shadow: 0 8px 22px rgba(11,91,231,.18); }
.member-card text { display: block; font-size: 18px; font-weight: 900; }
.member-card small { display: block; margin-top: 7px; max-width: 230px; font-size: 12px; opacity: .9; }
.member-card button { width: 68px; height: 34px; border-radius: 18px; background: #fff; color: #0b5be7; font-size: 13px; font-weight: 900; }
.list { margin-top: 22px; border-radius: 22px; overflow: hidden; background: #fff; box-shadow: 0 6px 18px rgba(31,45,71,.07); }
.row { position: relative; height: 68px; padding: 0 18px; display: flex; align-items: center; justify-content: space-between; border-bottom: 1px solid #edf0f6; box-sizing: border-box; color: #171b25; font-size: 18px; }
.row:last-child { border-bottom: 0; }
.row view { display: flex; align-items: center; }
.licon { width: 34px; color: #3d4557; font-size: 24px; }
.extra { color: #858b9b; font-size: 14px; margin-right: 12px; }
.arrow { color: #6a7283; font-size: 30px; }
.dot { position: absolute; left: 55px; top: 19px; width: 10px; height: 10px; border-radius: 999px; background: #d81f2a; }
.invite { margin-top: 24px; height: 96px; padding: 0 24px; border-radius: 24px; background: #e6e8ef; display: flex; align-items: center; justify-content: space-between; }
.invite text { display: block; color: #151a24; font-size: 20px; font-weight: 900; }
.invite small { display: block; margin-top: 8px; color: #4e5667; font-size: 13px; }
.invite button { width: 88px; height: 42px; border-radius: 22px; background: #0b5be7; color: #fff; font-size: 14px; font-weight: 900; }
</style>
