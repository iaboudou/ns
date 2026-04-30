'use client';

import { useState, useEffect } from 'react';
import Link from 'next/link';
import { Users } from 'lucide-react';
import styles from "./page.module.css";
import { GetFriends, FollowUser } from './actions';
import FollowSuggestions from "@/app/(protected)/feed/followsuggestions/followsuggestions";

const BASE = "http://localhost:4001";

function Friends() {
  const [friends, setFriends] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchData = async () => {
      setLoading(true);
      await GetFriends(setFriends);
      setLoading(false);
    };
    fetchData();
  }, []);

  const handleAction = async (userId) => {
    await FollowUser(userId);
    await GetFriends(setFriends);
  };

  return (
    <div className={styles.WRAPPER}>
      <div className={styles.MIDDLE}>
        <h1 className={styles.title}>Friends</h1>

        {loading ? (
          <div className={styles.loader}>Loading...</div>
        ) : (
          <div className={styles.content}>
            <div className={styles.section}>
              <h2 className={styles.sectionTitle}>Mutual Friends ({friends.length})</h2>
              <div className={styles.list}>
                {friends.length > 0 ? friends.map((user) => (
                  <div key={user.id} className={styles.row}>
                    <Link href={`/profile/${user.id}`} className={styles.userLink}>
                      {user.profile_image ? (
                        <img src={`${BASE}/pics/${user.profile_image}`} alt="" className={styles.avatar} />
                      ) : (
                        <div className={styles.avatarPlaceholder}><Users size={20} /></div>
                      )}
                      <span className={styles.userName}>{user.firstname} {user.lastname}</span>
                    </Link>
                    <button className={styles.actionBtn} onClick={() => handleAction(user.id)}>
                      Unfollow
                    </button>
                  </div>
                )) : <p className={styles.emptyMsg}>No friends found.</p>}
              </div>
            </div>
          </div>
        )}
      </div>
      
      <div className={styles.RIGHT}>
        <FollowSuggestions />
      </div>
    </div>
  );
}

export default Friends;
