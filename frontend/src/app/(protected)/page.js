import styles from "@/app/(protected)/page.module.css";
import FeedClient from "@/app/(protected)/feed/feedclient";
import FollowSuggestions from "@/app/(protected)/feed/followsuggestions/followsuggestions";

export default function Home() {
  return (
    <div className={styles.WRAPPER}>
      <div className={styles.MIDDLE}>
        <FeedClient></FeedClient>
      </div>
      <div className={styles.RIGHT}>
        <FollowSuggestions />
      </div>
    </div>
  );
}
