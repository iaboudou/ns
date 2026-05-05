import { SendDecision } from '@/_lib/group';
import styles from '@/components/groups/styles/groups.module.css';

export default function DisplayGroupInvite({ group, setGroups }) {
  const handleDecision = async (decision) => {
    SendDecision(group.id, decision)
      .then(() => {
        if (decision === 'accepted') alert(`you are now a member of ${group.title}`);
        setGroups((prev) => prev.filter((oldGgroups) => oldGgroups.id !== group.id));
      })
      .catch((err) => alert(err.message));
  };

  return (
    <div className={styles.groupCard}>
      <h3>{group.title}</h3>
      <p>{group.description}</p>

      <div className={styles.inviteActions}>
        <button className={styles.acceptBtn} onClick={() => handleDecision('accepted')}>
          Accept
        </button>
        <button className={styles.rejectBtn} onClick={() => handleDecision('rejected')}>
          Decline
        </button>
      </div>
    </div>
  );
}
