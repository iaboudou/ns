const BASE = "http://localhost:4001";
import { handleUnauthorized } from "@/_lib/redirect";

export async function GetFriends(setFriends) {
  try {
    const res = await fetch(`${BASE}/api/getfriends`, {
      credentials: "include"
    });

    if (handleUnauthorized(res)) return false;

    if (!res.ok) {
      return false;
    }

    const data = await res.json();
    const friends = data.users || data.data || [];

    setFriends(Array.isArray(friends) ? friends : []);
    return true;
    return true;

  } catch {
    return false;
  }
}

export async function GetSuggestions(setSuggestions) {
  try {
    const res = await fetch(`${BASE}/api/getsuggestionfollowers`, {
      credentials: "include"
    });

    if (handleUnauthorized(res)) return false;

    if (!res.ok) {
      return false;
    }

    const data = await res.json();

    if (!data || !data.data || data.data === "no suggestions") {
      setSuggestions([]);
      return true;
    }

    setSuggestions(data.data);
    return true;

  } catch {
    return false;
  }
}

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
