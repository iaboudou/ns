"use client";

import Link from "next/link";
import styles from "@/app/(protected)/notifications/notifications.module.css";
import { timeAgo } from "@/_lib/timeago";

export default function DisplayNotification({ notif, port, setNotifications }) {
  const handleRead = () => {
    if (notif.is_read) return;

    port.postMessage({
      type: "send",
      payload: {
        type: "mark_notif_read",
        id: notif.id,
      },
    });

    setNotifications((prev) =>
      prev.map((n) => (n.id === notif.id ? { ...n, is_read: true } : n)),
    );
  };

  const generate = (n) => {
    switch (n.type) {
      case "event":
        return {
          text: " created an event",
          link: `/groups/${n.group_id}/events`,
        };
      case "group_invite":
        return { text: " sent you a group invite", link: `/groups/invites` };
      case "group_request":
        return {
          text: " requested to join your group",
          link: `/groups/${n.group_id}/requests`,
        };
      case "follow_request":
        return {
          text: " wants to follow you",
          link: `/profile/me?section=requests`,
        };
    }
  };

  return (
    <Link
      href={generate(notif).link}
      className={`${styles.item} ${!notif.is_read ? styles.unread : ""}`}
      onClick={handleRead}
    >
      <div className={styles.avatarWrapper}>
        <img
          src={
            notif.sender_profile
              ? `/pics/${notif.sender_profile}`
              : "/default.jpg"
          }
          alt={
            notif.sender_nickname
              ? notif.sender_nickname
              : notif.sender_fullname
          }
          className={styles.avatar}
          onError={(e) => {
            e.target.src = "/default.jpg";
          }}
        />
      </div>

      <div className={styles.content}>
        <p className={styles.text}>
          <span className={styles.name}>
            {notif.sender_nickname
              ? notif.sender_nickname
              : notif.sender_fullname}
          </span>
          {generate(notif).text}
        </p>
        <span className={styles.time}>{timeAgo(notif.created_at)}</span>
      </div>

      {!notif.is_read && <div className={styles.unreadDot} />}
    </Link>
  );
}
