export function normalizeProfileForm(user) {
  return {
    nickname: textOrEmpty(user?.nickname),
    avatarUrl: textOrEmpty(user?.avatarUrl),
    homeCity: textOrEmpty(user?.homeCity),
    currentCity: textOrEmpty(user?.currentCity),
  };
}

export function buildProfilePayload(form) {
  const normalized = normalizeProfileForm(form);
  return {
    nickname: normalized.nickname,
    avatar_url: normalized.avatarUrl,
    home_city: normalized.homeCity,
    current_city: normalized.currentCity,
  };
}

export function applyProfileToStoredUser(user, profile) {
  const base = user && typeof user === "object" ? { ...user } : {};
  const normalized = normalizeProfileForm(profile);
  return {
    ...base,
    nickname: normalized.nickname,
    avatarUrl: normalized.avatarUrl,
    homeCity: normalized.homeCity,
    currentCity: normalized.currentCity,
  };
}

function textOrEmpty(value) {
  return typeof value === "string" ? value.trim() : "";
}
