"use client"
import styles from "./profile.module.css"
import { useState, useEffect } from "react";
import FeedClient from "@/app/(protected)/feed/feedclient";
import { fetchFollowers, fetchFollowing } from "../actions";
import { fetchPersonalInfo } from "@/_lib/personal_info"
import { about } from "../about/about";
import { useRouter, useSearchParams } from "next/navigation"
import { User, Lock } from 'lucide-react';


const BASE = "http://localhost:4001"
export default function ProfilePage() {

    const router = useRouter()

    // get the params from url
    const searchParams = useSearchParams()
    const s = searchParams.get("section") || "posts"

    // states
    const [user, setUser] = useState(null);
    const [section, setSection] = useState(s);
    const [followers, setFollowers] = useState([]);
    const [following, setFollowing] = useState([]);

    useEffect(() => {
        const path = window.location.pathname
        const uuid = path.split("/")[2]
        if (!uuid) return

        (async () => {
            let data = await fetchPersonalInfo(uuid)
            setUser(data)
        })()
    }, []);

    useEffect(() => {
        if (!user?.id) return

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
        }
    }, [section, user?.id]);

    // loading
    if (!user) return <div className={styles.FATHER}>Loading...</div>;

    let main_content

    let INFO = {
        section,
        profileId: user?.id
    }

    switch (section) {
        case "about":
            main_content = about(user);
            break;
        case "posts":
            main_content = <FeedClient INFO={INFO} />;
            break;
    }

    //
    const hasAccess = user.is_public || user.is_freind;

    //  
    const imageURL = user.profile_image;
    const fullImageURL = imageURL ? `${BASE}/pics/${imageURL}` : "";
    return (
        <div className={styles.FATHER}>

            <div className={styles.pageWrapper}>

                {/* HEADER */}
                <header className={styles.header}>
                    <div className={styles.coverPhoto}>
                        {!user.is_freind && <button className={styles.infoText}>follow </button>}
                    </div>

                    <div className={styles.profileInfo}>
                        {
                            fullImageURL ?
                                <img className={styles.profileAvatar} src={fullImageURL} />
                            :
                            <User className={styles.profileAvatar} />
                        }
                        <div className={styles.nameandprivacybuttoncontainer}>
                            <div className={styles.flnname}>
                                <h1 className={styles.profileName}>{user.firstname + " " + user.lastname}</h1>
                                <h5 className={styles.profileNickname}>{user.nickname}</h5>
                            </div>
                            <span className={styles.s}></span>

                        </div>
                        <p className={styles.profileBio}>{user.aboutme}</p>
                    </div>
                </header>

                {/* NAVIGATION */}
                {
                    hasAccess && (
                        <nav className={styles.NAV}>
                            <ul className={styles.navbar}>
                                <button className={`${styles.navBtn} ${section === "posts" ? styles.navBtnActive : ""}`} onClick={
                                    () => {
                                        setSection("posts")
                                        router.push(`?section=posts`)
                                    }}>Posts</button>
                                <button className={`${styles.navBtn} ${section === "about" ? styles.navBtnActive : ""}`} onClick={
                                    () => {
                                        setSection("about")
                                        router.push(`?section=about`)
                                    }}>About</button>
                            </ul>
                        </nav>
                    )
                }
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