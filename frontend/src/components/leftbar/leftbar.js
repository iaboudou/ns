"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import styles from "./leftbar.module.css";
import {
  Handshake,
  Users,
  MessageSquare,
  CircleUser,
  Home,
  LogOut,
  Bell,
} from "lucide-react";

import { useWebSocket } from "@/lib/UseWebsocket";

export default function Leftbar() {
  const router = useRouter();
  const { port, unreadNotifCount } = useWebSocket();
  const [loading, setLoading] = useState(false);

  const handleLogout = async () => {
    if (loading) return;
    setLoading(true);
    try {
      await fetch("http://localhost:4001/api/logout", {
        method: "POST",
        credentials: "include",
      });
      localStorage.removeItem("user");
      router.push("/login");
    } catch (error) {
      localStorage.removeItem("user");
      router.push("/login");
    } finally {
      setLoading(false);
      if (port) port.postMessage({ type: "logout" });
    }
  };

  useEffect(() => {
    if (port) {
      port.postMessage({
        type: "send",
        payload: {
          type: "unread_notif",
        },
      });
    }
  }, [port]);

  const updateLastseen = async () => {
    const resp = await fetch("http://localhost:4001/api/last-notif", {
      method: "PATCH",
      credentials: "include",
    });

    if (!resp.ok) throw new Error("failed to set last read notif");
  };

  return (
    <nav className={styles.leftbar}>
      <div className={styles.barElementContainer}>
        <Link
          href="/"
          className={styles.buttonLink}
          onClick={() => window.scrollTo({ top: 0, behavior: "smooth" })}
          title="Home"
        >
          <Home />
        </Link>

        <Link href="/groups/joins" className={styles.buttonLink} title="Groups">
          <Users />
        </Link>

        <Link href="/chat" className={styles.buttonLink} title="Chat">
          <MessageSquare />
        </Link>

        <Link
          href="/notifications"
          className={styles.buttonLink}
          title="Notifications"
          onClick={() => {
            updateLastseen()
              .then(() => {
                port.postMessage({ type: "notifs_seen" });
              })
              .catch(console.error);
          }}
        >
          <Bell />
          {unreadNotifCount > 0 && (
            <div className={styles.NOTIF}>{unreadNotifCount}</div>
          )}
        </Link>

        <Link href="/profile/me" className={styles.buttonLink} title="Profile">
          <CircleUser />
        </Link>
        <button
          type="button"
          onClick={handleLogout}
          className={styles.buttonLink}
          title="Logout"
          aria-label="Logout"
          disabled={loading}
        >
          {loading ? "..." : <LogOut />}
        </button>
      </div>
    </nav>
  );
}
