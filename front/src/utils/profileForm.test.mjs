import assert from "node:assert/strict";
import { describe, it } from "node:test";
import { applyProfileToStoredUser, buildProfilePayload, normalizeProfileForm } from "./profileForm.js";

describe("profileForm", () => {
  it("normalizes empty profile form values", () => {
    assert.deepEqual(normalizeProfileForm(null), {
      nickname: "",
      avatarUrl: "",
      homeCity: "",
      currentCity: "",
    });
  });

  it("builds trimmed profile payload", () => {
    assert.deepEqual(buildProfilePayload({
      nickname: "  小王  ",
      avatarUrl: " https://example.com/a.png ",
      homeCity: " 成都 ",
      currentCity: " 上海 ",
    }), {
      nickname: "小王",
      avatar_url: "https://example.com/a.png",
      home_city: "成都",
      current_city: "上海",
    });
  });

  it("applies saved profile fields to stored user", () => {
    assert.deepEqual(applyProfileToStoredUser({
      id: "1",
      phone: "13800138000",
      nickname: "旧昵称",
    }, {
      nickname: "新昵称",
      avatarUrl: "https://example.com/new.png",
      homeCity: "重庆",
      currentCity: "杭州",
    }), {
      id: "1",
      phone: "13800138000",
      nickname: "新昵称",
      avatarUrl: "https://example.com/new.png",
      homeCity: "重庆",
      currentCity: "杭州",
    });
  });
});
