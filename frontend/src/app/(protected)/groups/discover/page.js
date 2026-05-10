"use client";

import { useEffect, useRef, useState } from "react";

import styles from "@/components/groups/styles/groups.module.css";
import cardStyles from "@/components/groups/styles/groups-cards.module.css";

import handleFetchGroups from "@/components/groups/utils/fetchGroups";
import DisplayNewGroup from "@/components/groups/discover/DisplayNewGroups";

export default function DiscoverGroups() {
  const [groups, setGroups] = useState([]);
  const [hasMore, setHasMore] = useState(true);
  const [loading, setLoading] = useState(true);
  const [search, setSearch] = useState("");

  const observerRef = useRef(null);

  const fetchParams = {
    groups,
    setGroups,
    setLoading,
    setHasMore,
    tab: "discover",
  };

  useEffect(() => {
    const timer = setTimeout(() => {
      setHasMore(true);

      handleFetchGroups({
        ...fetchParams,
        search,
        isReset: true,
      });
    }, 400);

    return () => clearTimeout(timer);
  }, [search]);

  useEffect(() => {
    if (!observerRef.current || !hasMore) return;

    const observer = new IntersectionObserver(
      async ([entry]) => {
        if (!entry.isIntersecting || loading) return;

        await handleFetchGroups({
          ...fetchParams,
          search,
          hasMore,
        });
      },
      {
        threshold: 0.1,
      },
    );

    observer.observe(observerRef.current);

    return () => observer.disconnect();
  }, [groups, hasMore, loading, search]);

  return (
    <>
      <input
        className={styles.input}
        type="text"
        maxLength={60}
        placeholder="which group are you looking for ?"
        value={search}
        onChange={(e) => setSearch(e.target.value)}
      />

      {groups.length === 0
        ? <p>
            There is not any group for now. You can come back later or better
            create your own !
          </p>
        : <div className={cardStyles.groupsList}>
            {groups.map((g) => (
              <DisplayNewGroup key={g.id} group={g} setGroups={setGroups} />
            ))}
          </div>}

      {hasMore && (
        <div ref={observerRef} className={styles.loadMoreTrigger}>
          {loading ? "Loading..." : ""}
        </div>
      )}
    </>
  );
}
