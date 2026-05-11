import { BASE_URL } from "@/config.mjs";
const BASE = typeof window === "undefined" ? BASE_URL : "";

// get personal info
export async function fetchPersonalInfo(uuid) {
  const res = await fetch(`${BASE}/api/getpersonalinfo?id=${uuid}`, {
    method: "GET",
    credentials: "include",
    cache: "no-store",
  });
  if (!res.ok) return {};
  let data = await res.json().catch(() => ({}));

  return data.user || {};
}
