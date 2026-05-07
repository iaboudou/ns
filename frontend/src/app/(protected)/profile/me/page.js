"use client";
const BASE = "http://localhost:4001";
import styles from "./me.module.css";
import { useState, useEffect } from "react";
import Link from "next/link";
import FeedClient from "@/app/(protected)/feed/feedclient";
import {
  fetchSwitchPrivacy,
  fetchFollowers,
  fetchFollowing,
  fetchRequests,
  ManageFollow,
} from "../actions";
import { about } from "../about/about";
import { useRouter, useSearchParams } from "next/navigation";
import { fetchPersonalInfo } from "@/_lib/personal_info";
import { Lock, User } from "lucide-react";

export default function ProfilePage() {
  const router = useRouter();

  //
  const searchParams = useSearchParams();
  const s = searchParams.get("section") || "posts";

  //
  const [user, setUser] = useState(null);
  const [section, setSection] = useState(s);
  const [followers, setFollowers] = useState([]);
  const [following, setFollowing] = useState([]);
  const [requests, setRequests] = useState([]);

  // about
  useEffect(() => {
    let uuid = JSON.parse(localStorage.getItem("user") || "null")?.ID;
    if (!uuid) return;
    (async () => {
      let data = await fetchPersonalInfo(uuid);
      setUser(data);
    })();
  }, []);

  // followers & following
  useEffect(() => {
    if (!user || !user?.id) return;
    if (section === "followers") {
      (async () => {
        const data = await fetchFollowers(user.id);
        setFollowers(data);
      })();
    } else if (section === "following") {
      (async () => {
        const data = await fetchFollowing(user.id);

        setFollowing(data);
      })();
    } else if (section === "requests") {
      (async () => {
        const data = await fetchRequests(user.id);

        setRequests(data);
      })();
    }
  }, [section, user?.id]);

  // loading
  if (!user) return <div>Loading...</div>;

  // account privacy
  async function handleSwitchAccountPrivacy() {
    let switched = await fetchSwitchPrivacy();
    if (switched) {
      setUser((prev) => ({ ...prev, is_public: !user.is_public }));
    }
  }

  let main_content;
  const renderUserList = (users, type) => {
    if (users.length === 0)
      return <div className={styles.noData}>No {type} yet.</div>;
    return (
      <div className={styles.userListContainer}>
        {users.map((u) => (
          <div key={u.id} className={styles.userWrapper}>
            <Link href={`/profile/${u.id}`} className={styles.userListItem}>
              {u.profile_image ? (
                <img
                  src={`${BASE}/pics/${u.profile_image}`}
                  className={styles.smallAvatar}
                />
              ) : (
                <User className={styles.placeholderIcon} />
              )}
              <span className={styles.userName}>
                {u.firstname} {u.lastname}
              </span>
            </Link>
            {type === "requests" && (
              <div className={styles.requestActions}>
                <button
                  className={`${styles.actionBtn} ${styles.acceptBtn}`}
                  onClick={async () =>
                    ManageFollow(u.id, "accepted")
                      .then(() => {
                        setRequests((prev) =>
                          prev.filter(
                            (oldRequester) => oldRequester.id !== u.id,
                          ),
                        );
                      })
                      .catch(() => {})
                  }
                >
                  Accept
                </button>
                <button
                  className={`${styles.actionBtn} ${styles.rejectBtn}`}
                  onClick={() =>
                    ManageFollow(u.id, "rejected")
                      .then(() => {
                        setRequests((prev) =>
                          prev.filter(
                            (oldRequester) => oldRequester.id !== u.id,
                          ),
                        );
                      })
                      .catch(() => {})
                  }
                >
                  Reject
                </button>
              </div>
            )}
          </div>
        ))}
      </div>
    );
  };

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
    case "followers":
      main_content = renderUserList(followers, "followers");
      break;
    case "following":
      main_content = renderUserList(following, "following");
      break;
    case "requests":
      main_content = renderUserList(requests, "requests");
  }

  //
  const imageURL = user.profile_image;
  const fullImageURL = imageURL ? `${BASE}/pics/${imageURL}` : "";
  return (
    <div className={styles.FATHER}>
      <div className={styles.pageWrapper}>
        {/* HEADER */}
        <header className={styles.header}>
          <div className={styles.coverPhoto}>
            <p className={styles.infoText}>
              {user.is_public
                ? "Your account is currently Public"
                : "Your account is currently Private"}{" "}
              <Lock />
            </p>
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
                <h5 className={styles.profileNickname}>{user.nickname}</h5>
              </div>
              <button
                className={styles.toggleBtn}
                onClick={handleSwitchAccountPrivacy}
              >
                {user.is_public ? "Switch to Private" : "Switch to Public"}
              </button>
            </div>
            <p className={styles.profileBio}>{user.aboutme}</p>
          </div>
        </header>

        {/* NAVIGATION */}
        {
          <>
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
                  className={`${styles.navBtn} ${section === "followers" ? styles.navBtnActive : ""}`}
                  onClick={() => {
                    setSection("followers");
                    router.push(`?section=followers`);
                  }}
                >
                  Followers
                </button>
                <button
                  className={`${styles.navBtn} ${section === "following" ? styles.navBtnActive : ""}`}
                  onClick={() => {
                    setSection("following");
                    router.push(`?section=following`);
                  }}
                >
                  Following
                </button>

                {!user.is_public && (
                  <button
                    className={`${styles.navBtn} ${section === "requests" ? styles.navBtnActive : ""}`}
                    onClick={() => {
                      setSection("requests");
                      router.push(`?section=requests`);
                    }}
                  >
                    requests
                  </button>
                )}

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
          </>
        }
      </div>

      {/* MAIN */}
      <div className={styles.MAIN}> {main_content} </div>
    </div>
  );
}
