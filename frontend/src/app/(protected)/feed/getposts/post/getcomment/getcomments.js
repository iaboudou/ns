"use client";

import { timeAgo } from "@/_lib/timeago";
import styles from "./getcomments.module.css";
import { Users } from "lucide-react";

export default function Comments({ COMMENTS }) {
  return (
    <div className={styles.commentscontainer}>
      {Array.isArray(COMMENTS.comments) &&
        COMMENTS.comments.map((comment) => (
          <div key={comment.id} className={styles.postcard}>
            <div className={styles.header}>
              {comment.profile_image
                ? <img
                    src={`/pics/${comment.profile_image}`}
                    className={styles.avatar}
                  />
                : <Users />}
              <div>
                <strong className={styles.NAME}>
                  {comment.firstname} {comment.lastname}
                </strong>
                <small className={styles.NAME}>
                  {timeAgo(comment.created_at)}
                </small>
              </div>
            </div>
            <p>{comment.content}</p>
            {comment.image_url && (
              <img src={comment.image_url} className={styles.commentimg} />
            )}
          </div>
        ))}
      {Number(COMMENTS.number_of_comments) > COMMENTS.comments.length && (
        <button className={styles.leadmore} onClick={COMMENTS.get10Comments}>
          load more ...
        </button>
      )}
    </div>
  );
}
