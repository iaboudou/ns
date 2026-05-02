import styles from '@/components/groups/styles/groups.module.css';
import Link from 'next/link';

export default function DisplayMyGroups({ group }) {
  return (
    <Link href={`/groups/${group.id}/posts`} className={styles.groupCard}>
      <h3>{group.title}</h3>
      <p>{group.description}</p>
    </Link>
  );
}
