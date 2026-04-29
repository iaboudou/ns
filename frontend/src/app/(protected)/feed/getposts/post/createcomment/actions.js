let BASE = process.env.BACKEND_URL;

export async function createcomment(state, post_id) {

    let formdata = new FormData()
    formdata.append("content", state.text)
    formdata.append("post_id", post_id)
    if (state.picture) formdata.append("image_url", state.picture)

    try {
        let res = await fetch(`${BASE}/api/createcomment`, {
            method: "POST",
            credentials: "include",
            body: formdata
        })
        if (!res.ok) {
            console.error("error creating comment");
        }

        const json = await res.json().catch(() => ({}));
        json.comment.created_at = Date.now().toLocaleString()
        return json.comment || null

    } catch {
        console.error("error creating comment");
        return null
    }

}