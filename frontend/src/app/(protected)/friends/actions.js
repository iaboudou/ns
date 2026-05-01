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

export async function GetFollowRequests(setRequests) {
  try {
    const res = await fetch(`${BASE}/api/get-follow-requests`, {
      credentials: "include",
    });

    if (handleUnauthorized(res)) return false;

    if (!res.ok) {
      return false;
    }

    const data = await res.json();
    setRequests(data.data || []);
    return true;

  } catch {
    return false;
  }
}

export async function ManageFollow(followerId, decision) {
  try {
    const res = await fetch(`${BASE}/api/manage-follow`, {
      method: "POST",
      credentials: "include",
      headers: {
        "Content-Type": "application/json"
      },
      body: JSON.stringify({
        follower_id: followerId,
        decision: decision
      })
    });

    if (handleUnauthorized(res)) return false;

    return res.ok;

  } catch {
    return false;
  }
}
