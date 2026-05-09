"use client";

import styles from "@/components/groups/styles/singleGroup.module.css";
import { SendDecision } from "@/_lib/group";
import { useRouter } from "next/navigation";
import { useWebSocket } from "@/lib/UseWebsocket";

export default function RequestCard({ user, groupId, setRequests }) {
  const { port } = useWebSocket();
  const routeur = useRouter();
  const handleDecision = async (decision) => {
    SendDecision(groupId, decision, user.id, "")
      .then(() => {
        setRequests((prev) => prev.filter((u) => u.id !== user.id));
        port.postMessage({
          type: "set_notif",
          notif: { type: "group_request", sender_id: user.id },
        });
        routeur.refresh();
      })
      .catch((err) => alert(err.message));
  };

  return (
    <div className={styles.requestCard}>
      <p className={styles.requestUser}>
        {user.nickname ? user.nickname : user.firstname + " " + user.lastname}
      </p>

      <div className={styles.requestActions}>
        <button
          className={styles.acceptBtn}
          onClick={() => handleDecision("accepted")}
        >
          Accept
        </button>

        <button
          className={styles.rejectBtn}
          onClick={() => handleDecision("rejected")}
        >
          Refuse
        </button>
      </div>
    </div>
  );
}
