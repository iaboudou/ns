'use client';
const BASE = "http://localhost:4001"
import Link from 'next/link';
import styles from './followsuggestions.module.css';
import { FollowUser, GetUsers } from './actions';
import { useState, useEffect } from 'react';
import { Users } from 'lucide-react';

export default function FollowSuggestions() {
  const [suggestions, setUsers] = useState([]);
  useEffect(() => {
    (async () => {
      await GetUsers(setUsers);
    })();
  }, []);

  // Handle follow click
  const handleFollow = async (userId) => {
    const message = await FollowUser(userId);
    if (!message) return;

    setUsers((prev) =>
      prev.map((user) => {
        if (user.id === userId) {
          if (message === 'follow have been successfully') {
            return { ...user, interactionStatus: 'following' };
          } else if (message === 'request have been sent') {
            return { ...user, interactionStatus: 'requested' };
          } else if (message === 'follow deleted' || message === 'follow request deleted') {
            return { ...user, interactionStatus: 'none' };
          }
        }
        return user;
      })
    );
  };

  if (!suggestions || suggestions.length == 0) return null;
  return (
    <div className={styles.followsuggestions}>
      <span className={styles.header}>People you may want to follow</span>

      <ul className={styles.list}>
        {suggestions.map((user) => {
          // get image
          const profileimage = user?.profile_image;
          const fullprofileimage = profileimage ? `${BASE}/pics/${profileimage}` : '';

          return (
            <li key={user.id} className={styles.item}>
              <Link href={`/profile/${user.id}`}>
                {
                  fullprofileimage ? <img src={fullprofileimage} className={styles.avatar} /> : <Users className={styles.placeholderIcon} />
                }
              </Link>

              <div className={styles.meta}>
                <span className={styles.name}>
                  {user.firstname} {user.lastname}
                </span>
              </div>

              <button className={styles.followBtn} type="button" onClick={() => handleFollow(user.id)}>
                {user.interactionStatus === 'following' ? 'unfollow' : user.interactionStatus === 'requested' ? 'requested' : user.account_privacy ? 'request' : 'follow'}
              </button>
            </li>
          );
        })}
      </ul>
    </div>
  );
}
