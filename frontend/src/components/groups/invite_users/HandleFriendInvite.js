import { getUsers } from '@/_lib/group';
import styles from '@/components/groups/styles/singleGroup.module.css';
import style from '@/components/groups/styles/groups.module.css';
import { useEffect, useState } from 'react';
import DisplayUser from './DisplayUsers';

export default function InviteFriends({ groupId }) {
  const [users, setUsers] = useState([]);
  const [hasMore, setHasmore] = useState(true);
  const [loading, setLoading] = useState(false);
  const [search, setSearch] = useState('');

  const handleFetchUsers = async (currentSearch, isReset = false) => {
    if (!hasMore && !isReset) return;

    const lastUser = isReset ? undefined : users.at(-1);
    const lastTime = lastUser?.created_at;
    const lastId = lastUser?.id;

    try {
      setLoading(true);
      const newUsers = await getUsers(lastTime, lastId, currentSearch, groupId);

      if (newUsers.length < 10) setHasmore(false);

      setUsers(
        isReset
          ? newUsers
          : (prev) => {
              const map = new Map(prev.map((u) => [u.id, u]));
              newUsers.forEach((u) => map.set(u.id, u));
              return Array.from(map.values());
            },
      );
    } catch (err) {
      alert(err.message);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    const timer = setTimeout(() => {
      setHasmore(true);
      handleFetchUsers(search, true);
    }, 400);

    return () => clearTimeout(timer);
  }, [search]);

  return (
    <>
      <p>Invite Friends</p>
      <input className={style.input} type="text" placeholder="looking for someone ?" value={search} onChange={(e) => setSearch(e.target.value)} />
      <div className={styles.userList}>
        {users.map((u) => (
          <DisplayUser key={u.id} u={u} groupId={groupId} setUsers={setUsers} />
        ))}
        {hasMore && (
          <button className={styles.loadMoreBtn} onClick={() => handleFetchUsers(search)} disabled={loading}>
            {loading ? 'Loading...' : 'Load More'}
          </button>
        )}
      </div>
    </>
  );
}
