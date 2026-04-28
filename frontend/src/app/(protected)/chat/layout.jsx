import styles from "./layout.module.css";
import UsersList from "@/components/chat/UsersList";

export default function Layout({ children }) {
  return (
    <div className={styles.container}>
      <div className={styles.left}>
        <UsersList />
      </div>

      <div className={styles.right}>
        {children}
      </div>
    </div>
  );
}