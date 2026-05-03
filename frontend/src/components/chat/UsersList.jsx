"use client";

import Link from "next/link";
import styles from "./UsersList.module.css";
import { useEffect, useState } from "react";
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
            onClick={() =>
              port.postMessage({
                type: "send",
                payload: {
                  type: "load_history",
                  receiver_Id: u.id,
                  last_read_time:
                    messages.length === 0 ? 0 : messages.at(-1).created_at,
                },
              })
            }
          >
            <div key={u.id}>
              {u.nickname ? u.nickname : u.fisrtname + " " + u.lastname}
              {onlineUsers.includes(u.id) && " (online)"}
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
