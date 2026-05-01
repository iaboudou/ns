const BASE = "http://localhost:4001"

// get personal info
export async function fetchPersonalInfo(uuid) {

    const res = await fetch(`${BASE}/api/getpersonalinfo?id=${uuid}`, {
        method: "GET",
        credentials: "include",
        cache: "no-store",
    })
    if (!res.ok) return {}
    let data = await res.json()

    return data.user || {}
}
