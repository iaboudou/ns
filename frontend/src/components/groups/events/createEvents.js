"use client";

import { CreateEvent } from "@/_lib/group";
import styles from "@/components/groups/styles/groups.module.css";
import { useState } from "react";

export default function CreateEvents({ groupId, setEvents }) {
  const [eventTitle, setEventTitle] = useState("");
  const [eventDescription, setEventDescription] = useState("");
  const [eventDateTime, setEventDateTime] = useState("");
  const [error, setError] = useState("");
  const [myVote, setMyVote] = useState("going"); // par défaut going

  const handleCreateEvent = async (e) => {
    e.preventDefault();

    if (
      eventTitle.trim() === "" ||
      eventDescription.trim() === "" ||
      eventDateTime.trim() === ""
    ) {
      setError("please fill all the fields");
      return;
    }

    if (myVote !== "going" && myVote !== "notgoing") {
      setError("you can only chose to go or not to go");
      return;
    }

    const date = new Date(eventDateTime);

    if (date.getTime() <= Date.now() + 5 * 60 * 1000) {
      setError("an event must be at least in 5 minutes");
      return;
    }

    const isoDate = date.toISOString();

    CreateEvent(
      eventTitle.trim(),
      eventDescription.trim(),
      isoDate,
      groupId,
      myVote,
    )
      .then((newEvents) => {
        setEventTitle("");
        setEventDescription("");
        setEventDateTime("");
        setError("");
        setEvents((prev) => [newEvents, ...prev]);
      })
      .catch((err) => {
        if (err.message === "Something wrong happened. Please try later")
          alert(err.message);
        else setError(err.message);
      });
  };

  return (
    <form className={styles.eventForm} onSubmit={handleCreateEvent}>
      <div className={styles.formGroup}>
        {error && <p className={styles.groupError}>{error}</p>}
        <input
          minLength={3}
          maxLength={80}
          className={styles.input}
          placeholder="Title"
          value={eventTitle}
          onChange={(e) => setEventTitle(e.target.value)}
        />
      </div>
      <div className={styles.formGroup}>
        <input
          minLength={10}
          maxLength={500}
          className={styles.input}
          placeholder="Description"
          value={eventDescription}
          onChange={(e) => setEventDescription(e.target.value)}
        />
      </div>
      <div className={styles.formGroup}>
        <input
          className={styles.input}
          type="datetime-local"
          value={eventDateTime}
          min={new Date(Date.now() + 5 * 60 * 1000).toISOString().slice(0, 16)}
          max="9999-12-31T23:59"
          onChange={(e) => setEventDateTime(e.target.value)}
        />
      </div>
      <div className={styles.formGroup}>
        <div style={{ display: "flex", gap: "10px" }}>
          <button
            type="button"
            className={styles.primaryButton}
            style={
              myVote === "going"
                ? { background: "#C8892A", color: "#0A0A0F" }
                : {
                    background: "transparent",
                    border: "1px solid #2a1f0f",
                    color: "#6b5a3e",
                  }
            }
            onClick={() => setMyVote("going")}
          >
            Going
          </button>
          <button
            type="button"
            className={styles.primaryButton}
            style={
              myVote === "notgoing"
                ? { background: "#8B1A1A", color: "#E8E0CC" }
                : {
                    background: "transparent",
                    border: "1px solid #2a1f0f",
                    color: "#6b5a3e",
                  }
            }
            onClick={() => setMyVote("notgoing")}
          >
            Not Going
          </button>
        </div>
      </div>
      <button className={styles.primaryButton}>Create</button>
    </form>
  );
}
