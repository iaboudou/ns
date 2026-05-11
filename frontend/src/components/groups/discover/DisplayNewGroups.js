import { SendGroupRequest } from "@/_lib/group";
import styles from "@/components/groups/styles/groups.module.css";
import cardStyles from "@/components/groups/styles/groups-cards.module.css";

export default function DisplayNewGroups({ group, setGroups }) {
  return (
    <div
      className={cardStyles.groupCard}
      style={
        group.img
          ? {
              backgroundImage: `url(/pics/${group.img})`,
            }
          : { backgroundImage: `url(/static/pub.png)` }
      }
    >
      <h3>{group.title}</h3>
      <p>{group.description}</p>
      <button
        className={`${styles.primaryButton} ${cardStyles.askToJoinBtn}`}
        onClick={async () => {
          SendGroupRequest(group.id)
            .then(() =>
              setGroups((prev) => prev.filter((g) => g.id !== group.id)),
            )
            .catch((err) => console.error(err.message));
        }}
      >
        ask to join
      </button>
    </div>
  );
}
