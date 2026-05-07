"use client";

import { useEffect, useRef, useState } from "react";

import styles from "@/components/groups/styles/groups.module.css";
import cardStyles from "@/components/groups/styles/groups-cards.module.css";

import handleFetchGroups from "@/components/groups/utils/fetchGroups";
import DisplayGroupInvite from "@/components/groups/group_invites/DisplayGroupInvites";

export default function DiscoverGroups() {
  const [groups, setGroups] = useState([]);
  const [hasMore, setHasMore] = useState(true);
  const [loading, setLoading] = useState(true);

  const observerRef = useRef(null);

  const fetchParams = {
    groups,
    setGroups,
    setLoading,
    setHasMore,
    tab: "invites",
  };

  useEffect(() => {
    handleFetchGroups({
      ...fetchParams,
      isReset: true,
    });
  }, []);

  useEffect(() => {
    if (!observerRef.current || !hasMore) return;

    const observer = new IntersectionObserver(
      async ([entry]) => {
        if (!entry.isIntersecting || loading) return;

        await handleFetchGroups({
          ...fetchParams,
          hasMore,
        });
      },
      {
        threshold: 0.1,
      },
    );

    observer.observe(observerRef.current);

    return () => observer.disconnect();
  }, [groups, hasMore, loading]);

  return (
    <>
      {groups.length === 0 ? (
        <p>You have no groups invites for now</p>
      ) : (
        <div className={cardStyles.groupsList}>
          {groups.map((g) => (
            <DisplayGroupInvite key={g.id} group={g} setGroups={setGroups} />
          ))}
        </div>
      )}

      {hasMore && (
        <div ref={observerRef} className={styles.loadMoreTrigger}>
          {loading ? "Loading..." : ""}
        </div>
      )}
    </>
  );
}
