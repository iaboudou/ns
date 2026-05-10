"use client";

import { useState } from "react";
import styles from "./createpost.module.css";
import { createpost, postIsValid } from "./actions";
import Link from "next/link";
import { usePathname } from "next/navigation";
import { SelectFreinds } from "./select_freinds";
import { User, Image as ImageIcon } from "lucide-react";
import { useWebSocket } from "@/lib/UseWebsocket";

export default function CreatePost({ CREATEPOST }) {
  const { myInfo } = useWebSocket();
  let path = usePathname();

  const [state, setState] = useState(() => {
    return {
      privacy: "public",
      text: "",
      picture: null,
      selectedUsers: [],
      group_id: path.split("/")?.[2],
    };
  });

  // profile picture on home
  const imageURL = myInfo?.profile_image;
  const fullImageURL = imageURL ? `/pics/${imageURL}` : null;

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (postIsValid(state)) {
      let post = await createpost(state);
      if (post) {
        // take the user info
        post.firstname = myInfo?.firstname;
        post.lastname = myInfo?.lastname;
        post.profile_image = myInfo?.profile_image;

        CREATEPOST.onPostCreated(post);
        CREATEPOST.setState((prev) => ({
          ...prev,
          nbrofPosts: prev.nbrofPosts + 1,
        }));
      }
      setState((prev) => ({
        ...prev,
        privacy: "public",
        text: "",
        picture: null,
        selectedUsers: [],
      }));
    }
  };

  let STATE = {
    state,
    setState,
  };
  return (
    <div className={styles.createpost}>
      <form onSubmit={handleSubmit}>
        <div className={styles.topRow}>
          <div className={styles.userInfo}>
            {/* profile picture */}
            <Link href={`/profile/me`}>
              {imageURL
                ? <img
                    className={styles.profileImage}
                    src={`${fullImageURL}`}
                    alt="profile"
                  />
                : <div className={styles.profileImage}>
                    <User className={styles.IMAGEICON} />
                  </div>}
            </Link>

            {/* user full name*/}
            <span className={styles.userName}>
              {myInfo?.firstname + " " + myInfo?.lastname}
            </span>
          </div>
          <button type="submit" className={styles.postBtn}>
            create Post
          </button>
        </div>

        {/* create post */}
        <textarea
          className={styles.textarea}
          placeholder="Describe your quest to the guild..."
          value={state.text}
          onChange={(e) => setState({ ...state, text: e.target.value })}
          maxLength={600}
        />

        {/* create post image */}
        <div className={styles.bottomRow}>
          <div className={styles.pictureinput}>
            <input
              type="file"
              accept="image/*"
              onChange={(e) =>
                setState({ ...state, picture: e.target.files?.[0] || null })
              }
              className={styles.fileInput}
              id="fileInput"
            />

            {/* the name of the image chosen */}
            <label
              htmlFor="fileInput"
              className={styles.fileLabel}
              title={state.picture ? state.picture.name : ""}
            >
              <ImageIcon
                size={16}
                color="#c4c4c4"
                className={styles.imageIcon}
              />
              <span className={styles.filename}>
                {state.picture ? state.picture.name : "Upload Image"}
              </span>
            </label>
          </div>

          {/* privacy */}
          {!path.includes("/groups") && (
            <>
              <select
                value={state.privacy}
                onChange={(e) =>
                  setState({ ...state, privacy: e.target.value })
                }
                className={styles.privacySelect}
              >
                <option value="public">Anyone can reply</option>
                <option value="followers">Followers</option>
                <option value="private">Private</option>
              </select>

              {/* select users who can see the post */}
              {state.privacy === "private" && <SelectFreinds STATE={STATE} />}
            </>
          )}
        </div>
      </form>
    </div>
  );
}
