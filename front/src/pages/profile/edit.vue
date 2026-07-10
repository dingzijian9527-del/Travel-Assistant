<template>
  <view class="page">
    <view class="phone-shell">
      <view class="topbar">
        <text class="back" @tap="goBack">‹</text>
        <text class="top-title">编辑资料</text>
        <text class="top-action" @tap="save">保存</text>
      </view>

      <view class="hero-card">
        <text class="hero-title">完善你的旅行档案</text>
        <text class="hero-desc">昵称、头像和常用城市会同步到首页、我的页和会员码页。</text>
      </view>

      <view class="form-card">
        <view class="avatar-section">
          <text class="label">头像</text>
          <view class="avatar-row">
            <image v-if="avatarPreview" class="avatar-preview" :src="avatarPreview" mode="aspectFill" />
            <view v-else class="avatar-placeholder">头像</view>
            <view class="avatar-actions">
              <button class="upload-btn" @tap="pickAvatar">{{ uploadingAvatar ? "上传中..." : "上传头像" }}</button>
              <text class="upload-tip">支持 jpg、png、webp，上传后会自动回填。</text>
            </view>
          </view>
        </view>

        <view class="field-block">
          <text class="label">昵称</text>
          <view class="input-box"><input class="input" v-model="form.nickname" placeholder="请输入昵称" maxlength="20" /></view>
        </view>

        <view class="field-grid">
          <view class="field-block half">
            <text class="label">家乡</text>
            <view class="input-box"><input class="input" v-model="form.homeCity" placeholder="请输入家乡" maxlength="20" /></view>
          </view>
          <view class="field-block half">
            <text class="label">常住地</text>
            <view class="input-box"><input class="input" v-model="form.currentCity" placeholder="请输入常住地" maxlength="20" /></view>
          </view>
        </view>
      </view>

      <button class="save-btn" @tap="save">保存资料</button>
    </view>
  </view>
</template>

<script setup>
import { computed, onMounted, ref } from "vue";
import { apiBase, getStoredUser, getToken, requestJSON, setStoredUser } from "../../utils/api.js";
import { applyProfileToStoredUser, buildProfilePayload, normalizeProfileForm } from "../../utils/profileForm.js";

const form = ref(normalizeProfileForm(getStoredUser()));
const saving = ref(false);
const uploadingAvatar = ref(false);
const avatarPreview = computed(() => form.value.avatarUrl || "");

onMounted(() => {
  uni.hideTabBar && uni.hideTabBar();
  loadProfile();
});

async function loadProfile() {
  try {
    const profile = await requestJSON("/api/v1/user/profile");
    form.value = normalizeProfileForm(profile);
    setStoredUser(applyProfileToStoredUser(getStoredUser(), profile));
  } catch (error) {
    console.warn("读取资料失败", error);
  }
}

async function pickAvatar() {
  if (uploadingAvatar.value) {
    return;
  }
  const token = getToken();
  if (!token) {
    uni.showToast({ title: "请先登录", icon: "none" });
    return;
  }

  uni.chooseImage({
    count: 1,
    sizeType: ["compressed"],
    sourceType: ["album", "camera"],
    success: async (result) => {
      const filePath = result.tempFilePaths?.[0];
      if (!filePath) {
        return;
      }
      uploadingAvatar.value = true;
      try {
        const uploadResult = await uploadAvatar(filePath, token);
        form.value.avatarUrl = uploadResult.url || "";
        uni.showToast({ title: "头像已上传", icon: "success" });
      } catch (error) {
        uni.showToast({ title: error.message || "头像上传失败", icon: "none" });
      } finally {
        uploadingAvatar.value = false;
      }
    },
  });
}

function uploadAvatar(filePath, token) {
  return new Promise((resolve, reject) => {
    uni.uploadFile({
      url: `${apiBase}/api/v1/upload/avatar`,
      filePath,
      name: "file",
      header: {
        Authorization: `Bearer ${token}`,
      },
      success: (response) => {
        try {
          const result = JSON.parse(response.data || "{}");
          if (response.statusCode !== 200 || result.code !== 0 || !result.data?.url) {
            reject(new Error(result.msg || "头像上传失败"));
            return;
          }
          resolve(result.data);
        } catch (error) {
          reject(new Error("解析上传结果失败"));
        }
      },
      fail: (error) => {
        reject(new Error(error.errMsg || "头像上传失败"));
      },
    });
  });
}

async function save() {
  if (saving.value) {
    return;
  }
  saving.value = true;
  try {
    const updatedUser = await requestJSON("/api/v1/user/profile", {
      method: "POST",
      data: buildProfilePayload(form.value),
    });
    setStoredUser(applyProfileToStoredUser(getStoredUser(), updatedUser));
    uni.showToast({
      title: "资料已保存",
      icon: "success",
      duration: 1200,
    });
    setTimeout(() => {
      uni.navigateBack();
    }, 300);
  } catch (error) {
    uni.showToast({ title: error.message || "保存失败", icon: "none" });
  } finally {
    saving.value = false;
  }
}

function goBack() {
  uni.navigateBack();
}
</script>

<style scoped>
.page { min-height: 100vh; background: #dfe9e6; }
.phone-shell { width: 100%; max-width: 430px; min-height: 100vh; margin: 0 auto; padding: calc(env(safe-area-inset-top) + 10px) 18px 32px; box-sizing: border-box; background: linear-gradient(180deg,#fbfdfb 0%,#eef6f3 48%,#f8faf7 100%); }
.topbar { height: 48px; display: flex; align-items: center; justify-content: space-between; }
.back { width: 38px; height: 38px; border-radius: 14px; background: #fffdfa; color: #0f5f75; display: flex; align-items: center; justify-content: center; font-size: 24px; box-shadow: 0 8px 20px rgba(16,32,51,.08); }
.top-title { color: #102033; font-size: 17px; font-weight: 900; }
.top-action { color: #0f5f75; font-size: 13px; font-weight: 900; }
.hero-card { margin-top: 14px; padding: 18px; border-radius: 24px; background: #102033; box-shadow: 0 16px 32px rgba(16,32,51,.16); }
.hero-title { display: block; color: #f8fbf9; font-size: 22px; font-weight: 900; }
.hero-desc { display: block; margin-top: 8px; color: rgba(248,251,249,.72); font-size: 13px; line-height: 1.6; }
.form-card { margin-top: 18px; padding: 20px 18px 6px; border-radius: 24px; background: rgba(255,255,252,.98); border: 1px solid rgba(224,232,229,.9); box-shadow: 0 10px 24px rgba(16,32,51,.07); }
.avatar-section { margin-bottom: 18px; }
.avatar-row { display: flex; align-items: center; gap: 14px; }
.avatar-preview,
.avatar-placeholder { width: 72px; height: 72px; border-radius: 999px; flex: 0 0 72px; }
.avatar-preview { background: #edf7f3; }
.avatar-placeholder { background: #edf2f7; color: #607386; display: flex; align-items: center; justify-content: center; font-size: 13px; font-weight: 800; }
.avatar-actions { min-width: 0; flex: 1; display: flex; flex-direction: column; gap: 8px; }
.upload-btn { width: 110px; height: 36px; margin: 0; border-radius: 18px; background: #102033; color: #f8fbf9; font-size: 13px; font-weight: 800; }
.upload-btn::after { border: 0; }
.upload-tip { color: #667a78; font-size: 12px; line-height: 1.6; }
.field-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 12px; }
.field-block { margin-bottom: 18px; }
.half { min-width: 0; }
.label { display: block; margin-bottom: 8px; color: #3d4f5d; font-size: 14px; font-weight: 700; }
.input-box { min-height: 48px; padding: 0 14px; border-radius: 14px; background: #f8fbf9; border: 1px solid rgba(205,218,214,.9); display: flex; align-items: center; }
.input { flex: 1; min-width: 0; color: #102033; font-size: 15px; }
.save-btn { height: 48px; margin-top: 24px; border-radius: 18px; background: #102033; color: #f8fbf9; font-size: 16px; font-weight: 800; box-shadow: 0 14px 28px rgba(16,32,51,.18); }
</style>
