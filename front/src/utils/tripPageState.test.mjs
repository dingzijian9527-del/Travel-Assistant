import assert from "node:assert/strict";
import { describe, it } from "node:test";
import { createTripPageState } from "./tripPageState.js";

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

describe("createTripPageState", () => {
  it("persists selected trip and loaded list", () => {
    const storage = createMemoryStorage();
    const state = createTripPageState(storage);
    const trips = [{ id: "trip-1", title: "成都 / 3天" }, { id: "trip-2", title: "杭州 / 2天" }];

    state.saveTrips(trips);
    state.selectTrip("trip-2");

    const restored = createTripPageState(storage);
    assert.deepEqual(restored.getTrips(), trips);
    assert.equal(restored.getSelectedTripId(), "trip-2");
  });

  it("persists planner draft and can clear it", () => {
    const storage = createMemoryStorage();
    const state = createTripPageState(storage);
    const draft = { destinationName: "重庆", startDateText: "7月10日", endDateText: "7月12日", people: 2 };

    state.savePlannerDraft(draft);
    assert.deepEqual(createTripPageState(storage).getPlannerDraft(), draft);

    state.clearPlannerDraft();
    assert.equal(createTripPageState(storage).getPlannerDraft(), null);
  });
});
