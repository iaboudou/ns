"use client";

let BASE = process.env.BACKEND_URL;
import { timeAgo } from "@/_lib/timeago";
import styles from "./getcomments.module.css";

export default function Comments({ COMMENTS }) {
  return (
    <div className={styles.commentscontainer}>
      {Array.isArray(COMMENTS.comments) &&
        COMMENTS.comments.map((comment) => (
          <div key={comment.id} className={styles.postcard}>
            <div className={styles.header}>
              {
                comment.profile_image ?
                  <img
                    src={`${BASE}/pics/${comment.profile_image}`}
                    className={styles.avatar}
                  />
                  : <img
                    src="/avatar.jpg"
                    className={styles.avatar}
                  />
              }
              <div>
                <strong className={styles.NAME}>{comment.firstname} {comment.lastname}</strong>
                <small>{timeAgo(comment.created_at)}</small>
              </div>
            </div>
            <p>{comment.content}</p>
            {comment.image_url && (
              <img src={`${BASE}${comment.image_url}`} className={styles.commentimg} />
            )}
          </div>
        ))}
      {Number(COMMENTS.number_of_comments) > COMMENTS.comments.length && (
        <button className={styles.leadmore} onClick={COMMENTS.get10Comments}>load more ...</button>
      )}
    </div>
  );
}
