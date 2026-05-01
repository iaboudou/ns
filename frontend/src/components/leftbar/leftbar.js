"use client";

import Link from 'next/link';
import { useRouter } from 'next/navigation';
import styles from './leftbar.module.css';
import { Handshake, Users, MessageSquare, Bell, CircleUser, Home, LogOut } from 'lucide-react';

export default function Leftbar() {
  const router = useRouter();

  const handleLogout = async () => {
    try {
      await fetch('http://localhost:4001/api/logout', {
        method: 'POST',
        credentials: 'include',
      });
      localStorage.removeItem('user');
      router.push('/login');
    } catch (error) {
      localStorage.removeItem('user');
      router.push('/login');
    }
  };

  return (
    <nav className={styles.leftbar}>
      <div className={styles.barElementContainer}>
        <Link
          href="/"
          className={styles.buttonLink}
          onClick={() => window.scrollTo({ top: 0, behavior: 'smooth' })}
          title="Home">
          <Home />
        </Link>
        <Link href="/friends" className={styles.buttonLink} title="Friends"><Handshake /></Link>
        <Link href="/groups/joins" className={styles.buttonLink} title="Groups"><Users /></Link>
        <Link href="/chat" className={styles.buttonLink} title="Chat"><MessageSquare /></Link>
        <Link href="/notifications" className={styles.buttonLink} title="Notifications"><Bell /></Link>
        <Link href="/profile/me" className={styles.buttonLink} title="Profile"><CircleUser /></Link>
        <button 
          type="button"
          onClick={handleLogout} 
          className={styles.buttonLink} 
          title="Logout"
          aria-label="Logout"
        >
          <LogOut />
        </button>
      </div>
    </nav>
  );
}