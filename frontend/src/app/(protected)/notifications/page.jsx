"use client";

import { useEffect, useState } from "react";
import { useWebSocket } from "@/lib/UseWebsocket";
import Link from "next/link";
import styles from "./notifications.module.css";
import { timeAgo } from "@/_lib/timeago";

const notifLabel = (n) => {
  return `sent you a notification`;
};

export default function NotificationsPage() {
  const { setUnreadNotifCount, port, notifications } = useWebSocket();
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (port) {
      port.postMessage({ type: "get_notifications" });
      setLoading(false);
    }
    setUnreadNotifCount(0);
  }, [port]);

  return (
    <div className={styles.page}>
      <div className={styles.header}>
        <h1 className={styles.title}>Notifications</h1>
        {notifications.length > 0 && (
          <span className={styles.count}>{notifications.length}</span>
        )}
      </div>

      <div className={styles.list}>
        {loading ? (
          <div className={styles.skeleton}>
            {[1, 2, 3].map((i) => (
              <div key={i} className={styles.skeletonItem} />
            ))}
          </div>
        ) : notifications.length === 0 ? (
          <div className={styles.empty}>
            <span className={styles.emptyIcon}>🔔</span>
            <p>No notifications yet</p>
          </div>
        ) : (
          notifications.map((n) => (
            <Link
              href={`/profile/${n.ref_id}`}
              key={n.id}
              className={`${styles.item} ${!n.is_read ? styles.unread : ""}`}
            >
              <div className={styles.avatarWrapper}>
                <img
                  src={
                    n.from_image
                      ? `http://localhost:4001/pics/${n.from_image}`
                      : ""
                  }
                  alt={n.from_name}
                  className={styles.avatar}
                  onError={(e) => {
                    e.target.src = "";
                  }}
                />
                <span className={styles.notifIcon}>
                  {}
                </span>
              </div>

              <div className={styles.content}>
                <p className={styles.text}>
                  <strong className={styles.name}>{n.from_name}</strong>{" "}
                  {notifLabel(n)}
                </p>
                <span className={styles.time}>{timeAgo(n.created_at)}</span>
              </div>

              {!n.is_read && <span className={styles.dot} />}
            </Link>
          ))
        )}
      </div>
    </div>
  );
}
