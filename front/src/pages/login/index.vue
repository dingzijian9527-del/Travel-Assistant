<template>
  <view class="login-page">
    <view class="phone-shell">
      <view class="login-header">
        <text class="brand">旅行助手</text>
        <text class="slogan">把目的地、时间和预算整理成路线</text>
      </view>
      <view class="login-card">
        <text class="card-title">欢迎回来</text>
        <view class="field-block">
          <text class="label">手机号</text>
          <view class="input-box"><input class="input" placeholder="请输入手机号码" v-model="phone" /></view>
        </view>
        <view class="field-block">
          <text class="label">密码</text>
          <view class="input-box"><input class="input" password placeholder="请输入登录密码" v-model="password" /></view>
        </view>
        <button class="primary-btn" @tap="login">登录</button>
        <view class="login-actions">
          <text class="link" @tap="goRegister">注册账号</text>
        </view>
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
  if (!phone.value.trim() || !password.value.trim()) { toast("请输入手机号和密码"); return; }
  try {
    const response = await fetch(`${apiBase}/api/v1/user/login`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ phone: phone.value.trim(), password: password.value })
    });
    const result = await response.json();
    if (!response.ok || result.code !== 0) { toast(result.msg || "登录失败"); return; }
    uni.setStorageSync("travel_token", result.data.token);
    uni.setStorageSync("travel_user", result.data.user);
    uni.switchTab({ url: "/pages/index/index" });
  } catch (error) { toast("网关服务暂不可用"); }
}
function goRegister() { uni.navigateTo({ url: "/pages/register/index" }); }
function toast(title) { uni.showToast({ title, icon: "none" }); }
</script>

<style scoped>
.login-page { min-height: 100vh; background: linear-gradient(180deg, #eaf1ff 0%, #f6f8fc 100%); }
.phone-shell { width: 100%; max-width: 430px; min-height: 100vh; margin: 0 auto; padding: calc(env(safe-area-inset-top) + 68px) 20px calc(env(safe-area-inset-bottom) + 28px); box-sizing: border-box; }
.login-header { text-align: center; margin-bottom: 34px; }
.brand { display: block; color: #0b5be7; font-size: 32px; font-weight: 900; letter-spacing: 0; }
.slogan { display: block; margin-top: 8px; color: #667085; font-size: 15px; }
.login-card { padding: 28px 20px 24px; border-radius: 22px; background: #fff; box-shadow: 0 8px 28px rgba(31,45,71,.08); }
.card-title { display: block; margin-bottom: 24px; color: #101828; font-size: 21px; font-weight: 800; text-align: center; }
.field-block { margin-bottom: 18px; }
.label { display: block; margin-bottom: 8px; color: #344054; font-size: 14px; font-weight: 700; }
.input-box { height: 52px; padding: 0 16px; border-radius: 14px; background: #f5f7fb; display: flex; align-items: center; border: 1px solid #eaecf0; }
.input { flex: 1; min-width: 0; color: #101828; font-size: 16px; }
.primary-btn { height: 50px; margin-top: 8px; border-radius: 25px; background: #0b5be7; color: #fff; font-size: 17px; font-weight: 800; box-shadow: 0 8px 20px rgba(11,91,231,.22); }
.login-actions { display: flex; justify-content: center; margin-top: 20px; }
.link { color: #0b5be7; font-size: 15px; font-weight: 700; }

/* 登录页视觉收敛 */
.login-page { background: #dfe9e6; }
.phone-shell {
  background: linear-gradient(180deg, #fbfdfb 0%, #eef6f3 52%, #f8faf7 100%);
}
.login-header {
  text-align: left;
  margin-bottom: 28px;
}
.brand {
  color: #102033;
  font-size: 34px;
  font-weight: 900;
}
.brand::after {
  content: "";
  display: block;
  width: 44px;
  height: 4px;
  margin-top: 12px;
  border-radius: 999px;
  background: #caa459;
}
.slogan {
  margin-top: 14px;
  color: #667a78;
}
.login-card {
  border-radius: 24px;
  background: rgba(255,255,252,.98);
  border: 1px solid rgba(224,232,229,.9);
  box-shadow: 0 16px 34px rgba(16,32,51,.09);
}
.card-title { color: #102033; text-align: left; }
.label { color: #3d4f5d; }
.input-box {
  background: #f8fbf9;
  border-color: rgba(205,218,214,.9);
}
.primary-btn {
  border-radius: 18px;
  background: #102033;
  color: #f8fbf9;
  box-shadow: 0 14px 28px rgba(16,32,51,.18);
}
.link {
  color: #0f5f75;
  font-weight: 800;
}
</style>
