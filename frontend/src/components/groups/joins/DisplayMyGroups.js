import cardStyles from "@/components/groups/styles/groups-cards.module.css";
import Link from "next/link";

export default function DisplayMyGroup({ group }) {
  return (
    // DisplayMyGroup.jsx
    <Link
      href={`/groups/${group.id}/posts`}
      className={cardStyles.groupCard}
      style={
        group.img
          ? {
              backgroundImage: `url(http://localhost:4001/pics/${group.img})`,
            }
          : {
              backgroundImage: `url(http://localhost:4001/pics/pub.png)`,
            }
      }
    >
      <h3>{group.title}</h3>
      <p>{group.description}</p>
    </Link>
  );
}
