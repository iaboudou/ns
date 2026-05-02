import { sendGroupInvite } from '@/_lib/group';
import styles from '@/components/groups/styles/singleGroup.module.css';

export default function DisplayUser({ u, groupId, setUsers }) {
  return (
    <div className={styles.userItem}>
      <div className={styles.userLeft}>
        <img src={u.profile_image ? `http://localhost:4001/pics/${u.profile_image}` : '/avatar.jpg'} className={styles.userAvatar} alt="" />

        <div className={styles.userText}>
          <div className={styles.userNickname}>{u.nickname}</div>
          <div className={styles.userName}>
            {u.firstname} {u.lastname}
          </div>
        </div>
      </div>

      <button
        className={styles.inviteBtn}
        onClick={() => {
          sendGroupInvite(groupId, u.id)
            .then(() => {
              setUsers((prev) => prev.filter((user) => user.id != u.id));
              alert('invite sent');
            })
            .catch((err) => alert(err.message));
        }}
      >
        Invite
      </button>
    </div>
  );
}
