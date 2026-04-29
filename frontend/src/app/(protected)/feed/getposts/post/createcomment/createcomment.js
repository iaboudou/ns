"use client";

import { useState } from "react";
import styles from "./createcomment.module.css";
import { createcomment } from "./actions";
import { Image as ImageIcon } from "lucide-react";

export default function CreateComment({ post, onCommentCreated }) {

  let postID = post.id
  // post.offset += 1
  // initialize state for 
  const [state, setState] = useState({
    text: "",
    picture: null,
  })
  const handleSubmit = async (e) => {
    e.preventDefault()
    console.log("the state :", state)
    let comment = await createcomment(state, postID)
    setState({ text: "", picture: null })
    onCommentCreated(postID, comment)
  }

  return (
    <div >
      <form onSubmit={handleSubmit} className={styles.commentform}>
        <textarea
          placeholder="Write a comment..."
          value={state.text}
          onChange={(e) => setState({ ...state, text: e.target.value })}
          className={styles.textarea}
        />

        <div className={styles.pictureandsubmitcontainer}>
          <input
            id={`imageInput-${postID}`}
            type="file"
            accept="image/*"
            onChange={(e) => setState({ ...state, picture: e.target.files?.[0] || null })}
          />

          <label htmlFor={`imageInput-${postID}`} title={state.picture ? state.picture.name : ""}>
            <ImageIcon className={styles.IMAGEICON} />
            <span className={styles.filename}>
              {state.picture ? state.picture.name : "Image"}
            </span>
          </label>

          <button type="submit">comment</button>
        </div>
      </form>
    </div>
  );
}