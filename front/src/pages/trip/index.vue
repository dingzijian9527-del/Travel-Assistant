<template>
  <view class="page">
    <view class="phone-shell with-nav">
      <view class="topbar"><text class="back" @tap="goBack">‹</text><text class="title">行程预览</text><image class="avatar" :src="avatar" mode="aspectFill" /></view>
      <view class="map-card"><image src="https://images.unsplash.com/photo-1519181245277-cffeb31da2e3?auto=format&fit=crop&w=760&q=80" mode="aspectFill" /><button class="map-btn" @tap="toast('已打开模拟地图')">▱ 查看地图</button><view class="map-mask"></view><text class="destination">成都 · 3天2夜</text><text class="sub-destination">太古里 / 熊猫基地 / 宽窄巷子</text></view>
      <view class="summary"><view><text>出发</text><strong>06月28日</strong></view><view><text>同行</text><strong>2位成人</strong></view><view><text>预算</text><strong>¥3,200</strong></view></view>
      <text class="section-title">每日计划</text>
      <view class="day-tabs"><text v-for="day in days" :key="day.day" :class="activeDay === day.day ? 'active' : ''" @tap="activeDay = day.day">第{{ day.day }}天</text></view>
      <view class="day-block"><view class="day-head"><text class="num">{{ currentDay.day }}</text><text class="day-title">第 {{ currentDay.day }} 天：{{ currentDay.title }}</text></view><view class="line-wrap"><view class="dash"></view><view class="cards"><view class="route"><text class="card-title">⌘ 路线方案</text><text class="card-text">{{ currentDay.route }}</text></view><view class="mini-grid"><view v-for="tip in currentDay.tips" :key="tip.title" class="mini"><text>{{ tip.icon }} {{ tip.title }}</text><small>{{ tip.text }}</small></view></view><view class="weather">☼ {{ currentDay.weather }}</view></view></view></view>
      <text class="section-title budget-title">▣ 预算构成</text>
      <view class="budget"><view class="budget-row"><text>酒店与交通</text><strong>¥1,760</strong></view><view class="bar"><view class="fill blue"></view></view><view class="budget-row"><text>餐饮与门票</text><strong class="green">¥980</strong></view><view class="bar"><view class="fill green-bg"></view></view><view class="total"><text>预计总花费</text><strong>¥3,200</strong></view></view>
      <text class="section-title risk">△ 出行提醒</text>
      <view class="risk-card"><text>ⓘ 熊猫基地建议 08:30 前入园，热门场馆排队时间较长。</text><text>▭ 周末春熙路与太古里人流较大，建议提前预订晚餐。</text></view>
      <button class="save" @tap="saveTrip">▱ {{ saved ? '已保存到我的行程' : '保存行程至我的' }}</button>
      <BottomNav active="trip" />
    </view>
  </view>
</template>

<script setup>
import { computed, onMounted, ref } from "vue";
import BottomNav from "../../components/BottomNav.vue";
const avatar = "https://images.unsplash.com/photo-1500648767791-00dcc994a43e?auto=format&fit=crop&w=120&q=80";
const activeDay = ref(1);
const saved = ref(false);
const days = [
  { day: 1, title: "抵达成都与夜游太古里", route: "双流机场 → 春熙路酒店入住 → 太古里慢逛 → 玉林路晚餐", tips: [{ icon: "▣", title: "交通建议", text: "机场地铁 10 号线转 3 号线，约 55 分钟。" }, { icon: "♨", title: "餐饮推荐", text: "晚餐选择玉林路小馆，避开 19:00 排队高峰。" }], weather: "多云，22-28°C，夜间适合步行。" },
  { day: 2, title: "熊猫基地与宽窄巷子", route: "熊猫基地 → 文殊院素斋 → 宽窄巷子 → 奎星楼街", tips: [{ icon: "♟", title: "亲子友好", text: "熊猫基地建议乘坐观光车，节省体力。" }, { icon: "▱", title: "预约提醒", text: "门票建议提前一天购买，带好身份证。" }], weather: "晴间多云，注意防晒补水。" },
  { day: 3, title: "都江堰或人民公园慢游", route: "都江堰半日游 → 返程伴手礼 → 机场/高铁站", tips: [{ icon: "≋", title: "备选路线", text: "若不想远行，可改成人民公园茶馆与鹤鸣茶社。" }, { icon: "▤", title: "返程建议", text: "预留 2.5 小时前往机场，避免晚高峰。" }], weather: "午后可能阵雨，建议备伞。" }
];
const currentDay = computed(() => days.find((item) => item.day === activeDay.value) || days[0]);
onMounted(() => uni.hideTabBar && uni.hideTabBar());
function toast(title) { uni.showToast({ title, icon: "none" }); }
function goBack() { uni.navigateBack(); }
function saveTrip() { saved.value = true; toast("行程已保存"); }
</script>

<style scoped>
.page { min-height: 100vh; background: #f5f7fb; }
.phone-shell { width: 390px; max-width: 100vw; min-height: 100vh; margin: 0 auto; padding: 18px 18px 0; box-sizing: border-box; background: #f5f7fb; }
.with-nav { padding-bottom: 120px; }
.topbar { height: 42px; display: flex; align-items: center; justify-content: space-between; color: #0b5be7; }
.back { font-size: 32px; line-height: 1; }
.title { font-size: 20px; font-weight: 900; }
.avatar { width: 32px; height: 32px; border-radius: 999px; }
.map-card { position: relative; height: 188px; margin-top: 18px; border-radius: 18px; overflow: hidden; box-shadow: 0 8px 20px rgba(31,45,71,.1); }
.map-card image { width: 100%; height: 100%; }
.map-mask { position: absolute; inset: 0; background: linear-gradient(180deg, rgba(0,0,0,.05), rgba(0,0,0,.62)); }
.map-btn { position: absolute; right: 12px; top: 14px; z-index: 2; height: 34px; padding: 0 14px; border-radius: 999px; background: #eef5ff; color: #0b5be7; font-size: 13px; font-weight: 900; }
.destination { position: absolute; left: 18px; bottom: 40px; z-index: 2; color: #fff; font-size: 30px; font-weight: 900; }
.sub-destination { position: absolute; left: 18px; bottom: 18px; z-index: 2; color: rgba(255,255,255,.92); font-size: 13px; }
.summary { display: grid; grid-template-columns: repeat(3, 1fr); gap: 10px; margin-top: 14px; }
.summary view { padding: 12px; border-radius: 14px; background: #fff; box-shadow: 0 4px 14px rgba(31,45,71,.06); }
.summary text { display: block; color: #667085; font-size: 12px; }
.summary strong { display: block; margin-top: 5px; color: #101828; font-size: 14px; }
.section-title { display: block; margin: 24px 0 14px; color: #171b25; font-size: 22px; font-weight: 900; }
.day-tabs { display: grid; grid-template-columns: repeat(3, 1fr); gap: 10px; margin-bottom: 16px; }
.day-tabs text { height: 36px; border-radius: 18px; background: #fff; color: #5b6472; display: flex; align-items: center; justify-content: center; font-size: 13px; font-weight: 900; }
.day-tabs .active { background: #0b5be7; color: #fff; }
.day-head { display: flex; align-items: center; gap: 12px; }
.num { width: 24px; height: 24px; border-radius: 999px; background: #0b5be7; color: #fff; display: flex; align-items: center; justify-content: center; font-size: 13px; font-weight: 900; }
.day-title { color: #171b25; font-size: 19px; font-weight: 900; }
.line-wrap { display: flex; margin: 14px 0 20px 31px; }
.dash { width: 2px; margin-right: 18px; border-left: 2px dashed #0b5be7; }
.cards { flex: 1; display: flex; flex-direction: column; gap: 12px; }
.route,.mini,.weather,.budget { padding: 15px; border-radius: 14px; background: #fff; box-shadow: 0 4px 14px rgba(31,45,71,.06); }
.card-title { display: block; color: #0b5be7; font-size: 15px; font-weight: 900; }
.card-text { display: block; margin-top: 9px; color: #222936; font-size: 14px; line-height: 1.55; }
.mini-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 10px; }
.mini text { display: block; color: #17733c; font-size: 13px; font-weight: 900; }
.mini small { display: block; margin-top: 7px; color: #222936; font-size: 12px; line-height: 1.45; }
.weather { background: #dfe8ff; color: #1e3260; font-size: 13px; line-height: 1.5; }
.budget-title { color: #0b5be7; }
.budget-row,.total { display: flex; align-items: center; justify-content: space-between; color: #222936; font-size: 14px; }
.budget-row strong,.total strong { color: #0b5be7; font-size: 16px; }
.budget-row .green { color: #08843b; }
.bar { height: 8px; margin: 10px 0 16px; border-radius: 999px; background: #edf0f6; }
.fill { height: 100%; border-radius: 999px; width: 56%; }
.blue { background: #0b5be7; }
.green-bg { background: #07843b; width: 34%; }
.total { padding-top: 14px; border-top: 1px solid #d9deea; font-size: 16px; }
.total strong { font-size: 24px; }
.risk { color: #171b25; }
.risk-card { padding: 15px; border-radius: 14px; background: #fff0e8; border: 1px solid #ffd0b4; color: #a74700; font-size: 13px; line-height: 1.7; display: flex; flex-direction: column; gap: 8px; }
.save { position: fixed; left: 50%; bottom: 82px; z-index: 25; width: 350px; max-width: calc(100vw - 40px); height: 54px; transform: translateX(-50%); border-radius: 30px; background: #0b5be7; color: #fff; font-size: 17px; font-weight: 900; box-shadow: 0 12px 24px rgba(11,91,231,.25); }
</style>
