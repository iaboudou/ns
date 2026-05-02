"use client";

import styles from "@/components/groups/styles//groups.module.css";
import Link from "next/link";

export default function GroupLayout({ children }) {
  return (
    <div className={styles.groupsWrapper}>
      <div className={styles.tabs}>
        <Link className={styles.tabButton} href="/groups/joins">
          My groups
        </Link>

        <Link className={styles.tabButton} href="/groups/discover">
          Discover
        </Link>

        <Link className={styles.tabButton} href="/groups/invites">
          Invitations
        </Link>
      </div>

      <div className={styles.content}>{children}</div>
    </div>
  );
}
