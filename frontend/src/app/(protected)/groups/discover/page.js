'use client';

import { useEffect, useState } from 'react';
import styles from '@/components/groups/styles/groups.module.css';
import handleFetchGroups from '@/components/groups/utils/fetchGroups';
import DisplayNewGroups from '@/components/groups/discover/DisplayNewGroups';

export default function DiscoverGroups() {
  const [groups, setGroups] = useState([]);
  const [hasMore, setHasMore] = useState(true);
  const [loading, setLoading] = useState(true);
  const [search, setSearch] = useState('');

  const fetchParams = { groups, setGroups, setLoading, setHasMore, tab: 'discover' };

  useEffect(() => {
    const timer = setTimeout(() => {
      setHasMore(true);
      handleFetchGroups({ ...fetchParams, search, isReset: true });
    }, 400);
    return () => clearTimeout(timer);
  }, [search]);

  return (
    <>
      <input className={styles.input} type="text" maxLength={60} placeholder="which group are you looking for ?" value={search} onChange={(e) => setSearch(e.target.value)} />
      {groups.length === 0 ? (
        <p>There is not any group for now. You can comme back later or better create your own !</p>
      ) : (
        <div className={styles.groupsList}>
          {groups.map((g) => (
            <DisplayNewGroups key={g.id} group={g} setGroups={setGroups} />
          ))}
        </div>
      )}
      {hasMore && <button onClick={() => handleFetchGroups({ ...fetchParams, search, hasMore })}> {loading ? 'Loading...' : 'Load More'}</button>}
    </>
  );
}
