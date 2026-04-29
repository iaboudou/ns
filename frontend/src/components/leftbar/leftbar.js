import Link from 'next/link';
import styles from './leftbar.module.css';

export default function Leftbar() {
  return (
    <nav className={styles.leftbar}>
      <div className={styles.barElementContainer}>
        <Link href="/" className={styles.buttonLink}></Link>
        <Link href="/friends" className={styles.buttonLink}></Link>
        <Link href="/groups/joins" className={styles.buttonLink}></Link>
        <Link href="/chat" className={styles.buttonLink}></Link>
        <Link href="/notifications" className={styles.buttonLink}></Link>
        <Link href={`/profile/me`} className={styles.buttonLink}></Link>
      </div>
    </nav>
  );
}
