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
  const [requests, setRequests] = useState([]);
  const [loadingRequests, setLoadingRequests] = useState(true);

  const fetchData = async () => {
    setLoading(true);
    await GetFriends(setFriends);
    setLoading(false);
  };

  const fetchRequests = async () => {
    setLoadingRequests(true);
    try {
      const res = await fetch(`${BASE}/api/get-follow-requests`, {
        credentials: "include",
      });
      if (res.ok) {
        const data = await res.json();
        setRequests(data.data || []);
      }
    } catch (e) {
      console.error(e);
    }
    setLoadingRequests(false);
  };

  useEffect(() => {
    fetchData();
    fetchRequests();
  }, []);

  const handleAction = async (userId) => {
    await FollowUser(userId);
    await GetFriends(setFriends);
  };

  const handleManageRequest = async (followerId, decision) => {
    const res = await fetch(`${BASE}/api/manage-follow`, {
      method: "POST",
      body: JSON.stringify({ follower_id: followerId, decision }),
      credentials: "include",
    });
    if (res.ok) {
      fetchRequests();
      await GetFriends(setFriends);
    }
  };

  return (
    <div className={styles.WRAPPER}>
      <div className={styles.MIDDLE}>
        <h1 className={styles.title}>Friends</h1>

        {/* Follow Requests Section */}
        {!loadingRequests && requests.length > 0 && (
          <div className={styles.section}>
            <h2 className={styles.sectionTitle}>
              Follow Requests ({requests.length})
            </h2>
            <div className={styles.list}>
              {requests.map((user) => (
                <div key={user.id} className={styles.row}>
                  <Link href={`/profile/${user.id}`} className={styles.userLink}>
                    {user.profile_image ? (
                      <img
                        src={`${BASE}/pics/${user.profile_image}`}
                        alt=""
                        className={styles.avatar}
                      />
                    ) : (
                      <div className={styles.avatarPlaceholder}>
                        <Users size={20} />
                      </div>
                    )}
                    <div>
                      <span className={styles.userName}>
                        {user.firstname} {user.lastname}
                      </span>
                      <span className={styles.subText}>Wants to follow you</span>
                    </div>
                  </Link>
                  <div className={styles.requestActions}>
                    <button
                      className={`${styles.actionBtn} ${styles.acceptBtn}`}
                      onClick={() => handleManageRequest(user.id, 'accepted')}
                    >
                      Accept
                    </button>
                    <button
                      className={`${styles.actionBtn} ${styles.rejectBtn}`}
                      onClick={() => handleManageRequest(user.id, 'rejected')}
                    >
                      Reject
                    </button>
                  </div>
                </div>
              ))}
            </div>
          </div>
        )}

        {/* Friends List Section */}
        {loading ? (
          <div className={styles.loader}>Loading...</div>
        ) : (
          <div className={styles.content}>
            <div className={styles.section}>
              <h2 className={styles.sectionTitle}>
                Mutual Friends ({friends.length})
              </h2>
              <div className={styles.list}>
                {friends.length > 0 ? (
                  friends.map((user) => (
                    <div key={user.id} className={styles.row}>
                      <Link href={`/profile/${user.id}`} className={styles.userLink}>
                        {user.profile_image ? (
                          <img
                            src={`${BASE}/pics/${user.profile_image}`}
                            alt=""
                            className={styles.avatar}
                          />
                        ) : (
                          <div className={styles.avatarPlaceholder}>
                            <Users size={20} />
                          </div>
                        )}
                        <span className={styles.userName}>
                          {user.firstname} {user.lastname}
                        </span>
                      </Link>
                      <button
                        className={styles.actionBtn}
                        onClick={() => handleAction(user.id)}
                      >
                        Unfollow
                      </button>
                    </div>
                  ))
                ) : (
                  <p className={styles.emptyMsg}>No friends found.</p>
                )}
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