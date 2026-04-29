import { handleUnauthorized } from "@/_lib/redirect";

let BASE = process.env.BACKEND_URL;

export async function fetchFriendsUsers(search = "") {
    try {
        const word = encodeURIComponent(search || "");
        const res = await fetch(`${BASE}/api/getfriends?q=${word}`, {
            credentials: "include",
        });

        if (handleUnauthorized(res)) return [];

        if (!res.ok) {
            return [];
        }

        const data = await res.json().catch(() => ({}));
        return data?.users || [];
    } catch {
        return [];
    }
}

export const createpost = async (state) => {

    const formData = new FormData();
    formData.append("text", state.text.trim());
    formData.append("privacy", state.privacy);
    formData.append("group_id", state.group_id || "");

    if (state.picture) formData.append("Image", state.picture);

    let r = (state.selectedUsers?.map(e => e.id) || []).join(",")
    if (state.privacy === "private") {
        formData.append("allowed_users", r);
    }

    try {
        const res = await fetch(`${BASE}/api/createpost`, {
            method: "POST",
            body: formData,
            credentials: "include",
        });

        if (handleUnauthorized(res)) return null;

        const data = await res.json().catch(() => ({}));

        data.post.created_at = "now"

        return data.post || null;
    } catch (err) {
        console.error(err);
        return null;
    }
};


export function postIsValid(state) {
    if (state.text == "" && state.picture == null) {
        return false
    }
    return true
}
