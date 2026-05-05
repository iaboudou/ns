'use client';

import { useEffect, useState } from 'react';
import styles from '@/components/groups/styles/groups.module.css';
import handleFetchGroups from '@/components/groups/utils/fetchGroups';
import DisplayGroupInvite from '@/components/groups/group_invites/DisplayGroupInvites';

export default function DiscoverGroups() {
  const [groups, setGroups] = useState([]);
  const [hasMore, setHasMore] = useState(true);
  const [loading, setLoading] = useState(true);

  const fetchParams = { groups, setGroups, setLoading, setHasMore, tab: 'invites' };

  useEffect(() => {
    handleFetchGroups({ ...fetchParams, isReset: true });
  }, []);

  return (
    <>
      {groups.length === 0 ? (
        <p>You have no groups invites for now</p>
      ) : (
        <div className={styles.groupsList}>
          {groups.map((g) => (
            <DisplayGroupInvite key={g.id} group={g} setGroups={setGroups} />
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
