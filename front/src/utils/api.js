import { ref } from "vue";

export const apiBase = "http://127.0.0.1:8080";

const USER_STORAGE_KEY = "travel_user";
export const storedUserRef = ref(readStoredUser());

export function getToken() {
  return uni.getStorageSync("travel_token") || "";
}

export function getStoredUser() {
  return storedUserRef.value;
}

export function setStoredUser(user) {
  if (!user || typeof user !== "object") {
    return;
  }
  uni.setStorageSync(USER_STORAGE_KEY, user);
  storedUserRef.value = user;
}

export function clearStoredUser() {
  uni.removeStorageSync(USER_STORAGE_KEY);
  storedUserRef.value = null;
}

export function syncStoredUser() {
  storedUserRef.value = readStoredUser();
  return storedUserRef.value;
}

export async function requestJSON(path, options = {}) {
  const {
    method = "GET",
    data,
    headers = {},
    withAuth = true,
  } = options;
  const requestHeaders = {
    "Content-Type": "application/json",
    ...headers,
  };
  if (withAuth && getToken()) {
    requestHeaders.Authorization = `Bearer ${getToken()}`;
  }
  const response = await fetch(`${apiBase}${path}`, {
    method,
    headers: requestHeaders,
    body: data === undefined ? undefined : JSON.stringify(data),
  });
  const result = await response.json();
  if (!response.ok || !result || result.code !== 0) {
    throw new Error(result?.msg || "请求失败");
  }
  return result.data;
}

function readStoredUser() {
  const user = uni.getStorageSync(USER_STORAGE_KEY);
  return user && typeof user === "object" ? user : null;
}
