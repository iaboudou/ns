"use client";

import { SendVote } from "@/_lib/group";
import styles from "@/components/groups/styles/EventCard.module.css";
import { useState } from "react";

export default function EventCard({ event, groupId }) {
  const [vote, setVote] = useState(event.voted);
  const [going, setGoing] = useState(event.goingCount);
  const [notgoing, setNotGoing] = useState(event.notGoingCount);

  const dateObj = new Date(event.date);
  const month = dateObj.toLocaleString("en", { month: "short" });
  const day = dateObj.getDate();
  const fullDate = dateObj.toLocaleString("en", {
    weekday: "long",
    month: "long",
    day: "numeric",
  });
  const time = dateObj.toLocaleString("en", {
    hour: "2-digit",
    minute: "2-digit",
  });

  const handleVote = async (newVote) => {
    if (vote === newVote) return;

    if (newVote === "going") {
      setGoing((g) => g + 1);
      if (vote === "notgoing") {
        setNotGoing((n) => n - 1);
      }
    }

    if (newVote === "notgoing") {
      setNotGoing((n) => n + 1);
      if (vote === "going") {
        setGoing((g) => g - 1);
      }
    }

    SendVote(event.id, groupId, newVote)
      .then(() => setVote(newVote))
      .catch((err) => alert(err.message));
  };

  return (
    <div className={styles.wrapper}>
      <div className={styles.card}>
        <div className={styles.banner}>
          <div className={styles.bannerDate}>
            <span className={styles.month}>{month}</span>
            <span className={styles.day}>{day}</span>
          </div>
          <div className={styles.bannerInfo}>
            <div className={styles.title}>{event.title}</div>
            <div className={styles.time}>
              <svg viewBox="0 0 24 24" className={styles.clockIcon}>
                <circle cx="12" cy="12" r="9" />
                <polyline points="12 7 12 12 15 15" />
              </svg>
              {fullDate} · {time}
            </div>
          </div>
        </div>

        <div className={styles.body}>
          <p className={styles.desc}>{event.description}</p>

          <div className={styles.divider} />

          <div className={styles.sectionLabel}>Attendance</div>
          <div className={styles.attendeesRow}>
            <div className={`${styles.pill} ${styles.going}`}>
              <span className={styles.dot} />
              <span>{going}</span> going
            </div>
            <div className={`${styles.pill} ${styles.notgoing}`}>
              <span className={styles.dot} />
              <span>{notgoing}</span> not going
            </div>
          </div>

          <div className={styles.sectionLabel}>Your response</div>
          <div className={styles.voteRow}>
            <button
              className={`${styles.voteBtn} ${vote === "going" ? styles.activeGoing : ""}`}
              onClick={() => handleVote("going")}
            >
              <span className={styles.btnDot} /> Going
            </button>
            <button
              className={`${styles.voteBtn} ${vote === "notgoing" ? styles.activeNotgoing : ""}`}
              onClick={() => handleVote("notgoing")}
            >
              <span className={styles.btnDot} /> Not going
            </button>
          </div>

          <div className={styles.creatorRow}>
            <div className={styles.creatorText}>
              Created by
              {event.authorNickname
                ? <p>{event.authorNickname}</p>
                : <p>{event.authorName}</p>}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
