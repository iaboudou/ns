"use client";
import { useState } from "react";
import styles from "@/app/(protected)/page.module.css";
import FeedClient from "@/app/(protected)/feed/feedclient";
import FollowSuggestions from "@/app/(protected)/feed/followsuggestions/followsuggestions";
import { UserPlus } from "lucide-react";

export default function Home() {
  const [show, setShow] = useState(false);

  return (
    <div className={styles.WRAPPER}>
      <div className={styles.MIDDLE}>
        
        <button 
          className={styles.circleBtn} 
          onClick={() => setShow(!show)} 
          title="Suggestions"
        >
          <UserPlus size={22} />
        </button>

        {show && (
          <>
            <div className={styles.overlay} onClick={() => setShow(false)} />
            <div className={styles.mobileSuggestions}>
              <FollowSuggestions />
            </div>
          </>
        )}

        <FeedClient />
      </div>
      
      <div className={styles.RIGHT}>
        <FollowSuggestions />
      </div>
    </div>
  );
}
