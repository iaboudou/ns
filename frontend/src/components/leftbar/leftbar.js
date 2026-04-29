import Link from 'next/link';
import styles from './leftbar.module.css';
import { Search, Handshake, Users, MessageSquare, Bell, CircleUser, Home } from 'lucide-react';

export default function Leftbar() {
  return (
    <nav className={styles.leftbar}>
      <div className={styles.barElementContainer}>
        <Link href="/" className={styles.buttonLink} title="Home"><Home/></Link>
        <Link href="/friends" className={styles.buttonLink} title="Friends"><Handshake/></Link>
        <Link href="/groups/joins" className={styles.buttonLink} title="Groups"><Users/></Link>
        <Link href="/chat" className={styles.buttonLink} title="Chat"><MessageSquare/></Link>
        <Link href="/notifications" className={styles.buttonLink} title="Notifications"><Bell/></Link>
        <Link href="/profile/me" className={styles.buttonLink} title="Profile"><CircleUser/></Link>
      </div>
    </nav>
  );
}