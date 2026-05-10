"use client";

import { use } from "react";
import { useEffect, useState } from "react";
import EventCard from "@/components/groups/events/DisplayEvents";
import CreateEvents from "@/components/groups/events/createEvents";
import styles from "@/components/groups/styles/singleGroup.module.css";
import handleFetchData from "@/components/groups/utils/FetchData";

export default function GroupEvents({ params }) {
  const { id } = use(params);
  const [events, setEvents] = useState([]);
  const [showCreateEvent, setShowCreateEvent] = useState(false);
  const [hasMore, setHasMore] = useState(true);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    handleFetchData(events, setEvents, setLoading, setHasMore, "events", id);
  }, [id]);

  return (
    <>
      <button
        className={styles.secondaryButton}
        onClick={() => setShowCreateEvent((v) => !v)}
      >
        Create Event
      </button>

      <div style={{ display: "flex", flexDirection: "column", gap: "16px" }}>
        {showCreateEvent && <CreateEvents groupId={id} setEvents={setEvents} />}

        {events.length === 0
          ? <p>There is no event for now but you can create your own !</p>
          : events.map((e) => <EventCard key={e.id} event={e} groupId={id} />)}
      </div>

      {hasMore && (
        <button
          className={styles.loadMoreBtn}
          onClick={() =>
            handleFetchData(
              events,
              setEvents,
              setLoading,
              setHasMore,
              "events",
              id,
            )
          }
          disabled={loading}
        >
          {loading ? "Loading..." : "Load More"}
        </button>
      )}
    </>
  );
}
