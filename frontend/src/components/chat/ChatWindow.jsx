import styles from "./ChatWindow.module.css";

export default function ChatWindow({ isSelected, userId }) {
  if (isSelected) {
    return (
      <div>
        <h1>User id: {userId}</h1>
      </div>
    );
  }

  return (
    <div className={styles.emptyShell}>
      <div className={styles.emptyState}>
        <h1>you did not select the chat yet</h1>
      </div>
    </div>
  );
}
