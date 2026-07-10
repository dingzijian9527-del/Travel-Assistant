const AI_PAGE_SESSION_KEY = "ai_page_session";

export function createAiPageState(storage = createStorageAdapter()) {
  return {
    getSession() {
      return readJSON(storage, AI_PAGE_SESSION_KEY, null);
    },
    saveSession(session) {
      if (!session || typeof session !== "object") {
        storage.removeItem(AI_PAGE_SESSION_KEY);
        return;
      }
      writeJSON(storage, AI_PAGE_SESSION_KEY, session);
    },
    clearSession() {
      storage.removeItem(AI_PAGE_SESSION_KEY);
    },
  };
}

function createStorageAdapter() {
  if (typeof uni !== "undefined" && uni && typeof uni.getStorageSync === "function") {
    return {
      getItem(key) {
        const value = uni.getStorageSync(key);
        return value === undefined || value === null || value === "" ? null : String(value);
      },
      setItem(key, value) {
        uni.setStorageSync(key, value);
      },
      removeItem(key) {
        uni.removeStorageSync(key);
      },
    };
  }

  const memory = new Map();
  return {
    getItem(key) {
      return memory.has(key) ? memory.get(key) : null;
    },
    setItem(key, value) {
      memory.set(key, value);
    },
    removeItem(key) {
      memory.delete(key);
    },
  };
}

function readJSON(storage, key, fallback) {
  const value = storage.getItem(key);
  if (!value) {
    return fallback;
  }
  try {
    return JSON.parse(value);
  } catch {
    return fallback;
  }
}

function writeJSON(storage, key, value) {
  storage.setItem(key, JSON.stringify(value));
}
