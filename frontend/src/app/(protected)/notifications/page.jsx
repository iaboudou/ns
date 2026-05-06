"use client";

import { useEffect, useState, useRef } from "react";
import { useWebSocket } from "@/lib/UseWebsocket";
import Link from "next/link";
import styles from "./notifications.module.css";
import { timeAgo } from "@/_lib/timeago";

const notifLabel = (n) => {
  if (n.type.startsWith("create_event:")) {
    return "created a new event in the group";
  }

  switch (n.type) {
    case "follow_request":
      return "sent you a follow request";
    case "follow_accepted":
      return "accepted your request";
    case "follow":
      return "started following you";
    case "like":
      return "liked your post";
    case "comment":
      return "commented on your post";
    case "message":
      return "sent you a message";
    case "group_request":
      return "sent you a group request";
    case "group_invite":
      return "sent you a group invite";
    default:
      return "sent you a notification";
  }
};

const notifLink = (n) => {
  if (n.type.startsWith("create_event:")) {
    const groupId = n.type.split(":")[1];
    return `/groups/${groupId}/events`;
  }

  switch (n.type) {
    case "follow":
      return `/profile/${n.ref_id}`;
    case "follow_request":
      return `/profile/${n.ref_id}`;
    case "follow_accepted":
      return `/profile/${n.ref_id}`;

    case "like":
    case "comment":
      return `/post/${n.ref_id}`;

    case "message":
      return `/chat/${n.ref_id}`;
    case "group_request":
      return `/groups/${n.group_id}/requests`;
    case "group_invite":
      console.log("from group_invite: ", n)
      return "/groups/invites";
    default:
      return "/";
  }
};

export default function NotificationsPage() {
  const { setUnreadNotifCount, port, notifications, portKey, sendFocus } = useWebSocket();
  const [loading, setLoading] = useState(true);
  const prevCount = useRef(notifications.length);
  useEffect(() => {
    sendFocus("notification");
    return () => { // if leave the notification tab
      if (port) {
        sendFocus("none");
      }
    };
  }, [port, portKey]);

  useEffect(() => {
    if (port) {
      port.postMessage({
        type: "send",
        payload: { type: "get_notifications" },
      });
      setLoading(false);
    }
    // setUnreadNotifCount(0);
  }, [port]);

  // scroll to top when a new notification arrives
  useEffect(() => {
    if (notifications.length > prevCount.current) {
      window.scrollTo({ top: 0, behavior: "smooth" });
    }
    prevCount.current = notifications.length;
  }, [notifications.length]);

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
          notifications.map((n) => {
            return (
              <Link
                href={notifLink(n)}
                key={n.id}
                className={`${styles.item} ${!n.is_read ? styles.unread : ""}`}
              >
                <div className={styles.avatarWrapper}>
                  <img
                    src={
                      n.from_image
                        ? `http://localhost:4001/pics/${n.from_image}`
                        : "/default.jpg"
                    }
                    alt={n.from_name}
                    className={styles.avatar}
                    onError={(e) => {
                      e.target.src = "/default.jpg";
                    }}
                  />
                </div>

                <div className={styles.content}>
                  <p className={styles.text}>
                    <span className={styles.name}>{n.from_name}</span>{" "}
                    {notifLabel(n)}
                  </p>
                  <span className={styles.time}>{timeAgo(n.created_at)}</span>
                </div>

                {!n.is_read && <div className={styles.unreadDot} />}
              </Link>
            );
          })
        )}
      </div>
    </div>
  );
}
