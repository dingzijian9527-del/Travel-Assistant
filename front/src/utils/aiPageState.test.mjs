import assert from "node:assert/strict";
import { describe, it } from "node:test";
import { createAiPageState } from "./aiPageState.js";

function createMemoryStorage() {
  const store = new Map();
  return {
    getItem(key) {
      return store.has(key) ? store.get(key) : null;
    },
    setItem(key, value) {
      store.set(key, value);
    },
    removeItem(key) {
      store.delete(key);
    },
  };
}

describe("createAiPageState", () => {
  it("persists assistant session and restores it after page re-entry", () => {
    const storage = createMemoryStorage();
    const state = createAiPageState(storage);
    const session = {
      draftMessage: "帮我看看重庆两天怎么安排",
      promptGroupIndex: 2,
      promptExpanded: true,
      pendingTripPlan: { userText: "重庆两天", replyText: "第1天 解放碑" },
      messages: [
        { id: "msg-1", role: "ai", text: "你好", time: "10:30" },
        { id: "msg-2", role: "user", text: "重庆两天", time: "现在" },
      ],
    };

    state.saveSession(session);

    const restored = createAiPageState(storage).getSession();
    assert.deepEqual(restored, session);
  });

  it("clears stored assistant session", () => {
    const storage = createMemoryStorage();
    const state = createAiPageState(storage);

    state.saveSession({ messages: [{ id: "msg-1", role: "ai", text: "你好", time: "10:30" }] });
    state.clearSession();

    assert.equal(createAiPageState(storage).getSession(), null);
  });
});
