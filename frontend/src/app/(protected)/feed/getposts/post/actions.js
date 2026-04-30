import { handleUnauthorized } from "@/_lib/redirect";
let BASE = "http://localhost:4001";

export async function createLikeServer(post_id) {

    let res = await fetch(`${BASE}/api/createreaction`, {
        credentials: "include",
        method: "POST",
        body: JSON.stringify({ post_id: post_id })
    })

    if (handleUnauthorized(res)) return;

    let data = await res.json()
}

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

        if (handleUnauthorized(res)) return null;

        if (!res.ok) {
            console.error("error creating comment");
        }
    } catch (e) {
        console.error(e)
    }
}

export async function fetchComments(post_id, offset) {
    let res = await fetch(`${BASE}/api/getcomments`, {
        credentials: "include",
        method: "POST",
        body: JSON.stringify({ offset: offset, post_id: post_id }),
    })

    if (handleUnauthorized(res)) return;

    if (!res.ok) {
        return
    }
    let json = await res.json()
    return json.comments
}

export function set_comment_in_state(setState, post_id, data) {
    setState(prev => {
        return ({

            ...prev,
            posts: prev.posts.map(p =>
                p.id === post_id
                    ? { ...p, comments: [...(p.comments || []), ...data], offset: (p?.offset || 0) + 10 }
                    : p
            )
        })
    })
}