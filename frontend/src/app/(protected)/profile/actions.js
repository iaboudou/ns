

//
export async function fetchSwitchPrivacy() {
  let res = await fetch(`/api/switchaccountprivacy`, {
    credentials: "include",
    method: "PATCH",
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
  const res = await fetch(`/api/follow?want=followers&id=${id}`, {
    credentials: "include",
  });
  if (!res.ok) return [];
  const data = await res.json();
  return data?.data || [];
}

//
export async function fetchFollowing(id) {
  const res = await fetch(`/api/follow?want=following&id=${id}`, {
    credentials: "include",
  });
  if (!res.ok) return [];
  const data = await res.json();
  return data?.data || [];
}

export async function fetchRequests(id) {
  const res = await fetch(`/api/follow?want=requests&id=${id}`, {
    credentials: "include",
  });
  if (!res.ok) return [];
  const data = await res.json();
  return data?.data || [];
}

export async function ManageFollow(followerId, decision) {
  try {
    const res = await fetch(`/api/manage-follow`, {
      method: "PATCH",
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
