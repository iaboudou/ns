import styles from '@/components/groups/styles/singleGroup.module.css';
import { cookies } from 'next/headers';
import Link from 'next/link';
import { GetGroup } from '@/_lib/group';
import GroupHeader from '@/components/groups/header/GroupHeader';
import { redirect } from 'next/navigation';

export default async function SingleGroupLayout({ children, params }) {
  const { id } = await params;
  const cookieStore = await cookies();
  const cookie = cookieStore.toString();

  let group;

  try {
    group = await GetGroup(id, cookie);
  } catch (err) {
    //show a not found or a server error later
    redirect('/groups/joins');
  }

  return (
    <div className={styles.mainWrapper}>
      <div className={styles.up}>
        <GroupHeader group={group} id={id} />

        <div className={styles.groupFeedBtn}>
          <Link className={styles.groupFeedLink} href={`/groups/${id}/posts`}>
            Posts
          </Link>
          <Link className={styles.groupFeedLink} href={`/groups/${id}/chats`}>
            Chat
          </Link>
          <Link className={styles.groupFeedLink} href={`/groups/${id}/events`}>
            Events
          </Link>

          {group.isCreator && (
            <Link className={styles.groupFeedLink} href={`/groups/${id}/requests`}>
            Requests
            </Link>
          )}
        </div>
      </div>

      <div className={styles.down}>
        <div className={styles.content}>{children}</div>
      </div>
    </div>
  );
}
