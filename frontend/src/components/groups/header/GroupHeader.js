"use client";

import { useState } from "react";
import styles from "@/components/groups/styles/singleGroup.module.css";
import InviteFriends from "@/components/groups/invite_users/HandleFriendInvite";
import AreYouSure from "../utils/AreYouSure";

export default function GroupHeader({ group, id }) {
  const [hoverTab, setHoverTab] = useState("");

  return (
    <div className={styles.groupInfo}>
      {/* Image */}
      <img
        className={styles.groupImage}
        src={
          group.img
            ? `http://localhost:4001/pics/${group.img}`
            : "http://localhost:4001/pics/pub.png"
        }
        alt={group.title}
      />

      <div className={styles.groupMeta}>
        <h1>{group.title}</h1>
        <p>{group.description}</p>
        <span className={styles.membersCount}>⚔ {group.members} members</span>

        <div className={styles.groupActions}>
          <button
            className={styles.secondaryButton}
            onClick={() => setHoverTab("invite friends")}
          >
            Invite Followers
          </button>
          {group.isCreator ? (
            <button
              className={styles.secondaryButton}
              onClick={() => setHoverTab("delete group")}
            >
              Delete Group
            </button>
          ) : (
            <button
              className={styles.secondaryButton}
              onClick={() => setHoverTab("leave group")}
            >
              Leave Group
            </button>
          )}
        </div>
      </div>

      {hoverTab === "invite friends" && (
        <div className={styles.popupOverlay} onClick={() => setHoverTab("")}>
          <div className={styles.popup} onClick={(e) => e.stopPropagation()}>
            <InviteFriends groupId={id} />
            <button
              className={styles.secondaryButton}
              onClick={() => setHoverTab("")}
            >
              Close
            </button>
          </div>
        </div>
      )}

      {hoverTab === "delete group" && (
        <AreYouSure Message={hoverTab} groupId={id} setHoverTab={setHoverTab} />
      )}
      {hoverTab === "leave group" && (
        <AreYouSure Message={hoverTab} groupId={id} setHoverTab={setHoverTab} />
      )}
    </div>
  );
}
