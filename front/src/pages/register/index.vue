<template>
  <view class="register-page">
    <view class="phone-shell">
      <view class="topbar">
        <text class="back" @tap="goBack">‹</text>
        <text class="top-title">创建账号</text>
        <text></text>
      </view>
      <view class="register-header">
        <text class="main-title">创建旅行档案</text>
        <text class="sub-title">用手机号同步你的行程、偏好和收藏</text>
      </view>
      <view class="register-card">
        <view class="field-block">
          <text class="label">手机号</text>
          <view class="input-box"><input class="input" placeholder="请输入手机号" v-model="phone" /></view>
        </view>
        <view class="field-block code-row">
          <view class="code-input"><input class="input" placeholder="请输入验证码" v-model="code" /></view>
          <button class="code-btn" :disabled="sendingCode" @tap="sendCode">{{ sendingCode ? "发送中" : "获取验证码" }}</button>
        </view>
        <view class="field-block">
          <text class="label">设置密码</text>
          <view class="input-box"><input class="input" password placeholder="建议包含字母与数字" v-model="password" /></view>
        </view>
      </view>
      <button class="primary-btn" @tap="register">注册</button>
    </view>
  </view>
</template>

<script setup>
import { ref } from "vue";
import { apiBase } from "../../utils/runtime.js";
const phone = ref("");
const code = ref("");
const password = ref("");
const sendingCode = ref(false);
function goBack() { uni.navigateBack(); }
async function sendCode() {
  const mobile = phone.value.trim();
  if (!/^1[3-9]\d{9}$/.test(mobile)) { toast("请输入正确手机号"); return; }
  if (sendingCode.value) return;
  sendingCode.value = true;
  try {
    const response = await fetch(`${apiBase}/api/v1/sms/register-code`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ phone: mobile })
    });
    const result = await response.json();
    if (!response.ok || result.code !== 0) { toast(result.msg || "验证码发送失败"); return; }
    if (result.data && result.data.code) { code.value = String(result.data.code); toast("开发环境验证码已填入"); return; }
    toast("验证码已发送");
  } catch (error) { toast("网关服务暂不可用"); }
  finally { sendingCode.value = false; }
}
async function register() {
  if (!phone.value.trim() || !password.value.trim()) { toast("请输入手机号和密码"); return; }
  if (!code.value.trim()) { toast("请输入验证码"); return; }
  try {
    const registerResponse = await fetch(`${apiBase}/api/v1/user/register`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ phone: phone.value.trim(), code: code.value.trim(), password: password.value, nickname: "", home_city: "", current_city: "" })
    });
    const registerResult = await registerResponse.json();
    if (!registerResponse.ok || registerResult.code !== 0) { toast(registerResult.msg || "注册失败"); return; }
    const loginResponse = await fetch(`${apiBase}/api/v1/user/login`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ phone: phone.value.trim(), password: password.value })
    });
    const loginResult = await loginResponse.json();
    if (!loginResponse.ok || loginResult.code !== 0) { toast(loginResult.msg || "注册成功，请返回登录"); return; }
    uni.setStorageSync("travel_token", loginResult.data.token);
    uni.setStorageSync("travel_user", loginResult.data.user);
    uni.switchTab({ url: "/pages/index/index" });
  } catch (error) { toast("网关服务暂不可用"); }
}
function toast(title) { uni.showToast({ title, icon: "none" }); }
</script>

<style scoped>
.register-page { min-height: 100vh; background: linear-gradient(180deg, #eaf1ff 0%, #f6f8fc 100%); }
.phone-shell { width: 100%; max-width: 430px; min-height: 100vh; margin: 0 auto; padding: calc(env(safe-area-inset-top) + 12px) 20px calc(env(safe-area-inset-bottom) + 34px); box-sizing: border-box; }
.topbar { min-height: 44px; display: flex; align-items: center; justify-content: space-between; }
.back { color: #0b5be7; font-size: 34px; line-height: 1; cursor: pointer; }
.top-title { color: #101828; font-size: 18px; font-weight: 800; }
.register-header { margin-top: 28px; margin-bottom: 28px; }
.main-title { display: block; color: #101828; font-size: 27px; font-weight: 800; }
.sub-title { display: block; margin-top: 6px; color: #667085; font-size: 14px; }
.register-card { padding: 24px 18px 6px; border-radius: 22px; background: #fff; box-shadow: 0 8px 28px rgba(31,45,71,.08); }
.field-block { margin-bottom: 20px; }
.label { display: block; margin-bottom: 8px; color: #344054; font-size: 14px; font-weight: 700; }
.input-box, .code-input { height: 52px; padding: 0 16px; border-radius: 14px; background: #f5f7fb; display: flex; align-items: center; border: 1px solid #eaecf0; }
.input { flex: 1; min-width: 0; color: #101828; font-size: 16px; }
.code-row { display: grid; grid-template-columns: minmax(0, 1fr) 112px; gap: 10px; align-items: end; }
.code-btn { height: 52px; border-radius: 14px; border: 1px solid #0b5be7; color: #0b5be7; font-size: 14px; font-weight: 800; background: #fff; }
.primary-btn { height: 50px; margin-top: 24px; border-radius: 25px; background: #0b5be7; color: #fff; font-size: 17px; font-weight: 800; box-shadow: 0 8px 20px rgba(11,91,231,.22); }

/* 注册页视觉收敛 */
.register-page { background: #dfe9e6; }
.phone-shell {
  background: linear-gradient(180deg, #fbfdfb 0%, #eef6f3 52%, #f8faf7 100%);
}
.back {
  color: #0f5f75;
  width: 40px;
  height: 40px;
  border-radius: 14px;
  background: rgba(255,255,252,.98);
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 8px 20px rgba(16,32,51,.08);
}
.top-title,
.main-title { color: #102033; }
.main-title { font-size: 30px; font-weight: 900; }
.sub-title { color: #667a78; }
.register-card {
  border-radius: 24px;
  background: rgba(255,255,252,.98);
  border: 1px solid rgba(224,232,229,.9);
  box-shadow: 0 16px 34px rgba(16,32,51,.09);
}
.label { color: #3d4f5d; }
.input-box,
.code-input {
  background: #f8fbf9;
  border-color: rgba(205,218,214,.9);
}
.code-btn {
  border-radius: 16px;
  border-color: #0f5f75;
  color: #0f5f75;
  background: #edf7f3;
}
.primary-btn {
  border-radius: 18px;
  background: #102033;
  color: #f8fbf9;
  box-shadow: 0 14px 28px rgba(16,32,51,.18);
}
</style>
