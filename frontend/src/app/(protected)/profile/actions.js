const BASE = "http://localhost:4001";

//
export async function fetchSwitchPrivacy() {
  let res = await fetch(`${BASE}/api/switchaccountprivacy`, {
    credentials: "include",
    method: "POST",
  });

  let data = await res?.text();
  if (!data) {
    return false;
  }
  if (res.headers.get("content-type").includes("application/json")) {
    data = JSON.parse(data);
    return true;
  }
  return false;
}

//
export async function fetchFollowers(id) {
  const res = await fetch(`${BASE}/api/follow?want=followers&id=${id}`, {
    credentials: "include",
  });
  if (!res.ok) return [];
  const data = await res.json();
  return data?.data || [];
}

//
export async function fetchFollowing(id) {
  const res = await fetch(`${BASE}/api/follow?want=following&id=${id}`, {
    credentials: "include",
  });
  if (!res.ok) return [];
  const data = await res.json();
  return data?.data || [];
}

export async function fetchRequests(id) {
  const res = await fetch(`${BASE}/api/follow?want=requests&id=${id}`, {
    credentials: "include",
  });
  if (!res.ok) return [];
  const data = await res.json();
  return data?.data || [];
}

export async function ManageFollow(followerId, decision) {
  try {
    const res = await fetch(`${BASE}/api/manage-follow`, {
      method: "POST",
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        follower_id: followerId,
        decision: decision,
      }),
    });

    return res.ok;
  } catch {
    return false;
  }
}
