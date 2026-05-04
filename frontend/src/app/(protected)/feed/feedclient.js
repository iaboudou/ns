'use client';

import { useEffect, useState, useRef } from 'react';
import CreatePost from '@/app/(protected)/feed/createpost/createpost';
import GetPosts from '@/app/(protected)/feed/getposts/getposts';
import { onPostCreated, onCommentCreated, loadPosts, setOpenComment, fetchNewPostsWhileScrooling } from './actions';
import { usePathname } from 'next/navigation';

// home page
// profile -> activity
// profile me -> posts
// profile user -> posts
// groups
export default function FeedClient({ INFO = {} }) {
  // fetch the posts based on the route
  const path = usePathname() || '';
  const isFetching = useRef(false);

  //
  const [state, setState] = useState({ posts: [], openComments: {}, loading: true, nbrofPosts: 0, should_stop_fetching: false });
  useEffect(() => {
    async function fetchposts() {
      try {
        let p = await loadPosts(state, path, INFO.section, INFO.profileId);
        p = p.posts;
        setState((prev) => ({ ...prev, posts: p }));
      } catch {
        setState((prev) => ({ ...prev, posts: [] }));
      } finally {
        setState((prev) => ({ ...prev, loading: false }));
      }
    }
    fetchposts();
  }, [INFO.section, path, INFO.profileId]);

  //
  let GETPOSTS = {
    state: state,
    setOpenComments: (postId) => setOpenComment(setState, postId),
    onCommentCreated: (post_id, comment) => {
      return onCommentCreated(setState, post_id, comment);
    },
    fetch: async () => {
      if (isFetching.current) return;
      isFetching.current = true;
      await fetchNewPostsWhileScrooling(setState, state, loadPosts, path, INFO.section, INFO.profileId);
      isFetching.current = false;
    },
    should_stop_fetching: state.should_stop_fetching,
    setState: setState,
  };

  let CREATEPOST = {
    onPostCreated: (newPost) => onPostCreated(setState, newPost),
    setState: setState,
  };

  return state.loading ? (
    <p>Loading...</p>
  ) : (
    <>
      {path == '/' || path.includes('/groups') ? <CreatePost CREATEPOST={CREATEPOST} /> : ''}
      {state.posts && state.posts?.length != 0 ? <GetPosts GETPOSTS={GETPOSTS} /> : 'no posts'}
    </>
  );
}
