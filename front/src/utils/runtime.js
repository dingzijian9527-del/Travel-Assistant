export function normalizeApiBase(value) {
  if (typeof value !== "string") {
    return "";
  }
  return value.trim().replace(/\/+$/, "");
}

export const apiBase = normalizeApiBase("http://127.0.0.1:8080");

export async function parseAPIResponse(response) {
  const contentType = String(response?.headers?.get?.("content-type") || "").toLowerCase();
  if (!contentType.includes("application/json")) {
    throw new Error("网关返回了非接口数据，请检查网关地址");
  }
  try {
    return await response.json();
  } catch {
    throw new Error("网关返回了无效接口数据，请稍后重试");
  }
}
