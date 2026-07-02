<template>
  <view class="register-page">
    <view class="phone-shell">
      <view class="topbar"><text class="back" @tap="goBack">‹</text><text class="top-title">创建账号</text><text></text></view>
      <view class="headline"><text class="main-title">手机号注册</text></view>
      <view class="form">
        <view class="field-block"><text class="label">手机号</text><view class="input-box"><text class="icon">▯</text><input class="input" placeholder="请输入手机号" v-model="phone" /></view></view>
        <view class="field-block code-row"><view class="code-input"><text class="icon">▱</text><input class="input" placeholder="请输入验证码" v-model="code" /></view><button class="code-btn" :disabled="sendingCode" @tap="sendCode">{{ sendingCode ? "发送中" : "获取验证码" }}</button></view>
        <view class="field-block"><text class="label">设置密码</text><view class="input-box"><text class="icon">▣</text><input class="input" password placeholder="建议包含字母与数字" v-model="password" /><text class="eye">◌</text></view></view>
      </view>
      <button class="primary" @tap="register">注册</button>
    </view>
  </view>
</template>

<script setup>
import { ref } from "vue";
const phone = ref("");
const code = ref("");
const password = ref("");
const sendingCode = ref(false);
const apiBase = "http://127.0.0.1:8080";
function goBack() { uni.navigateBack(); }
async function sendCode() {
  const mobile = phone.value.trim();
  if (!/^1[3-9]\d{9}$/.test(mobile)) { toast("请输入正确手机号"); return; }
  if (sendingCode.value) { return; }
  sendingCode.value = true;
  try {
    const response = await fetch(`${apiBase}/api/v1/sms/register-code`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ phone: mobile })
    });
    const result = await response.json();
    if (!response.ok || result.code !== 0) {
      toast(result.msg || "验证码发送失败");
      return;
    }
    if (result.data && result.data.code) {
      code.value = String(result.data.code);
      toast("开发环境验证码已填入");
      return;
    }
    toast("验证码已发送");
  } catch (error) {
    toast("网关服务暂不可用");
  } finally {
    sendingCode.value = false;
  }
}
async function register() {
  if (!phone.value.trim() || !password.value.trim()) { toast("请输入手机号和密码"); return; }
  if (!code.value.trim()) { toast("请输入验证码"); return; }
  try {
    const registerResponse = await fetch(`${apiBase}/api/v1/user/register`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        phone: phone.value.trim(),
        code: code.value.trim(),
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
function toast(title) { uni.showToast({ title, icon: "none" }); }
</script>

<style scoped>
.register-page { min-height: 100vh; background: #f8f9ff; }
.phone-shell { width: 390px; max-width: 100vw; min-height: 100vh; margin: 0 auto; padding: 22px 20px 42px; box-sizing: border-box; background: #f8f9ff; }
.topbar { height: 42px; display: flex; align-items: center; justify-content: space-between; color: #0b5be7; }
.back { font-size: 34px; line-height: 1; }
.top-title { font-size: 20px; font-weight: 900; }
.headline { margin-top: 72px; }
.main-title { display: block; color: #161a23; font-size: 32px; line-height: 1.18; font-weight: 900; }
.form { margin-top: 34px; padding: 28px 22px 8px; border-radius: 18px; background: #fff; border: 1px solid #dce2ee; box-shadow: 0 14px 32px rgba(33,45,72,.08); }
.field-block { margin-bottom: 22px; }
.label { display: block; margin: 0 0 10px 6px; color: #303746; font-size: 15px; font-weight: 900; }
.input-box,.code-input { height: 58px; padding: 0 18px; border-radius: 16px; background: #f0f1f7; display: flex; align-items: center; }
.icon,.eye { color: #858b9b; font-size: 24px; }
.input { flex: 1; min-width: 0; margin-left: 14px; color: #1f2430; font-size: 17px; }
.code-row { display: grid; grid-template-columns: 1fr 118px; gap: 14px; align-items: end; }
.code-btn { height: 58px; border-radius: 16px; border: 1px solid #0b5be7; color: #0b5be7; font-size: 15px; font-weight: 900; background: #fff; }
.primary { height: 58px; margin-top: 28px; border-radius: 30px; background: #0b5be7; color: #fff; font-size: 20px; font-weight: 900; box-shadow: 0 16px 28px rgba(11,91,231,.24); }
</style>
