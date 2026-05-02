'use client';

import { useState } from 'react';
import styles from '@/components/groups/styles/singleGroup.module.css';
import InviteFriends from '@/components/groups/invite_users/HandleFriendInvite';
import AreYouSure from '../utils/AreYouSure';

export default function GroupHeader({ group, id }) {
  const [hoverTab, setHoverTab] = useState('');

  return (
    <div className={styles.groupInfo}>
      <h1>{group.title}</h1>
      <p>{group.description}</p>
      <p>Members: {group.members}</p>

      <>
        <button className={styles.secondaryButton} onClick={() => setHoverTab('invite friends')}>
          invite friends
        </button>
        {group.isCreator ? (
          <button className={styles.secondaryButton} onClick={() => setHoverTab('delete group')}>
            delete group
          </button>
        ) : (
          <button className={styles.secondaryButton} onClick={() => setHoverTab('leave group')}>
            leave group
          </button>
        )}
      </>

      {hoverTab === 'invite friends' && (
        <div className={styles.popupOverlay} onClick={() => setHoverTab('')}>
          <div className={styles.popup} onClick={(e) => e.stopPropagation()}>
            <InviteFriends groupId={id} />
            <button className={styles.secondaryButton} onClick={() => setHoverTab('')}>
              Close
            </button>
          </div>
        </div>
      )}

      {hoverTab === 'delete group' && <AreYouSure Message={hoverTab} groupId={id} setHoverTab={setHoverTab} />}
      {hoverTab === 'leave group' && <AreYouSure Message={hoverTab} groupId={id} setHoverTab={setHoverTab} />}
    </div>
  );
}
