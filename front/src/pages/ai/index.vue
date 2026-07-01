<template>
  <view class="page">
    <view class="phone-shell">
      <view class="chat-header">
        <view class="bot-icon">▣<text></text></view>
        <view class="bot-info"><text class="name">旅行智能体</text><text class="status">在线 · 擅长行程、攻略、酒店与美食</text></view>
        <image class="avatar" :src="avatar" mode="aspectFill" />
      </view>

      <scroll-view class="chat-scroll" scroll-y :scroll-into-view="lastId">
        <view v-for="item in messages" :key="item.id" :id="item.id" :class="['message-wrap', item.role]">
          <view :class="['bubble', item.role]">{{ item.text }}</view>
          <text class="time">{{ item.time }}</text>
          <view v-if="item.card" class="plan-card">
            <image class="plan-image" :src="item.card.image" mode="aspectFill" />
            <view class="rating">☆ {{ item.card.score }}</view>
            <view class="plan-body">
              <text class="plan-title">{{ item.card.title }}</text>
              <view class="tag-row"><text v-for="tag in item.card.tags" :key="tag">{{ tag }}</text></view>
              <text class="plan-desc">{{ item.card.desc }}</text>
            </view>
          </view>
        </view>
        <view class="service"><view class="service-icon hotel">▤</view><view><text>成都博舍酒店</text><small>太古里步行可达 · ¥1,280 起</small></view></view>
        <view class="service"><view class="service-icon food">♨</view><view><text>玉林路本地小馆</text><small>晚餐建议提前 2 小时排队</small></view></view>
      </scroll-view>

      <scroll-view class="quick-row" scroll-x show-scrollbar="false"><view class="quick-inner"><button v-for="item in prompts" :key="item" @tap="fill(item)">{{ item }}</button></view></scroll-view>
      <view class="input-bar"><text class="mic" @tap="toast('语音输入模拟')">♩</text><text class="plus" @tap="toast('可添加预算、日期、同行人')">＋</text><input v-model="message" placeholder="问我任何旅行问题..." @confirm="send" /><button class="send-btn" :disabled="isSending" @tap="send">发送</button></view>
      <BottomNav active="ai" />
    </view>
  </view>
</template>

<script setup>
import { nextTick, onMounted, ref } from "vue";
import BottomNav from "../../components/BottomNav.vue";
const message = ref("");
const lastId = ref("msg-1");
const isSending = ref(false);
const chatStreamUrl = "http://127.0.0.1:8080/api/v1/ai-stream";
const avatar = "https://images.unsplash.com/photo-1500648767791-00dcc994a43e?auto=format&fit=crop&w=120&q=80";
const prompts = ["帮我规划三天成都行程", "预算 3000 元去哪玩", "杭州亲子周末怎么安排", "推荐广州早茶路线"];
const messages = ref([
  { id: "msg-1", role: "ai", text: "你好！我是你的专属旅行智能体。无论你想去哪里，我都能帮你做国内出行路线、酒店建议、当地攻略和预算拆分。今天想聊聊哪里的旅行？", time: "10:30" },
  { id: "msg-2", role: "user", text: "我想计划下周去成都玩三天，预算比较充足，喜欢美食和悠闲的节奏。", time: "10:31" },
  { id: "msg-3", role: "ai", text: "成都非常适合慢节奏旅行！我先为你规划一个“锦官城悠享之旅”，节奏松弛，餐食会优先安排本地口碑店。", time: "10:31", card: { title: "成都 3 日悠享路线", score: "4.9", tags: ["第1天：太古里与锦里", "第2天：熊猫基地与宽窄巷子", "第3天：都江堰或茶馆慢游"], desc: "已为你预留美食、茶馆、夜游和弹性午休时间，适合不赶路的旅行方式。", image: "https://images.unsplash.com/photo-1519181245277-cffeb31da2e3?auto=format&fit=crop&w=760&q=80" } }
]);
onMounted(() => uni.hideTabBar && uni.hideTabBar());
function fill(text) { message.value = text; }
async function send() {
  if (isSending.value) return;
  const content = message.value.trim();
  if (!content) { toast("请输入内容"); return; }
  const token = uni.getStorageSync("travel_token");
  if (!token) {
    toast("请先登录后再使用智能体");
    uni.navigateTo({ url: "/pages/login/index" });
    return;
  }
  const userId = `msg-${Date.now()}`;
  const replyId = `${userId}-reply`;
  messages.value.push({ id: userId, role: "user", text: content, time: "现在" });
  messages.value.push({ id: replyId, role: "ai", text: "", time: "现在" });
  message.value = "";
  isSending.value = true;
  await nextTick();
  lastId.value = replyId;
  try {
    const response = await fetch(chatStreamUrl, {
      method: "POST",
      headers: { "Content-Type": "application/json", "Authorization": `Bearer ${token}` },
      body: JSON.stringify({ message: content })
    });
    if (!response.ok) throw new Error("智能体服务响应异常");
    const reply = messages.value.find((item) => item.id === replyId);
    if (!reply) return;
    if (!response.body) {
      reply.text = await response.text();
      return;
    }
    const reader = response.body.getReader();
    const decoder = new TextDecoder("utf-8");
    while (true) {
      const { value, done } = await reader.read();
      if (done) break;
      reply.text += decoder.decode(value, { stream: true });
      await nextTick();
      lastId.value = replyId;
    }
    reply.text += decoder.decode();
  } catch (error) {
    const reply = messages.value.find((item) => item.id === replyId);
    if (reply) reply.text = "智能体服务暂不可用，请确认网关已启动后再试。";
    toast("智能体服务暂不可用");
  } finally {
    isSending.value = false;
    await nextTick();
    lastId.value = replyId;
  }
}
function toast(title) { uni.showToast({ title, icon: "none" }); }
</script>

<style scoped>
.page { min-height: 100vh; background: #f5f7fb; }
.phone-shell { position: relative; width: 390px; max-width: 100vw; height: 100vh; margin: 0 auto; overflow: hidden; background: #f5f7fb; }
.chat-header { height: 82px; padding: 16px 18px 12px; display: flex; align-items: center; gap: 12px; border-bottom: 1px solid #e6e9f2; box-sizing: border-box; background: rgba(255,255,255,.92); }
.bot-icon { position: relative; width: 48px; height: 48px; border-radius: 999px; background: #0b5be7; color: #fff; display: flex; align-items: center; justify-content: center; font-size: 25px; font-weight: 900; }
.bot-icon text { position: absolute; right: 2px; bottom: 3px; width: 9px; height: 9px; border-radius: 999px; background: #11a346; border: 2px solid #f7f8ff; }
.bot-info { flex: 1; min-width: 0; }
.name { display: block; color: #0b5be7; font-size: 22px; font-weight: 900; }
.status { display: block; margin-top: 3px; color: #08843b; font-size: 12px; font-weight: 800; }
.avatar { width: 44px; height: 44px; border-radius: 999px; }
.chat-scroll { height: calc(100vh - 232px); padding: 18px; box-sizing: border-box; }
.message-wrap { margin-bottom: 18px; }
.message-wrap.user { text-align: right; }
.bubble { display: inline-block; max-width: 318px; padding: 16px; border-radius: 16px; font-size: 16px; line-height: 1.55; text-align: left; white-space: pre-wrap; box-sizing: border-box; box-shadow: 0 4px 12px rgba(32,44,70,.06); }
.bubble.ai { background: #fff; border: 1px solid #cfd6e6; color: #171b25; border-bottom-left-radius: 4px; }
.bubble.user { background: #0b5be7; color: #fff; border-bottom-right-radius: 4px; }
.time { display: block; margin-top: 8px; color: #747b8d; font-size: 12px; font-weight: 700; }
.plan-card { position: relative; margin-top: 14px; overflow: hidden; border-radius: 16px; background: #fff; border: 1px solid #d6dbe8; box-shadow: 0 6px 18px rgba(32,44,70,.08); }
.plan-image { width: 100%; height: 188px; }
.rating { position: absolute; right: 12px; top: 14px; padding: 7px 12px; border-radius: 999px; background: rgba(255,255,255,.9); color: #171b25; font-size: 14px; font-weight: 900; }
.plan-body { padding: 18px; }
.plan-title { display: block; color: #171b25; font-size: 22px; font-weight: 900; }
.tag-row { display: flex; flex-wrap: wrap; gap: 8px; margin-top: 12px; }
.tag-row text { padding: 6px 10px; border-radius: 999px; background: #e8f0ff; color: #0b5be7; font-size: 12px; font-weight: 800; }
.plan-desc { display: block; margin-top: 14px; color: #394153; font-size: 14px; line-height: 1.6; }
.service { margin-top: 12px; padding: 14px; border-radius: 14px; background: #fff; display: flex; align-items: center; gap: 14px; box-shadow: 0 4px 14px rgba(32,44,70,.06); }
.service-icon { width: 46px; height: 46px; border-radius: 12px; display: flex; align-items: center; justify-content: center; font-size: 23px; }
.hotel { background: #e8f0ff; color: #0b5be7; }
.food { background: #fff2df; color: #b45b00; }
.service text { display: block; color: #242b39; font-size: 15px; font-weight: 800; }
.service small { display: block; margin-top: 4px; color: #737b8d; font-size: 12px; }
.quick-row { position: fixed; left: 50%; bottom: 118px; z-index: 25; width: 390px; max-width: 100vw; transform: translateX(-50%); white-space: nowrap; padding: 0 18px; box-sizing: border-box; }
.quick-inner { display: inline-flex; gap: 10px; }
.quick-row button { height: 38px; padding: 0 15px; border-radius: 20px; background: #fff; border: 1px solid #d7deeb; color: #151b27; font-size: 13px; box-shadow: 0 4px 12px rgba(32,44,70,.05); }
.input-bar { position: fixed; left: 50%; bottom: 82px; z-index: 26; width: 354px; max-width: calc(100vw - 36px); height: 56px; transform: translateX(-50%); border-radius: 28px; background: #fff; border: 1px solid #cbd1e3; box-shadow: 0 10px 24px rgba(34,45,72,.12); display: flex; align-items: center; padding: 0 8px 0 14px; box-sizing: border-box; }
.mic,.plus { width: 26px; color: #697184; font-size: 23px; margin-right: 6px; text-align: center; flex: 0 0 26px; }
.input-bar input { flex: 1; min-width: 0; height: 100%; color: #252c3a; font-size: 14px; }
.send-btn { width: 58px; height: 40px; flex: 0 0 58px; border-radius: 20px; background: #0b5be7; color: #fff; font-size: 14px; font-weight: 900; box-shadow: 0 6px 14px rgba(11,91,231,.22); }
</style>

