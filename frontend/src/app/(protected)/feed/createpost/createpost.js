'use client';

import { useState } from 'react';
import styles from './createpost.module.css';
import { createpost, postIsValid } from './actions';
import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { SelectFreinds } from './select_freinds';
import { User, Image as ImageIcon } from 'lucide-react';

export default function CreatePost({ CREATEPOST }) {

  let path = usePathname();

  const [state, setState] = useState(() => {
    const a = JSON.parse(localStorage.getItem("user") || "null")
    return {
      privacy: 'public',
      text: '',
      picture: null,
      selectedUsers: [],
      porsonel_info: a,
      group_id: path.split("/")?.[2]
    }
  })

  // profile picture on home
  const imageURL = state.porsonel_info?.Avatar;
  const fullImageURL = imageURL ? `http://localhost:4001/pics/${imageURL}` : null;

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (postIsValid(state)) {
      let post = await createpost(state);

      CREATEPOST.onPostCreated(post);
      CREATEPOST.setState((prev) => ({ ...prev, nbrofPosts: prev.nbrofPosts + 1 }));
      setState((prev) => ({ ...prev, privacy: 'public', text: '', picture: null, selectedUsers: [] }));
    }
  };

  let STATE = {
    state,
    setState,
  }
  return (
    <div className={styles.createpost}>
      <form onSubmit={handleSubmit}>
        <div className={styles.topRow}>
          <div className={styles.userInfo}>
            {/* profile picture */}
            <Link href={`/profile/me`}>
              {imageURL ? (
                <img className={styles.profileImage} src={`${fullImageURL}`} alt="profile" />
              ) : (
                <div className={styles.profileImage} >
                   <User className={styles.IMAGEICON}/>
                </div>
              )}
            </Link>

            {/* user full name*/}
            <span className={styles.userName}>{state.porsonel_info?.Firstname + ' ' + state.porsonel_info?.Lastname}</span>
          </div>
          <button type="submit" className={styles.postBtn}>create Post</button>
        </div>

        {/* create post */}
        <textarea className={styles.textarea} placeholder="What's on your mind" value={state.text} onChange={(e) => setState({ ...state, text: e.target.value })} />

        {/* create post image */}
        <div className={styles.bottomRow}>
          <div className={styles.pictureinput}>
            <input type="file" accept="image/*" onChange={(e) => setState({ ...state, picture: e.target.files?.[0] || null })} className={styles.fileInput} id="fileInput" />

            {/* the name of the image chosen */}
            <label htmlFor="fileInput" className={styles.fileLabel} title={state.picture ? state.picture.name : ""}>
              <ImageIcon size={16} color="#c4c4c4" style={{ flexShrink: 0 }} />
              <span className={styles.filename}>
                {state.picture ? state.picture.name : 'Upload Image'}
              </span>
            </label>
          </div>

          {/* privacy */}
          {!path.includes('/groups') && (
            <>
              <select
                value={state.privacy}
                onChange={(e) => setState({ ...state, privacy: e.target.value })}
                className={styles.privacySelect}
              >
                <option value="public">Anyone can reply</option>
                <option value="followers">Followers</option>
                <option value="private">Private</option>
              </select>

              {/* select users who can see the post */}
              {state.privacy === "private" && (
                <SelectFreinds STATE={STATE} />
              )}
            </>
          )}
        </div>
      </form>
    </div>
  );
}
