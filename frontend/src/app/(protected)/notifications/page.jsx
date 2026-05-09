"use client";

import { useEffect, useState, useRef } from "react";
import { useWebSocket } from "@/lib/UseWebsocket";
import styles from "./notifications.module.css";
import DisplayNotification from "@/components/notifications/DisplayNotification";

export default function NotificationsPage() {
  const { port, notifications, setNotifications } = useWebSocket();
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (port) {
      port.postMessage({
        type: "send",
        payload: { type: "get_notifications" },
      });
    }
  }, [port]);

  useEffect(() => setLoading(false), [notifications]);

  return (
    <div className={styles.page}>
      <div className={styles.header}>
        <h1 className={styles.title}>Notifications</h1>
      </div>

      <div className={styles.list}>
        {loading ? (
          <div className={styles.loading}>
            <div className={styles.spinner}></div>
            <p>Loading notifications...</p>
          </div>
        ) : notifications.length === 0 ? (
          <div className={styles.empty}>
            <div className={styles.emptyIcon}></div>
            <p>No notifications yet</p>
            <span>We'll let you know when something happens!</span>
          </div>
        ) : (
          notifications.map((notif) => (
            <DisplayNotification
              key={notif.id}
              notif={notif}
              port={port}
              setNotifications={setNotifications}
            />
          ))
        )}
      </div>
    </div>
  );
}
