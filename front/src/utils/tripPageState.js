const TRIPS_KEY = "trip_page_trips";
const SELECTED_TRIP_ID_KEY = "trip_page_selected_trip_id";
const SELECTED_TRIP_KEY = "trip_page_selected_trip";
const PLANNER_DRAFT_KEY = "trip_page_planner_draft";

export function createTripPageState(storage = createStorageAdapter()) {
  return {
    getTrips() {
      return readJSON(storage, TRIPS_KEY, []);
    },
    saveTrips(trips) {
      writeJSON(storage, TRIPS_KEY, Array.isArray(trips) ? trips : []);
    },
    getSelectedTripId() {
      return readText(storage, SELECTED_TRIP_ID_KEY, "");
    },
    selectTrip(tripId) {
      if (!tripId) {
        storage.removeItem(SELECTED_TRIP_ID_KEY);
        return;
      }
      storage.setItem(SELECTED_TRIP_ID_KEY, String(tripId));
    },
    getSelectedTrip() {
      return readJSON(storage, SELECTED_TRIP_KEY, null);
    },
    saveSelectedTrip(trip) {
      if (!trip || typeof trip !== "object") {
        storage.removeItem(SELECTED_TRIP_KEY);
        return;
      }
      writeJSON(storage, SELECTED_TRIP_KEY, trip);
    },
    clearTripState() {
      storage.removeItem(TRIPS_KEY);
      storage.removeItem(SELECTED_TRIP_ID_KEY);
      storage.removeItem(SELECTED_TRIP_KEY);
    },
    getPlannerDraft() {
      return readJSON(storage, PLANNER_DRAFT_KEY, null);
    },
    savePlannerDraft(draft) {
      if (!draft || typeof draft !== "object") {
        storage.removeItem(PLANNER_DRAFT_KEY);
        return;
      }
      writeJSON(storage, PLANNER_DRAFT_KEY, draft);
    },
    clearPlannerDraft() {
      storage.removeItem(PLANNER_DRAFT_KEY);
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

function readText(storage, key, fallback) {
  const value = storage.getItem(key);
  return value === null || value === undefined || value === "" ? fallback : String(value);
}

function readJSON(storage, key, fallback) {
  const value = storage.getItem(key);
  if (!value) {
    return fallback;
  }
  try {
    return JSON.parse(value);
  } catch (error) {
    return fallback;
  }
}

function writeJSON(storage, key, value) {
  storage.setItem(key, JSON.stringify(value));
}
