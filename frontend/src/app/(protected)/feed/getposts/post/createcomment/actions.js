

export async function createcomment(state, post_id) {

    let formdata = new FormData()
    formdata.append("content", state.text)
    formdata.append("post_id", post_id)
    if (state.picture) formdata.append("image_url", state.picture)

    try {
        let res = await fetch(`/api/createcomment`, {
            method: "POST",
            credentials: "include",
            body: formdata
        })
        if (!res.ok) {
            return null
        }

        const json = await res.json().catch(() => ({}));
        json.comment.created_at = Date.now().toLocaleString()
        return json.comment || null

    } catch {
        return null
    }

}