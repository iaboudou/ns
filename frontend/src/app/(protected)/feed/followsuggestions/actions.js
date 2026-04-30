let BASE = "http://localhost:4001"
import { handleUnauthorized } from "@/_lib/redirect";

/**
 * Get suggested users
 */
export async function GetUsers(setUsers) {
  try {
    const res = await fetch(`${BASE}/api/getsuggestionfollowers`, {
      credentials: "include"
    });

    if (handleUnauthorized(res)) return false;

    if (!res.ok) {
      return false;
    }

    const data = await res.json();

    if (!data || !data.data) {
      return false;
    }

    if (data.data === "no suggestions") {
      return true;
    }

    setUsers(data.data);
    return true;

  } catch {
    return false;
  }
}

/**
 * Follow a user
 */
export async function FollowUser(userId) {
  try {
    const res = await fetch(`${BASE}/api/follow`, {
      method: "POST",
      credentials: "include",
      headers: {
        "Content-Type": "application/json"
      },
      body: JSON.stringify({
        followed_id: userId
      })
    });

    if (handleUnauthorized(res)) return null;

    if (!res.ok) {
      return null;
    }

    const data = await res.json();
    return data?.message || null;

  } catch {
    return null;
  }
}