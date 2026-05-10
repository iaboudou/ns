"use client";

import { useEffect, useRef, useState } from "react";
import styles from "@/components/groups/styles/groups.module.css";
import cardStyles from "@/components/groups/styles/groups-cards.module.css";
import HandleCreateGroup from "@/components/groups/joins/CreateGroup";
import DisplayMyGroup from "@/components/groups/joins/DisplayMyGroups";
import handleFetchGroups from "@/components/groups/utils/fetchGroups";

export default function DiscoverGroups() {
  const [groups, setGroups] = useState([]);
  const [showCreateGroupForm, setShowCreateGroupForm] = useState(false);
  const [hasMore, setHasMore] = useState(true);
  const [loading, setLoading] = useState(true);

  const observerRef = useRef(null);

  const fetchParams = {
    groups,
    setGroups,
    setLoading,
    setHasMore,
    tab: "mine",
  };

  useEffect(() => {
    handleFetchGroups({ ...fetchParams, isReset: true });
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
      <div className={cardStyles.groupsList}>
        {showCreateGroupForm
          ? <HandleCreateGroup
              setGroups={setGroups}
              setShowCreateGroupForm={setShowCreateGroupForm}
            />
          : <button
              type="button"
              className={cardStyles.createGroupCard}
              onClick={() => setShowCreateGroupForm((v) => !v)}
            >
              <span className={cardStyles.createGroupIcon}>✦</span>
              <span className={cardStyles.createGroupLabel}>Create Group</span>
            </button>}

        {groups.map((g) => (
          <DisplayMyGroup key={g.id} group={g} />
        ))}
      </div>

      {hasMore && (
        <div ref={observerRef} className={styles.loadMoreTrigger}>
          {loading ? "Loading..." : ""}
        </div>
      )}
    </>
  );
}
