export function normalizeApiBase(value) {
  if (typeof value !== "string") {
    return "";
  }
  return value.trim().replace(/\/+$/, "");
}

const buildEnv = typeof import.meta.env === "object" && import.meta.env ? import.meta.env : {};
export const apiBase = normalizeApiBase(buildEnv.VITE_API_BASE_URL);
