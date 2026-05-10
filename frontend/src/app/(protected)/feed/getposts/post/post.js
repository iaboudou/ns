import styles from "./post.module.css";
import Comments from "@/app/(protected)/feed/getposts/post/getcomment/getcomments";
import CreateComment from "@/app/(protected)/feed/getposts/post/createcomment/createcomment";
import { fetchComments, set_comment_in_state } from "./actions";
import { timeAgo } from "@/_lib/timeago";
import Link from "next/link";
import { usePathname } from "next/navigation";
import { MessageCircle, Users } from "lucide-react";
import { useWebSocket } from "@/lib/UseWebsocket";

export default function Post({ POST }) {
  const { myInfo } = useWebSocket();
  //post image
  const imageURL = POST.post.image_url;
  const fullImageURL = imageURL || null;

  // profile image
  const profileimage = POST.post.profile_image;
  const fullprofileimage = profileimage ? `/pics/${profileimage}` : "";

  if (!POST.post) return null;

  //  fetch 10 comments from the server for a post starting at its offset,
  async function get10Comments(post_id) {
    let data = await fetchComments(post_id, POST.post?.offset || 0);
    if (!data?.length) return;

    if (Array.isArray(data))
      set_comment_in_state(POST.setState, post_id, data, POST.state);
  }

  let COMMENTS = {
    comments: POST.post.comments || [],
    number_of_comments: POST.post.number_of_comments,
    get10Comments: () => get10Comments(POST.post.id),
  };

  let path = usePathname();
  let id = myInfo?.id || "";

  return (
    <div className={styles.postcard}>
      <div className={styles.header}>
        {
          <Link
            href={
              id != POST.post.user_id
                ? `/profile/${POST.post.user_id}`
                : `/profile/me`
            }
          >
            {fullprofileimage
              ? <img src={fullprofileimage} className={styles.profileImg} />
              : <Users className={styles.placeholderIcon} />}
          </Link>
        }

        {/*  */}
        <div className={styles.meta}>
          <div className={styles.nameRow}>
            <span className={styles.name}>
              <Link
                className={styles.linkprofilename}
                href={
                  id != POST.post.user_id
                    ? `/profile/${POST.post.user_id}`
                    : `/profile/me`
                }
              >
                {POST.post.firstname} {POST.post.lastname}
              </Link>
            </span>
            {!path.includes("/group") && (
              <span className={styles.handle}>@{POST.post.privacy ?? ""}</span>
            )}
          </div>
          <span className={styles.date}>{timeAgo(POST.post.created_at)}</span>
        </div>
      </div>

      {POST.post.content && (
        <p className={styles.content}>{POST.post.content}</p>
      )}

      {fullImageURL && <img src={fullImageURL} className={styles.postImage} />}

      <div className={styles.postActions}>
        <button
          className={styles.actionBtn}
          type="button"
          onClick={() => {
            POST.setOpenComments(POST.post.id);
            if (!POST.post?.firs_10_comments_already_fetched) {
              POST.post.firs_10_comments_already_fetched = true;
              get10Comments(POST.post.id);
            }
          }}
        >
          <MessageCircle size={14} />
          opinions ({POST.post.number_of_comments || 0})
        </button>
      </div>

      {POST.openComments[POST.post.id] && (
        <div className={styles.commentsSection}>
          <CreateComment
            post={POST.post}
            onCommentCreated={POST.onCommentCreated}
          />
          <Comments COMMENTS={COMMENTS} />
        </div>
      )}
    </div>
  );
}
