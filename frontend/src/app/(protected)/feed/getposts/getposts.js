"use client";

import { useEffect, useRef } from "react";
import Post from "@/app/(protected)/feed/getposts/post/post";
import styles from "./getposts.module.css";

export default function GetPosts({ GETPOSTS }) {
  let ref = useRef(null);

  useEffect(() => {
    let observer = new IntersectionObserver((entr) => {
      if (entr[0].isIntersecting) {
        if (GETPOSTS.should_stop_fetching) {
          observer.disconnect();
          return;
        }
        GETPOSTS.fetch();
      }
    });
    observer.observe(ref.current);

    return () => observer.disconnect();
  }, [GETPOSTS.fetch, GETPOSTS.should_stop_fetching]);

  return (
    <div className={styles.wrapper}>
      {(GETPOSTS.state.posts || []).map((post) => {
        let POST = {
          post: post,
          openComments: GETPOSTS.state.openComments,
          setOpenComments: () => GETPOSTS.setOpenComments(post.id),
          onCommentCreated: GETPOSTS.onCommentCreated,
          setState: GETPOSTS.setState,
          state: GETPOSTS.state,
        };

        return <Post key={post.id} POST={POST} />;
      })}
      <span ref={ref} className={styles.leading} />
    </div>
  );
}
