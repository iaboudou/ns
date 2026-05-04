
let BASE = process.env.BACKEND_URL;
// insert post to the state
export function onPostCreated(setState, newPost) {
  setState((prev) => ({
    ...prev,
    posts: [newPost, ...(Array.isArray(prev?.posts) ? prev?.posts : [])],
  }));
}

// insert comment to the state
export function onCommentCreated(setState, post_id, comment) {
  setState((prev) => ({
    ...prev,
    posts: prev.posts.map((p) =>
    (p.id === post_id
      ? {
        ...p,
        comments: [comment, ...(p.comments || [])],
        number_of_comments: p.number_of_comments + 1,
        offset: (p?.offset || 0) + 1,
      } : p)),
  }));
}

//
export async function loadPosts(state, pathname, section, profileId) {

  let user_id = profileId || pathname.split('/')?.[2] || '';
  let group_id = pathname.split('/')?.[2] || '';
  let page = pathname.split('/')?.[1] || '';

  // get the current page
  switch (true) {
    case page == "":
      page = "home";
      break;
    case page == "profile" && pathname.split('/')?.[2] == "me" && section == "activity":
      page = "profile-me-activity"
      break;
    case page == "profile" && pathname.split('/')?.[2] == "me" && section == "posts":
      page = "profile-me-posts"
      break;
    case page == "profile" && section == "posts":
      page = "profille-other-posts"
      break;
    case page == "groups":
      page = "goups"
      break
  }

  const params = new URLSearchParams({
    page,
    offset: state.nbrofPosts,
    user_id,
    section: section || "",
    group_id,
  })

  const res = await fetch(`http://localhost:4001/api/getposts?${params.toString()}`, {
    method: "GET",
    credentials: "include",
  })
  const json = await res.json().catch(() => { });

  if (!res.ok) {
    return { posts: [] };
  }
  return json || { posts: [] };
}

//
export function setOpenComment(setState, postId) {
  return setState((prev) => ({
    ...prev,
    openComments: { ...prev.openComments, [postId]: !prev.openComments[postId] },
  }));
}

//
export async function fetchNewPostsWhileScrooling(setState, state, loadPosts, path, section, profileId) {
  //
  await new Promise(resolve => setTimeout(resolve, 500));

  let newState = { ...state, nbrofPosts: state.nbrofPosts + 10 };
  const res = await loadPosts(newState, path, section, profileId);
  if (!res || !Array.isArray(res.posts) || res.posts.length === 0) {
    setState((prev) => ({ ...prev, should_stop_fetching: true }));
    return;
  }
  newState.posts.push(...res.posts);

  setState((prev) => ({
    ...prev,
    posts: newState.posts,
    nbrofPosts: newState.nbrofPosts,
  }));
}
