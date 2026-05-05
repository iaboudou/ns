'use client';

import { useEffect, useState } from 'react';
import styles from '@/components/groups/styles/groups.module.css';
import HandleCreateGroup from '@/components/groups/joins/CreateGroup';
import handleFetchGroups from '@/components/groups/utils/fetchGroups';
import DisplayMyGroup from '@/components/groups/joins/DisplayMyGroups';

export default function DiscoverGroups() {
  const [groups, setGroups] = useState([]);
  const [showCreateGroupForm, setShowCreateGroupForm] = useState(false);
  const [hasMore, setHasMore] = useState(true);
  const [loading, setLoading] = useState(true);

  const fetchParams = { groups, setGroups, setLoading, setHasMore, tab: 'mine' };

  useEffect(() => {
    handleFetchGroups({ ...fetchParams, isReset: true });
  }, []);

  return (
    <>
      <button onClick={() => setShowCreateGroupForm((v) => !v)}>create group</button>
      {showCreateGroupForm && <HandleCreateGroup setGroups={setGroups} />}

      {groups.length === 0 ? (
        <p>you're not in any group. You can join one or better create your own !</p>
      ) : (
        <div className={styles.groupsList}>
          {groups.map((g) => (
            <DisplayMyGroup key={g.id} group={g} />
          ))}
        </div>
      )}

      {hasMore && (
        <button className={styles.loadMoreBtn} onClick={() => handleFetchGroups({ ...fetchParams, hasMore })} disabled={loading}>
          {loading ? 'Loading...' : 'Load More'}
        </button>
      )}
    </>
  );
}
