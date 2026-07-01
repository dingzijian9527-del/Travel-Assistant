<template>
  <view class="page">
    <view class="phone-shell with-nav">
      <view class="header">
        <view>
          <text class="hello">上午好，张先生</text>
          <text class="city">当前位置：成都 · 今日适合城市漫游</text>
        </view>
        <view class="header-actions"><text class="notice" @tap="toast('暂无新通知')">♧</text><image class="avatar" :src="avatar" mode="aspectFill" @tap="goProfile" /></view>
      </view>

      <view class="search">
        <text>⌕</text>
        <input placeholder="问智能体：成都周末怎么玩" v-model="keyword" @confirm="askAi" />
        <button @tap="askAi">问助手</button>
      </view>

      <view class="quick-tabs">
        <text v-for="item in quickTabs" :key="item" :class="['quick-tab', activeQuick === item ? 'active' : '']" @tap="activeQuick = item">{{ item }}</text>
      </view>

      <view class="planner">
        <view class="planner-head"><text>生成智能行程</text><small>参考真实国内旅行场景</small></view>
        <view class="form-grid">
          <view class="big cell" @tap="selectNextDestination"><text>目的地</text><strong>{{ selectedDestination.name }}</strong><small>{{ selectedDestination.theme }}</small></view>
          <view class="cell" @tap="changeDays"><text>日期/天数</text><strong>{{ dayCount }}天 {{ dateRange }}</strong></view>
          <view class="cell" @tap="togglePeople"><text>随行人员</text><strong>{{ peopleText }}</strong></view>
          <view class="cell" @tap="toggleBudget"><text>预算级别</text><strong>{{ budgetText }}</strong></view>
          <view class="cell"><text>兴趣偏好</text><view class="tag-line"><em v-for="tag in selectedDestination.tags" :key="tag">{{ tag }}</em></view></view>
        </view>
        <button class="generate" @tap="goAi">生成我的专属行程  →</button>
      </view>

      <view class="section-row"><text>智能行程预览</text><text class="more" @tap="goTrip">查看完整版</text></view>
      <view class="preview-card" @tap="goTrip">
        <image class="preview-img" :src="selectedDestination.image" mode="aspectFill" />
        <view class="cover"></view>
        <view class="weather">☼ {{ selectedDestination.weather }}</view>
        <text class="preview-title">{{ selectedDestination.title }}</text>
        <view class="days"><text v-for="day in days" :key="day" :class="['day', day === activeDay ? 'active' : '']" @tap.stop="activeDay = day">{{ day }}</text></view>
        <view class="facts"><text>{{ selectedDestination.traffic }}</text><text>{{ selectedDestination.food }}</text><text>{{ budgetAmount }}</text></view>
        <view class="tip">{{ selectedDestination.tip }}</view>
      </view>

      <view class="section-row"><text>国内热门推荐</text><view><text class="pill">周末</text><text class="pill light">亲子</text></view></view>
      <scroll-view class="recommend-scroll" scroll-x show-scrollbar="false">
        <view class="recommend-row">
          <view v-for="item in recs" :key="item.title" class="rec-card" @tap="selectDestination(item.name)">
            <image :src="item.image" mode="aspectFill" />
            <text class="score">{{ item.score }} 分</text>
            <text class="rec-title">{{ item.title }}</text>
            <text class="rec-desc">{{ item.desc }}</text>
            <text class="price">{{ item.price }}</text>
          </view>
        </view>
      </scroll-view>
      <BottomNav active="explore" />
    </view>
  </view>
</template>

<script setup>
import { computed, ref, onMounted } from "vue";
import BottomNav from "../../components/BottomNav.vue";
const keyword = ref("");
const avatar = "https://images.unsplash.com/photo-1500648767791-00dcc994a43e?auto=format&fit=crop&w=160&q=80";
const quickTabs = ["自由行", "周末游", "亲子游", "美食游"];
const activeQuick = ref("自由行");
const activeDay = ref("第1天");
const dayCount = ref(3);
const people = ref(2);
const budget = ref("舒适型");
const selectedIndex = ref(0);
const destinations = [
  { name: "成都", theme: "美食慢游 · 茶馆巷弄", title: "成都三日慢生活路线", weather: "22°C · 多云", traffic: "地铁+打车", food: "8处餐食", tip: "建议上午去熊猫基地，下午回到宽窄巷子喝盖碗茶，晚餐安排玉林路小馆。", tags: ["川菜", "茶馆", "熊猫"], image: "https://images.unsplash.com/photo-1519181245277-cffeb31da2e3?auto=format&fit=crop&w=760&q=80" },
  { name: "杭州", theme: "西湖骑行 · 江南园林", title: "杭州西湖与良渚文化路线", weather: "19°C · 晴", traffic: "地铁+骑行", food: "6处餐食", tip: "西湖建议清晨骑行，避开热门码头人流；下午安排良渚博物院更从容。", tags: ["西湖", "园林", "龙井"], image: "https://images.unsplash.com/photo-1500530855697-b586d89ba3ee?auto=format&fit=crop&w=760&q=80" },
  { name: "厦门", theme: "海边散步 · 老城小吃", title: "厦门鼓浪屿与环岛路路线", weather: "26°C · 晴", traffic: "轮渡+公交", food: "7处餐食", tip: "鼓浪屿船票建议提前一天预约，傍晚留给环岛路和沙坡尾更舒适。", tags: ["海岛", "小吃", "亲子"], image: "https://images.unsplash.com/photo-1507525428034-b723cf961d3e?auto=format&fit=crop&w=760&q=80" }
];
const recs = [
  { name: "成都", title: "成都三日慢游", desc: "熊猫基地、茶馆、川菜小馆", score: "4.9", price: "¥1680起", image: "https://images.unsplash.com/photo-1519181245277-cffeb31da2e3?auto=format&fit=crop&w=360&q=80" },
  { name: "杭州", title: "杭州西湖周末", desc: "西湖骑行、灵隐、龙井问茶", score: "4.8", price: "¥1299起", image: "https://images.unsplash.com/photo-1500530855697-b586d89ba3ee?auto=format&fit=crop&w=360&q=80" },
  { name: "厦门", title: "厦门海边假日", desc: "鼓浪屿、沙坡尾、环岛路", score: "4.7", price: "¥1890起", image: "https://images.unsplash.com/photo-1507525428034-b723cf961d3e?auto=format&fit=crop&w=360&q=80" }
];
const days = computed(() => Array.from({ length: dayCount.value }, (_, index) => `第${index + 1}天`));
const selectedDestination = computed(() => destinations[selectedIndex.value]);
const peopleText = computed(() => people.value === 1 ? "1位成人" : `${people.value}位成人`);
const budgetText = computed(() => budget.value);
const budgetAmount = computed(() => budget.value === "经济型" ? "约 ¥1800" : budget.value === "舒适型" ? "约 ¥3200" : "约 ¥5600");
const dateRange = computed(() => dayCount.value === 3 ? "06.28-06.30" : "07.03-07.07");
onMounted(() => uni.hideTabBar && uni.hideTabBar());
function toast(title) { uni.showToast({ title, icon: "none" }); }
function askAi() { uni.navigateTo({ url: `/pages/ai/index?question=${encodeURIComponent(keyword.value || '帮我规划一次国内旅行')}` }); }
function goAi() { uni.navigateTo({ url: "/pages/ai/index" }); }
function goTrip() { uni.switchTab({ url: "/pages/trip/index" }); }
function goProfile() { uni.switchTab({ url: "/pages/profile/index" }); }
function selectNextDestination() { selectedIndex.value = (selectedIndex.value + 1) % destinations.length; activeDay.value = "第1天"; }
function selectDestination(name) { const index = destinations.findIndex((item) => item.name === name); if (index >= 0) selectedIndex.value = index; }
function changeDays() { dayCount.value = dayCount.value === 3 ? 5 : 3; activeDay.value = "第1天"; }
function togglePeople() { people.value = people.value === 2 ? 4 : 2; }
function toggleBudget() { budget.value = budget.value === "舒适型" ? "经济型" : budget.value === "经济型" ? "高端型" : "舒适型"; }
</script>

<style scoped>
.page { min-height: 100vh; background: #f5f7fb; }
.phone-shell { width: 390px; max-width: 100vw; min-height: 100vh; margin: 0 auto; padding: 18px 16px 0; box-sizing: border-box; background: #f5f7fb; }
.with-nav { padding-bottom: 104px; }
.header { display: flex; align-items: center; justify-content: space-between; }
.hello { display: block; color: #101828; font-size: 24px; font-weight: 900; }
.city { display: block; margin-top: 6px; color: #667085; font-size: 12px; }
.header-actions { display: flex; align-items: center; gap: 14px; }
.notice { width: 34px; height: 34px; border-radius: 999px; background: #fff; color: #344054; display: flex; align-items: center; justify-content: center; box-shadow: 0 4px 12px rgba(31,45,71,.08); }
.avatar { width: 38px; height: 38px; border-radius: 999px; border: 2px solid #fff; box-shadow: 0 4px 12px rgba(31,45,71,.12); }
.search { height: 48px; margin-top: 18px; padding: 0 10px 0 14px; border-radius: 14px; background: #fff; display: flex; align-items: center; gap: 10px; color: #344054; box-shadow: 0 6px 18px rgba(31,45,71,.07); }
.search input { flex: 1; min-width: 0; font-size: 14px; }
.search button { width: 72px; height: 34px; border-radius: 17px; background: #0b5be7; color: #fff; font-size: 13px; font-weight: 900; }
.quick-tabs { display: flex; gap: 10px; margin-top: 14px; }
.quick-tab { padding: 8px 12px; border-radius: 999px; background: #fff; color: #526071; font-size: 13px; font-weight: 800; }
.quick-tab.active { background: #e8f0ff; color: #0b5be7; }
.planner { margin-top: 18px; padding: 16px; border-radius: 20px; background: linear-gradient(135deg, #0964dc 0%, #23b6ff 54%, #56e583 100%); color: #0d1830; box-shadow: 0 12px 26px rgba(11,91,231,.18); }
.planner-head { display: flex; justify-content: space-between; align-items: flex-end; margin-bottom: 14px; }
.planner-head text { color: #fff; font-size: 21px; font-weight: 900; }
.planner-head small { color: rgba(255,255,255,.86); font-size: 11px; }
.form-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 10px; }
.cell { min-height: 76px; padding: 12px; border-radius: 12px; background: rgba(255,255,255,.92); box-shadow: inset 0 0 0 1px rgba(255,255,255,.55); }
.big { grid-column: 1 / 3; }
.cell text { display: block; color: #5b6472; font-size: 12px; font-weight: 800; }
.cell strong { display: block; margin-top: 8px; color: #151a25; font-size: 18px; }
.cell small { display: block; margin-top: 4px; color: #667085; font-size: 12px; }
.tag-line { display: flex; flex-wrap: wrap; gap: 5px; }
em { display: inline-flex; margin-top: 10px; padding: 4px 8px; border-radius: 999px; background: #e8f0ff; color: #0b5be7; font-size: 12px; font-style: normal; }
.generate { height: 52px; margin-top: 14px; border-radius: 14px; background: #0b5be7; color: #fff; font-size: 18px; font-weight: 900; box-shadow: 0 10px 20px rgba(11,91,231,.24); }
.section-row { display: flex; align-items: center; justify-content: space-between; margin: 24px 2px 14px; color: #111722; font-size: 22px; font-weight: 900; }
.more { color: #0b5be7; font-size: 13px; font-weight: 800; }
.preview-card { position: relative; overflow: hidden; border-radius: 18px; background: #fff; box-shadow: 0 8px 24px rgba(31,45,71,.09); }
.preview-img { width: 100%; height: 164px; }
.cover { position: absolute; left: 0; right: 0; top: 0; height: 164px; background: linear-gradient(180deg, rgba(0,0,0,.1), rgba(0,0,0,.62)); }
.weather { position: absolute; left: 16px; top: 106px; color: #fff; font-size: 13px; z-index: 1; }
.preview-title { position: absolute; left: 16px; top: 126px; color: #fff; z-index: 1; font-size: 22px; font-weight: 900; }
.days { display: grid; grid-template-columns: repeat(5, 1fr); gap: 8px; padding: 14px 16px; }
.day { height: 42px; border-radius: 9px; background: #f0f2f8; display: flex; align-items: center; justify-content: center; color: #192033; font-size: 12px; font-weight: 900; }
.day.active { background: #eaf2ff; border: 1px solid #0b5be7; color: #0b5be7; }
.facts { display: grid; grid-template-columns: 1fr 1fr 1fr; gap: 6px; padding: 0 16px 12px; color: #303746; font-size: 12px; }
.tip { margin: 0 16px 16px; padding: 13px; border-radius: 12px; background: #e8fff0; color: #174125; font-size: 13px; line-height: 1.5; }
.pill { display: inline-flex; margin-left: 8px; padding: 5px 11px; border-radius: 999px; background: #e8f0ff; color: #0b5be7; font-size: 12px; }
.pill.light { color: #4b5363; background: #eef0f6; }
.recommend-scroll { white-space: nowrap; }
.recommend-row { display: inline-flex; gap: 14px; }
.rec-card { position: relative; width: 210px; height: 224px; border-radius: 16px; overflow: hidden; background: #fff; box-shadow: 0 6px 18px rgba(31,45,71,.08); }
.rec-card image { width: 100%; height: 118px; background: #e6ebf5; }
.score { position: absolute; right: 10px; top: 10px; padding: 4px 8px; border-radius: 999px; background: rgba(255,255,255,.9); color: #0b5be7; font-size: 12px; font-weight: 900; }
.rec-title { display: block; margin: 14px 14px 5px; color: #151a24; font-size: 18px; font-weight: 900; }
.rec-desc { display: block; margin: 0 14px; color: #4d5565; font-size: 12px; line-height: 1.35; }
.price { display: block; margin: 9px 14px 0; color: #ff6b00; font-size: 16px; font-weight: 900; }
</style>
