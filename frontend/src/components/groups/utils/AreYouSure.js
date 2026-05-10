"use client";

import { DeleteGroup, LeaveGroup } from "@/_lib/group";
import styles from "@/components/groups/styles/singleGroup.module.css";
import { useRouter } from "next/navigation";

export default function AreYouSure({ Message, groupId, setHoverTab }) {
  const router = useRouter();

  return (
    <div className={styles.popupOverlay}>
      <div className={styles.popup}>
        <p>Are you sure you want to {Message}?</p>

        <div className={styles.buttonswrapper}>
          <button
            onClick={(e) => {
              e.preventDefault();
              Message === "delete group"
                ? DeleteGroup(groupId).catch((err) => alert(err.message))
                : LeaveGroup(groupId).catch((err) => {
                    alert(err.message);
                    router.replace("/groups/joins");
                  });
              router.replace("/groups/joins");
            }}
          >
            Yes
          </button>
          <button
            onClick={(e) => {
              e.preventDefault();
              setHoverTab("");
            }}
          >
            No
          </button>
        </div>
      </div>
    </div>
  );
}
