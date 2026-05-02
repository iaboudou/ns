import { SendGroupRequest } from '@/_lib/group';
import styles from '@/components/groups/styles/groups.module.css';

export default function DisplayNewGroups({ group, setGroups }) {
  return (
    <div className={styles.groupCard}>
      <h3>{group.title}</h3>
      <p>{group.description}</p>
      <button
        className={styles.primaryButton}
        onClick={async () => {
          SendGroupRequest(group.id)
            .then(() => setGroups((prev) => prev.filter((oldGgroup) => oldGgroup.id !== group.id)))
            .catch((err) => alert(err.message));
        }}
      >
        ask to join
      </button>
    </div>
  );
}
