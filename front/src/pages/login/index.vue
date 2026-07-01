<template>
  <view class="login-page">
    <view class="phone-shell">
      <view class="hero">
        <view class="plane-badge">✈</view>
        <text class="brand">旅行助手</text>
        <text class="title">欢迎回来</text>
        <text class="subtitle">继续规划您的下一次国内旅程</text>
      </view>

      <view class="login-card">
        <view class="field-block"><text class="label">手机号</text><view class="input-box"><text class="icon">▯</text><input class="input" placeholder="请输入手机号码" v-model="phone" /></view></view>
        <view class="field-block"><view class="label-row"><text class="label">密码</text><text class="link" @tap="toast('已发送找回密码提示')">忘记密码</text></view><view class="input-box"><text class="icon">▣</text><input class="input" password placeholder="请输入登录密码" v-model="password" /><text class="eye">◉</text></view></view>
        <button class="primary" @tap="login">登录</button>
        <text class="code-login" @tap="toast('验证码登录功能演示')">验证码登录</text>
        <view class="divider"><view></view><text>快捷登录</text><view></view></view>
        <view class="social-row"><button class="social" @tap="toast('微信登录演示')"><text class="wechat">▰</text>微信</button><button class="social" @tap="toast('支付宝登录演示')"><text class="alipay">▦</text>支付宝</button></view>
      </view>
      <view class="register-link"><text>还没有账号？</text><text @tap="goRegister">立即注册</text></view>
    </view>
  </view>
</template>

<script setup>
import { ref } from "vue";
const phone = ref("");
const password = ref("");
const apiBase = "http://127.0.0.1:8080";
async function login() {
  if (!phone.value.trim() || !password.value.trim()) {
    toast("请输入手机号和密码");
    return;
  }
  try {
    const response = await fetch(`${apiBase}/api/v1/user/login`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ phone: phone.value.trim(), password: password.value })
    });
    const result = await response.json();
    if (!response.ok || result.code !== 0) {
      toast(result.msg || "登录失败");
      return;
    }
    uni.setStorageSync("travel_token", result.data.token);
    uni.setStorageSync("travel_user", result.data.user);
    uni.switchTab({ url: "/pages/index/index" });
  } catch (error) {
    toast("网关服务暂不可用");
  }
}
function goRegister() { uni.navigateTo({ url: "/pages/register/index" }); }
function toast(title) { uni.showToast({ title, icon: "none" }); }
</script>

<style scoped>
.login-page { min-height: 100vh; background: #e8f1fb; }
.phone-shell { width: 390px; max-width: 100vw; min-height: 100vh; margin: 0 auto; background: linear-gradient(180deg, #dcecff 0%, #f7f9ff 100%); overflow: hidden; }
.hero { height: 390px; padding-top: 58px; text-align: center; background: linear-gradient(180deg, rgba(218,235,255,.82), rgba(245,248,255,.78)), url("https://images.unsplash.com/photo-1500530855697-b586d89ba3ee?auto=format&fit=crop&w=900&q=80"); background-size: cover; background-position: center; }
.plane-badge { width: 84px; height: 84px; margin: 0 auto 24px; border-radius: 999px; background: #0b5be7; color: #fff; display: flex; align-items: center; justify-content: center; font-size: 40px; box-shadow: 0 18px 36px rgba(11,91,231,.25); }
.brand { display: block; color: #0758d9; font-size: 32px; font-weight: 900; }
.title { display: block; margin-top: 18px; color: #303746; font-size: 28px; font-weight: 900; }
.subtitle { display: block; margin-top: 8px; color: #7a8497; font-size: 15px; }
.login-card { margin: -32px 20px 0; padding: 32px 26px 28px; border-radius: 24px; background: rgba(255,255,255,.94); box-shadow: 0 22px 40px rgba(71,91,126,.14); }
.field-block { margin-bottom: 22px; }
.label-row { display: flex; align-items: center; justify-content: space-between; }
.label { display: block; margin-bottom: 10px; color: #3d4351; font-size: 14px; font-weight: 900; }
.link { color: #0758d9; font-size: 14px; font-weight: 900; }
.input-box { height: 58px; padding: 0 18px; border-radius: 16px; background: #f3f4f9; display: flex; align-items: center; }
.icon,.eye { color: #818899; font-size: 24px; }
.input { flex: 1; min-width: 0; margin-left: 14px; color: #1f2430; font-size: 18px; }
.primary { height: 60px; margin-top: 8px; border-radius: 30px; background: #0b5be7; color: #fff; font-size: 22px; font-weight: 900; box-shadow: 0 16px 28px rgba(11,91,231,.28); }
.code-login { display: block; margin-top: 24px; text-align: center; color: #7d8495; font-size: 15px; font-weight: 800; }
.divider { display: flex; align-items: center; gap: 16px; margin: 24px 0 18px; color: #757d90; font-size: 14px; font-weight: 800; }
.divider view { flex: 1; height: 1px; background: #c9d0df; }
.social-row { display: grid; grid-template-columns: 1fr 1fr; gap: 16px; }
.social { height: 52px; border-radius: 14px; border: 1px solid #c7cede; color: #222833; font-size: 17px; font-weight: 900; background: rgba(255,255,255,.78); }
.wechat { color: #09873b; margin-right: 10px; }
.alipay { color: #0b5be7; margin-right: 10px; }
.register-link { margin-top: 26px; text-align: center; color: #667085; font-size: 15px; }
.register-link text:last-child { margin-left: 10px; color: #0b5be7; font-weight: 900; }
</style>
