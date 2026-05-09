import { SendDecision } from "@/_lib/group";
import styles from "@/components/groups/styles/groups.module.css";
import cardStyles from "@/components/groups/styles/groups-cards.module.css";
import { useWebSocket } from "@/lib/UseWebsocket";

export default function DisplayGroupInvite({ group, setGroups }) {
  console.log(group);
  const { port } = useWebSocket();

  const handleDecision = async (decision) => {
    SendDecision(group.id, decision, "", group.invited_by)
      .then(() => {
        if (decision === "accepted")
          alert(`you are now a member of ${group.title}`);
        setGroups((prev) => prev.filter((g) => g.id !== group.id));
        port.postMessage({
          type: "set_notif",
          notif: { type: "group_invite", sender_id: group.invited_by },
        });
      })
      .catch((err) => alert(err.message));
  };

  return (
    <div
      className={cardStyles.groupCard}
      style={
        group.img
          ? {
              backgroundImage: `url(/pics/${group.img})`,
            }
          : { backgroundImage: `url(/pics/pub.png)` }
      }
    >
      <h3>{group.title}</h3>
      <p>{group.description}</p>
      <div className={styles.inviteActions}>
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
          Decline
        </button>
      </div>
    </div>
  );
}
