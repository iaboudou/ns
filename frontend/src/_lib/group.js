const BASE =
  typeof window === "undefined"
    ? "http://backend:4001/api/groups"
    : "/api/groups";

export async function GetGroups(tab, search = "", lastId = "", lastTime = "") {
  const resp = await fetch(
    `${BASE}?want=${tab}&search=${search}&last=${lastTime}&lastId=${lastId}`,
    {
      method: "GET",
      credentials: "include",
    },
  );

  const result = await resp.json();

  if (result.code === 200) return result.data;
  else throw new Error(result.message); //error possible: 500
}

export async function GetGroup(id, cookie) {
  const resp = await fetch(`${BASE}/${id}`, {
    method: "GET",
    headers: { cookie },
  });

  const res = await resp.json();

  if (res.code === 200) return res.data;
  else throw new Error(res.message); //error possible: 500/404
}

export async function GetData(tab, groupId, lastTime = "", lastId = "") {
  const resp = await fetch(
    `${BASE}/${groupId}/${tab}?last=${lastTime}&lastId=${lastId}`,
    {
      method: "GET",
      credentials: "include",
    },
  );

  const res = await resp.json();

  if (res.code === 200) return res.data;
  else throw new Error(res.message); //error possible: 500/404/403
}

export async function getUsers(
  lastTime = "",
  lastId = "",
  search = "",
  groupId,
) {
  const resp = await fetch(
    `${BASE}/${groupId}/users?want=users&last=${lastTime}&lastId=${lastId}&search=${search}`,
    {
      method: "GET",
      credentials: "include",
    },
  );
  const res = await resp.json();

  if (res.code === 200) return res.data;
  else throw new Error(res.message); //error possible: 500/404
}

export async function CreateGroup(formData) {
  const resp = await fetch(`${BASE}/`, {
    method: "POST",
    body: formData,
    credentials: "include",
    next: { revalidate: 0 },
  });

  const result = await resp.json();

  if (result.code === 201) return result.data;
  else throw new Error(result.message); //error possible: 500/400/409
}

export async function CreateEvent(title, description, date, groupId, vote) {
  const resp = await fetch(`${BASE}/${groupId}/events`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      title,
      description,
      date,
      vote,
    }),
    credentials: "include",
  });

  const result = await resp.json();

  if (result.code === 201) return result.data;
  else throw new Error(result.message); //error possible: 500/400/409
}

export async function SendGroupRequest(groupId) {
  const resp = await fetch(`${BASE}/${groupId}/requests`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    credentials: "include",
  });

  const result = await resp.json();

  if (result.code == 200) return;
  else throw new Error(result.message); // 404/500
}

export async function sendGroupInvite(groupId, userId) {
  const resp = await fetch(`${BASE}/${groupId}/invites`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ invitedUser: userId }),
    credentials: "include",
  });

  const result = await resp.json();

  if (result.code == 200) return;
  else throw new Error(result.message); // 400/404/500
}

export async function SendVote(eventId, groupId, vote) {
  const resp = await fetch(`${BASE}/${groupId}/events/${eventId}`, {
    headers: { "Content-Type": "application/json" },
    method: "PATCH",
    body: JSON.stringify({
      vote,
    }),
    credentials: "include",
  });

  const result = await resp.json();

  if (result.code === 200) return;
  else throw new Error(result.message); //error possible: 500/404/400
}
export async function SendDecision(
  groupId,
  decision,
  requesterId = "",
  invitedBy = "",
) {
  const resp = await fetch(`${BASE}/${groupId}/requests`, {
    method: "PATCH",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      sender: requesterId,
      decision,
      invited_by: invitedBy,
    }),
    credentials: "include",
  });

  const result = await resp.json();

  if (result.code === 200) return;
  else throw new Error(result.message); //error possible: 500/404/400
}

export async function DeleteGroup(groupId) {
  const resp = await fetch(`${BASE}/${groupId}`, {
    method: "DELETE",
    headers: {
      "Content-Type": "application/json",
    },
    credentials: "include",
  });

  const result = await resp.json();

  if (result.code === 200) return;
  else throw new Error(result.message); //error possible: 500/403/404
}

export async function LeaveGroup(groupId) {
  const resp = await fetch(`${BASE}/${groupId}/me`, {
    method: "DELETE",
    headers: {
      "Content-Type": "application/json",
    },
    credentials: "include",
  });

  const result = await resp.json();

  if (result.code === 200) return;
  else throw new Error(result.message); //error possible: 500/404
}
