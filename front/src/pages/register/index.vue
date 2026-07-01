<template>
  <view class="register-page">
    <view class="phone-shell">
      <view class="topbar"><text class="back" @tap="goBack">‹</text><text class="top-title">创建账号</text><text></text></view>
      <view class="headline"><text class="main-title">开启您的精彩旅程</text><text class="sub-title">加入旅行助手，保存偏好并生成更懂您的国内行程</text></view>
      <view class="form">
        <view class="field-block"><text class="label">手机号</text><view class="input-box"><text class="icon">▯</text><input class="input" placeholder="请输入手机号" v-model="phone" /></view></view>
        <view class="field-block code-row"><view class="code-input"><text class="icon">▱</text><input class="input" placeholder="请输入验证码" v-model="code" /></view><button class="code-btn" @tap="toast('验证码已发送')">获取验证码</button></view>
        <view class="field-block"><text class="label">设置密码</text><view class="input-box"><text class="icon">▣</text><input class="input" password placeholder="建议包含字母与数字" v-model="password" /><text class="eye">◌</text></view></view>
        <view class="agree" @tap="agreed = !agreed"><view :class="['checkbox', agreed ? 'checked' : '']">{{ agreed ? '✓' : '' }}</view><text>已阅读并同意《用户协议》与《隐私政策》，授权旅行助手管理我的行程信息。</text></view>
      </view>
      <view class="prefs"><view class="pref-head"><text>旅行偏好（可选）</text><text class="multi">多选</text></view><view class="pref-grid"><view v-for="item in prefs" :key="item" :class="['pref', selectedPrefs.includes(item) ? 'active' : '']" @tap="togglePref(item)">{{ item }}</view></view></view>
      <button class="primary" @tap="register">注册并开始规划  →</button>
      <view class="login-link"><text>已有账号？</text><text class="blue" @tap="goLogin">去登录</text></view>
    </view>
  </view>
</template>

<script setup>
import { ref } from "vue";
const phone = ref("");
const code = ref("");
const password = ref("");
const agreed = ref(false);
const apiBase = "http://127.0.0.1:8080";
const prefs = ["自然风光", "美食探索", "亲子游", "城市漫步", "小众秘境", "预算优先"];
const selectedPrefs = ref(["美食探索", "城市漫步"]);
function goBack() { uni.navigateBack(); }
function goLogin() { uni.navigateTo({ url: "/pages/login/index" }); }
async function register() {
  if (!agreed.value) { toast("请先同意协议"); return; }
  if (!phone.value.trim() || !password.value.trim()) { toast("请输入手机号和密码"); return; }
  try {
    const registerResponse = await fetch(`${apiBase}/api/v1/user/register`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        phone: phone.value.trim(),
        password: password.value,
        nickname: phone.value.trim() ? `旅行者${phone.value.trim().slice(-4)}` : "",
        home_city: "",
        current_city: ""
      })
    });
    const registerResult = await registerResponse.json();
    if (!registerResponse.ok || registerResult.code !== 0) {
      toast(registerResult.msg || "注册失败");
      return;
    }
    const loginResponse = await fetch(`${apiBase}/api/v1/user/login`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ phone: phone.value.trim(), password: password.value })
    });
    const loginResult = await loginResponse.json();
    if (!loginResponse.ok || loginResult.code !== 0) {
      toast(loginResult.msg || "注册成功，请返回登录");
      return;
    }
    uni.setStorageSync("travel_token", loginResult.data.token);
    uni.setStorageSync("travel_user", loginResult.data.user);
    uni.switchTab({ url: "/pages/index/index" });
  } catch (error) {
    toast("网关服务暂不可用");
  }
}
function togglePref(item) { selectedPrefs.value = selectedPrefs.value.includes(item) ? selectedPrefs.value.filter((pref) => pref !== item) : [...selectedPrefs.value, item]; }
function toast(title) { uni.showToast({ title, icon: "none" }); }
</script>

<style scoped>
.register-page { min-height: 100vh; background: #f8f9ff; }
.phone-shell { width: 390px; max-width: 100vw; min-height: 100vh; margin: 0 auto; padding: 22px 20px 42px; box-sizing: border-box; background: #f8f9ff; }
.topbar { height: 42px; display: flex; align-items: center; justify-content: space-between; color: #0b5be7; }
.back { font-size: 34px; line-height: 1; }
.top-title { font-size: 20px; font-weight: 900; }
.headline { margin-top: 48px; }
.main-title { display: block; color: #161a23; font-size: 32px; line-height: 1.18; font-weight: 900; }
.sub-title { display: block; margin-top: 14px; color: #646b7a; font-size: 18px; line-height: 1.45; }
.form { margin-top: 32px; }
.field-block { margin-bottom: 22px; }
.label { display: block; margin: 0 0 10px 6px; color: #303746; font-size: 15px; font-weight: 900; }
.input-box,.code-input { height: 58px; padding: 0 18px; border-radius: 16px; background: #f0f1f7; display: flex; align-items: center; }
.icon,.eye { color: #858b9b; font-size: 24px; }
.input { flex: 1; min-width: 0; margin-left: 14px; color: #1f2430; font-size: 17px; }
.code-row { display: grid; grid-template-columns: 1fr 118px; gap: 14px; align-items: end; }
.code-btn { height: 58px; border-radius: 16px; border: 1px solid #0b5be7; color: #0b5be7; font-size: 15px; font-weight: 900; background: #fff; }
.agree { display: flex; align-items: flex-start; gap: 12px; margin-top: 8px; color: #606879; font-size: 14px; line-height: 1.42; }
.checkbox { width: 20px; height: 20px; flex: 0 0 20px; border-radius: 6px; background: #e3e5ec; color: #fff; display: flex; align-items: center; justify-content: center; font-size: 14px; }
.checkbox.checked { background: #0b5be7; }
.prefs { margin-top: 24px; }
.pref-head { display: flex; align-items: center; justify-content: space-between; color: #171b24; font-size: 22px; font-weight: 900; }
.multi { color: #606879; font-size: 14px; width: 24px; line-height: 1.2; }
.pref-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 14px; margin-top: 18px; }
.pref { min-height: 44px; padding: 0 10px; border-radius: 22px; border: 1px solid #c9cedd; color: #626978; display: flex; align-items: center; justify-content: center; font-size: 16px; font-weight: 800; background: #fbfcff; }
.pref.active { border-color: #0b5be7; background: #e8f0ff; color: #0b5be7; }
.primary { height: 58px; margin-top: 44px; border-radius: 30px; background: #0b5be7; color: #fff; font-size: 20px; font-weight: 900; box-shadow: 0 16px 28px rgba(11,91,231,.24); }
.login-link { margin-top: 24px; text-align: center; color: #626978; font-size: 18px; font-weight: 700; }
.blue { margin-left: 16px; color: #0b5be7; font-size: 22px; font-weight: 900; }
</style>
