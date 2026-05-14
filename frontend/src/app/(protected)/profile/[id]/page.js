"use client";
import styles from "./profile.module.css";
import { useState, useEffect } from "react";
import FeedClient from "@/app/(protected)/feed/feedclient";
import { fetchPersonalInfo } from "@/_lib/personal_info";
import { about } from "../about/about";
import { useRouter, useSearchParams } from "next/navigation";
import { User, Lock } from "lucide-react";
import { FollowUser } from "@/app/(protected)/feed/followsuggestions/actions";

export default function ProfilePage() {
  const router = useRouter();

  // get the params from url
  const searchParams = useSearchParams();
  const s = searchParams.get("section") || "posts";

  // states
  const [user, setUser] = useState(null);
  const [section, setSection] = useState(s);
  const [interactionStatus, setInteractionStatus] = useState();
  const [hasAccess, setAccess] = useState(false);

  useEffect(() => {
    const path = window.location.pathname;
    const uuid = path.split("/")[2];
    if (!uuid) return;

    (async () => {
      let data = await fetchPersonalInfo(uuid);
      setUser(data);
      setAccess(data.is_public || data.is_freind);
      if (data?.interaction_status) {
        setInteractionStatus(data.interaction_status);
      }
    })();
  }, []);

  useEffect(() => {
    if (!user?.id) return;
  }, [section, user?.id]);

  const handleFollow = async (e) => {
    e.preventDefault();
    if (!user?.id) return;
    const message = await FollowUser(user.id);
    console.log(message);
    if (!message) return;

    if (message === "follow have been successfully") {
      setAccess(true);
      setInteractionStatus("following");
      setUser((prev) => ({ ...prev, is_freind: true }));
    } else if (message === "request have been sent") {
      // window.location.reload();
      setInteractionStatus("requested");
      // router.replace(window.location.pathname)
      setAccess(false);
    } else if (message === "follow deleted") {
      setInteractionStatus("none");
      setUser((prev) => ({ ...prev, is_freind: false }));
    } else if (message === "follow request deleted") {
      setInteractionStatus("none");
      setUser((prev) => ({ ...prev, is_freind: false }));
      setAccess(false);
    }
  };

  // loading
  if (!user) return <div className={styles.FATHER}>Loading...</div>;

  let main_content;

  let INFO = {
    section,
    profileId: user?.id,
  };

  switch (section) {
    case "about":
      main_content = about(user);
      break;
    case "posts":
      main_content = <FeedClient INFO={INFO} />;
      break;
  }

  //
  //
  const imageURL = user.profile_image;
  const fullImageURL = imageURL ? `/pics/${imageURL}` : "";
  return (
    <div className={styles.FATHER}>
      <div className={styles.pageWrapper}>
        {/* HEADER */}
        <header className={styles.header}>
          <div className={styles.coverPhoto}>
            <button className={styles.infoText} onClick={handleFollow}>
              {interactionStatus === "following" || user.is_freind
                ? "unfollow"
                : interactionStatus === "requested"
                  ? "requested"
                  : !user.is_public
                    ? "request"
                    : "follow"}
            </button>
          </div>

          <div className={styles.profileInfo}>
            {fullImageURL ? (
              <img className={styles.profileAvatar} src={fullImageURL} />
            ) : (
              <User className={styles.profileAvatar} />
            )}
            <div className={styles.nameandprivacybuttoncontainer}>
              <div className={styles.flnname}>
                <h1 className={styles.profileName}>
                  {user.firstname + " " + user.lastname}
                </h1>
                {hasAccess && (
                  <h5 className={styles.profileNickname}>{user.nickname}</h5>
                )}
              </div>
              <span className={styles.s}></span>
            </div>
            {hasAccess && <p className={styles.profileBio}>{user.aboutme}</p>}
          </div>
        </header>

        {/* NAVIGATION */}
        {hasAccess && (
          <nav className={styles.NAV}>
            <ul className={styles.navbar}>
              <button
                className={`${styles.navBtn} ${section === "posts" ? styles.navBtnActive : ""}`}
                onClick={() => {
                  setSection("posts");
                  router.push(`?section=posts`);
                }}
              >
                Posts
              </button>
              <button
                className={`${styles.navBtn} ${section === "about" ? styles.navBtnActive : ""}`}
                onClick={() => {
                  setSection("about");
                  router.push(`?section=about`);
                }}
              >
                About
              </button>
            </ul>
          </nav>
        )}
      </div>

      {/* MAIN */}
      <div className={styles.MAIN}>
        {hasAccess ? (
          main_content
        ) : (
          <div className={styles.privateMessage}>
            <Lock size={48} className={styles.lockIcon} />
            <h2>This Profile is Private</h2>
          </div>
        )}
      </div>
    </div>
  );
}
