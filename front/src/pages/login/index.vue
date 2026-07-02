<template>
  <view class="login-page">
    <view class="phone-shell">
      <view class="login-card">
        <text class="brand">旅行助手</text>
        <text class="title">手机号密码登录</text>
        <view class="field-block"><text class="label">手机号</text><view class="input-box"><input class="input" placeholder="请输入手机号码" v-model="phone" /></view></view>
        <view class="field-block"><text class="label">密码</text><view class="input-box"><input class="input" password placeholder="请输入登录密码" v-model="password" /></view></view>
        <button class="primary" @tap="login">登录</button>
        <button class="secondary" @tap="goRegister">注册账号</button>
      </view>
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
.login-page { min-height: 100vh; background: #f5f7fb; }
.phone-shell { width: 390px; max-width: 100vw; min-height: 100vh; margin: 0 auto; padding: 108px 20px 32px; box-sizing: border-box; background: #f5f7fb; overflow: hidden; }
.login-card { padding: 34px 26px 30px; border-radius: 18px; background: #fff; border: 1px solid #dce2ee; box-shadow: 0 14px 32px rgba(33,45,72,.08); }
.brand { display: block; color: #0b5be7; font-size: 30px; font-weight: 900; text-align: center; }
.title { display: block; margin: 12px 0 34px; color: #303746; font-size: 20px; font-weight: 900; text-align: center; }
.field-block { margin-bottom: 22px; }
.label { display: block; margin-bottom: 10px; color: #3d4351; font-size: 14px; font-weight: 900; }
.input-box { height: 58px; padding: 0 18px; border-radius: 16px; background: #f3f4f9; display: flex; align-items: center; }
.input { flex: 1; min-width: 0; color: #1f2430; font-size: 18px; }
.primary { height: 60px; margin-top: 8px; border-radius: 30px; background: #0b5be7; color: #fff; font-size: 22px; font-weight: 900; box-shadow: 0 16px 28px rgba(11,91,231,.28); }
.secondary { height: 56px; margin-top: 16px; border-radius: 28px; border: 1px solid #0b5be7; background: #fff; color: #0b5be7; font-size: 18px; font-weight: 900; }
</style>
