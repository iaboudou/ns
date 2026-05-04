"use client";

import { useState } from "react";
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
import { usePathname } from "next/navigation";
import { useEffect } from "react";

export default function Leftbar() {
  const router = useRouter();
  const { port, unreadNotifCount, setUnreadNotifCount } = useWebSocket();
  const [loading, setLoading] = useState(false);
  const pathname = usePathname();

  // reset badge when user visits the notifications page
  useEffect(() => {
    if (pathname === "/notifications") {
      setUnreadNotifCount(0);
    }
  }, [pathname]);

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
        <Link href="/friends" className={styles.buttonLink} title="Friends">
          <Handshake />
        </Link>
        <Link href="/groups/joins" className={styles.buttonLink} title="Groups">
          <Users />
        </Link>
        <Link href="/chat" className={styles.buttonLink} title="Chat">
          <MessageSquare />
        </Link>
        
        <Link href="/notifications" className={styles.buttonLink} style={{ position: "relative" }} title="Notifications">
          <Bell />
          {unreadNotifCount > 0 && (
            <span style={{
              position: "absolute",
              top: "-6px",
              right: "-6px",
              background: "#ef4444",
              color: "white",
              fontSize: "10px",
              fontWeight: "bold",
              borderRadius: "999px",
              padding: "1px 5px",
              lineHeight: "1.4",
              minWidth: "16px",
              textAlign: "center",
            }}>
              {unreadNotifCount}
            </span>
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
