<template>
  <view class="page">
    <view class="phone-shell with-nav">
      <view class="topbar"><text class="logo">✈</text><text class="title">发现灵感</text><text class="search-icon" @tap="runSearch">⌕</text></view>
      <view class="search"><text>⌕</text><input v-model="keyword" placeholder="搜索目的地、攻略、景点" @confirm="runSearch" /><button @tap="runSearch">搜索</button></view>
      <scroll-view class="category-scroll" scroll-x show-scrollbar="false"><view class="category-row"><view v-for="item in categories" :key="item.name" :class="['category', activeCategory === item.name ? 'active' : '']" @tap="activeCategory = item.name"><view>{{ item.icon }}</view><text>{{ item.name }}</text></view></view></scroll-view>
      <view class="banner" @tap="toast('已加入收藏')"><image :src="banner.image" mode="aspectFill" /><view class="mask"></view><view class="copy"><text class="tag">当季推荐</text><text class="banner-title">{{ banner.title }}</text><text class="desc">{{ banner.desc }}</text></view></view>
      <view class="section-row"><text>{{ activeCategory }}推荐</text><text class="more" @tap="toast('已为你刷新推荐')">换一批 ›</text></view>
      <view class="city-grid"><view v-for="city in filteredCities" :key="city.name" class="city" @tap="openCity(city)"><view class="img-wrap"><image :src="city.image" mode="aspectFill" /><text>★ {{ city.rating }}</text></view><view class="city-body"><strong>{{ city.name }}</strong><small>{{ city.feature }}</small><small>⌖ 距您 {{ city.distance }}</small></view></view></view>
      <view class="section-row"><text>热门攻略</text><text class="more">{{ keywordLabel }}</text></view>
      <view class="guide-list"><view v-for="guide in filteredGuides" :key="guide.title" class="guide" @tap="toast('打开攻略：' + guide.title)"><image :src="guide.image" mode="aspectFill" /><view><strong>{{ guide.title }}</strong><view class="guide-tags"><text v-for="tag in guide.tags" :key="tag">{{ tag }}</text></view><view class="meta"><text>{{ guide.author }}</text><text>◎ {{ guide.views }}</text></view></view></view></view>
      <BottomNav active="explore" />
    </view>
  </view>
</template>

<script setup>
import { computed, ref, onMounted } from "vue";
import BottomNav from "../../components/BottomNav.vue";
const keyword = ref("");
const searched = ref("");
const activeCategory = ref("城市漫游");
const categories = [{ name: "城市漫游", icon: "▥" }, { name: "山水度假", icon: "≋" }, { name: "亲子周末", icon: "♟" }, { name: "美食打卡", icon: "♨" }, { name: "古镇人文", icon: "▣" }];
const banner = { title: "川西秋色自驾线", desc: "成都出发，四姑娘山、新都桥、康定三日环线", image: "https://images.unsplash.com/photo-1500530855697-b586d89ba3ee?auto=format&fit=crop&w=760&q=80" };
const cities = [
  { name: "成都 · 慢生活", category: "城市漫游", rating: "4.9", distance: "0km", feature: "茶馆、川菜、熊猫基地", image: "https://images.unsplash.com/photo-1519181245277-cffeb31da2e3?auto=format&fit=crop&w=420&q=80" },
  { name: "杭州 · 西湖骑行", category: "城市漫游", rating: "4.8", distance: "1600km", feature: "西湖、灵隐、龙井村", image: "https://images.unsplash.com/photo-1500530855697-b586d89ba3ee?auto=format&fit=crop&w=420&q=80" },
  { name: "桂林 · 山水画卷", category: "山水度假", rating: "4.9", distance: "980km", feature: "漓江、阳朔、遇龙河", image: "https://images.unsplash.com/photo-1500530855697-b586d89ba3ee?auto=format&fit=crop&w=420&q=80" },
  { name: "厦门 · 海边亲子", category: "亲子周末", rating: "4.7", distance: "1900km", feature: "鼓浪屿、环岛路、沙坡尾", image: "https://images.unsplash.com/photo-1507525428034-b723cf961d3e?auto=format&fit=crop&w=420&q=80" },
  { name: "广州 · 早茶地图", category: "美食打卡", rating: "4.8", distance: "1500km", feature: "早茶、骑楼、珠江夜游", image: "https://images.unsplash.com/photo-1525755662778-989d0524087e?auto=format&fit=crop&w=420&q=80" },
  { name: "苏州 · 园林人文", category: "古镇人文", rating: "4.8", distance: "1700km", feature: "园林、评弹、平江路", image: "https://images.unsplash.com/photo-1500530855697-b586d89ba3ee?auto=format&fit=crop&w=420&q=80" }
];
const guides = [
  { title: "成都三日不赶路：从熊猫基地到玉林小馆", author: "锦城旅行家", views: "2.6w", category: "城市漫游", tags: ["成都", "美食"], image: "https://images.unsplash.com/photo-1519181245277-cffeb31da2e3?auto=format&fit=crop&w=180&q=80" },
  { title: "杭州周末避开人潮：西湖清晨骑行路线", author: "江南慢游", views: "1.8w", category: "城市漫游", tags: ["杭州", "骑行"], image: "https://images.unsplash.com/photo-1500530855697-b586d89ba3ee?auto=format&fit=crop&w=180&q=80" },
  { title: "广州早茶十家本地人常去茶楼", author: "行走的小七", views: "3.1w", category: "美食打卡", tags: ["广州", "早茶"], image: "https://images.unsplash.com/photo-1525755662778-989d0524087e?auto=format&fit=crop&w=180&q=80" },
  { title: "厦门亲子三天：海边、轮渡和老城小吃", author: "周末亲子指南", views: "1.4w", category: "亲子周末", tags: ["厦门", "亲子"], image: "https://images.unsplash.com/photo-1507525428034-b723cf961d3e?auto=format&fit=crop&w=180&q=80" }
];
const filteredCities = computed(() => cities.filter((city) => city.category === activeCategory.value || city.name.includes(searched.value)).slice(0, 4));
const filteredGuides = computed(() => guides.filter((guide) => guide.category === activeCategory.value || guide.title.includes(searched.value)).slice(0, 3));
const keywordLabel = computed(() => searched.value ? `搜索：${searched.value}` : "真实体验");
onMounted(() => uni.hideTabBar && uni.hideTabBar());
function runSearch() { searched.value = keyword.value.trim(); if (searched.value) toast(`已搜索：${searched.value}`); }
function openCity(city) { keyword.value = city.name.split(" · ")[0]; toast(`已选择${city.name}`); }
function toast(title) { uni.showToast({ title, icon: "none" }); }
</script>

<style scoped>
.page { min-height: 100vh; background: #f5f7fb; }
.phone-shell { width: 390px; max-width: 100vw; min-height: 100vh; margin: 0 auto; padding: 18px 16px 0; box-sizing: border-box; background: #f5f7fb; }
.with-nav { padding-bottom: 104px; }
.topbar { height: 42px; display: flex; align-items: center; justify-content: space-between; color: #0b5be7; }
.logo,.search-icon { font-size: 24px; font-weight: 900; }
.title { color: #0b5be7; font-size: 20px; font-weight: 900; }
.search { height: 48px; margin-top: 14px; padding: 0 10px 0 16px; border-radius: 14px; background: #fff; display: flex; align-items: center; gap: 10px; color: #3a4354; box-shadow: 0 6px 18px rgba(31,45,71,.07); }
.search input { flex: 1; min-width: 0; font-size: 14px; }
.search button { width: 62px; height: 34px; border-radius: 17px; background: #0b5be7; color: #fff; font-weight: 900; }
.category-scroll { margin-top: 22px; white-space: nowrap; }
.category-row { display: inline-flex; gap: 14px; }
.category { width: 72px; text-align: center; color: #303746; }
.category view { width: 56px; height: 56px; margin: 0 auto 8px; border-radius: 999px; background: #fff; color: #0b5be7; display: flex; align-items: center; justify-content: center; font-size: 25px; box-shadow: 0 4px 14px rgba(31,45,71,.06); }
.category text { font-size: 12px; }
.category.active view { background: #0b5be7; color: #fff; box-shadow: 0 8px 18px rgba(11,91,231,.24); }
.category.active text { color: #0b5be7; font-weight: 900; }
.banner { position: relative; height: 162px; margin-top: 24px; border-radius: 18px; overflow: hidden; box-shadow: 0 8px 20px rgba(31,45,71,.1); }
.banner image { width: 100%; height: 100%; }
.mask { position: absolute; inset: 0; background: linear-gradient(90deg, rgba(0,0,0,.6), rgba(0,0,0,.05)); }
.copy { position: absolute; left: 18px; bottom: 18px; color: #fff; }
.tag { display: inline-flex; padding: 6px 10px; border-radius: 999px; background: #0b5be7; font-size: 12px; font-weight: 900; }
.banner-title { display: block; margin-top: 10px; font-size: 22px; font-weight: 900; }
.desc { display: block; margin-top: 4px; max-width: 270px; font-size: 13px; line-height: 1.45; }
.section-row { display: flex; align-items: center; justify-content: space-between; margin: 24px 2px 14px; color: #171b25; font-size: 20px; font-weight: 900; }
.more { color: #0b5be7; font-size: 13px; font-weight: 800; }
.city-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 14px; }
.city { overflow: hidden; border-radius: 16px; background: #fff; box-shadow: 0 6px 18px rgba(31,45,71,.07); }
.img-wrap { position: relative; height: 132px; }
.img-wrap image { width: 100%; height: 100%; }
.img-wrap text { position: absolute; right: 9px; top: 9px; padding: 5px 10px; border-radius: 999px; background: rgba(255,255,255,.92); color: #0b5be7; font-size: 12px; font-weight: 900; }
.city-body { padding: 13px; }
.city-body strong { display: block; color: #171b25; font-size: 16px; }
.city-body small { display: block; margin-top: 7px; color: #576071; font-size: 12px; }
.guide-list { display: flex; flex-direction: column; gap: 14px; }
.guide { display: flex; gap: 14px; padding: 14px; border-radius: 16px; background: #fff; box-shadow: 0 6px 18px rgba(31,45,71,.07); }
.guide image { width: 86px; height: 86px; border-radius: 10px; flex: 0 0 86px; }
.guide strong { display: block; color: #171b25; font-size: 15px; line-height: 1.45; }
.guide-tags { display: flex; gap: 6px; margin-top: 8px; }
.guide-tags text { padding: 3px 7px; border-radius: 999px; background: #e8f0ff; color: #0b5be7; font-size: 11px; }
.meta { display: flex; justify-content: space-between; margin-top: 10px; color: #687083; font-size: 12px; }
</style>
