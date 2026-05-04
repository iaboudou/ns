"use client";

import Link from "next/link";
import { useParams } from "next/navigation";
import styles from "./UsersList.module.css";
import { useEffect, useState, useRef } from "react";
import { useWebSocket } from "@/lib/UseWebsocket";

export default function UsersList() {
  const [users, setUsers] = useState([]);
  const [hasMore, setHasmore] = useState(true);
  const [loading, setLoading] = useState(false);
  const [search, setSearch] = useState("");
  const { onlineUsers, port, messages } = useWebSocket();

  const handleFetchUsers = async (currentSearch, isReset = false) => {
    if (!hasMore && !isReset) return;

    const lastUser = isReset ? undefined : users.at(-1);
    const lastTime = lastUser?.created_at;
    const lastId = lastUser?.id;

    try {
      setLoading(true);

      const resp = await fetch(
        `http://localhost:4001/api/getUsers?last=${lastTime}&lastId=${lastId}&search=${currentSearch}`,
        {
          method: "GET",
          credentials: "include",
        },
      );

      const res = await resp.json();

      if (res.code !== 200) {
        alert(res.message);
        return;
      }

      const newUsers = res.data;

      if (newUsers.length < 10) setHasmore(false);

      setUsers(
        isReset
          ? newUsers
          : (prev) => {
              const map = new Map(prev.map((u) => [u.id, u]));
              newUsers.forEach((u) => map.set(u.id, u));
              return Array.from(map.values());
            },
      );
    } catch (err) {
      alert(err.message);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    const timer = setTimeout(() => {
      setHasmore(true);
      handleFetchUsers(search, true);
    }, 400);

    return () => clearTimeout(timer);
  }, [search]);

  const params = useParams();
  const lastProcessedMsgId = useRef(null);

  // Reorder users when a new message arrives
  useEffect(() => {
    if (messages.length === 0) return;
    const lastMsg = messages[messages.length - 1];
    
    // only move to top if it's a new live message, not history loading
    if (!lastMsg.isNew || lastMsg.id === lastProcessedMsgId.current) return;
    lastProcessedMsgId.current = lastMsg.id;

    setUsers((prev) => {
      const userIndex = prev.findIndex((u) => u.id === lastMsg.sender_id || u.id === lastMsg.receiver_id);
      if (userIndex === -1) {
         handleFetchUsers(search, true);
         return prev;
      }

      const newUsers = [...prev];
      const [oldUser] = newUsers.splice(userIndex, 1);
      
      const user = { ...oldUser };
      
      if (lastMsg.sender_id === user.id && params.id !== user.id) {
        user.unread_count = (user.unread_count || 0) + 1;
      }

      newUsers.unshift(user);
      return newUsers;
    });
  }, [messages, params.id]);

  return (
    <div className={styles.userContainer}>
      <h3>Users</h3>
      <div className={styles.userList}>
        <input
          className={styles.input}
          type="text"
          placeholder="looking for someone ?"
          value={search}
          onChange={(e) => setSearch(e.target.value)}
        />
        {users.map((u) => (
          <Link
            key={u.id}
            className={styles.userItem}
            href={`/chat/${u.id}`}
            onClick={() => {
              // Reset unread count when clicking
              setUsers(prev => prev.map(usr => usr.id === u.id ? { ...usr, unread_count: 0 } : usr));
              port.postMessage({
                type: "send",
                payload: {
                  type: "load_history",
                  receiver_Id: u.id,
                  last_read_time: 0,
                },
              });
            }}
          >
            <div className={styles.userInfo}>
              <span className={styles.fullname}>{u.firstname} {u.lastname}</span>
              {u.nickname && <span className={styles.nickname}>@{u.nickname}</span>}
            </div>
            <div className={styles.userStatus}>
              {u.unread_count > 0 && (
                <span className={styles.unreadBadge}>{u.unread_count}</span>
              )}
              {onlineUsers.includes(u.id) && (
                <span className={styles.onlineBtn} />
              )}
            </div>
          </Link>
        ))}

        {hasMore && (
          <button onClick={() => handleFetchUsers(search)} disabled={loading}>
            {loading ? "Loading..." : "Load More"}
          </button>
        )}
      </div>
    </div>
  );
}
