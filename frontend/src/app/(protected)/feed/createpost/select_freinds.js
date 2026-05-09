'use client';

import { useEffect, useState } from "react";
import styles from "./createpost.module.css";
import { fetchFriendsUsers } from "./actions";
import { User } from 'lucide-react';

export function SelectFreinds({ STATE }) {
    const { state, setState } = STATE;
    const [usersfetched, setUsersFetched] = useState([])

    useEffect(() => {
        let active = true;
        const t = setTimeout(async () => {
            if (state.privacy !== "private") return;
            const data = await fetchFriendsUsers(state.search);
            if (!active) return;
            setUsersFetched(data);
        }, 300);

        return () => {
            active = false;
            clearTimeout(t);
        };
    }, [state.search, state.privacy, setState]);

    if (state.privacy !== "private") return null;

    return (
        <div className={styles.selectFriendInput}>
            <input
                className={styles.searchInput}
                placeholder="Select friends..."
                value={state.search || ""}
                onChange={(e) =>
                    setState((prev) => ({ ...prev, search: e.target.value }))
                }
                maxLength={30}
            />

            {/* selected users */}
            <div className={styles.selectedusers}>

                {(state.selectedUsers || []).length == 0
                    ? <div className={styles.NouserSelected}></div>
                    : state.selectedUsers.map((u) => {

                        return (
                            <div key={"selected" + u.id} className={styles.selecteduserscontainer}>
                                {u.profile_image ? (
                                    <img
                                        src={`/pics/${u.profile_image}`}
                                        alt=""
                                        className={styles.selecteduseravatar}
                                    />
                                ) : (
                                    <div className={styles.selecteduseravatar}>
                                        <User size={12} />
                                    </div>
                                )}
                                <span>
                                    {u.firstname} {u.lastname}
                                </span>
                                <button
                                    type="button"
                                    onClick={() => {
                                        setState(prev => ({
                                            ...prev,
                                            selectedUsers: prev.selectedUsers.filter(t => t.id !== u.id)
                                        }))
                                    }}>x</button>
                            </div>
                        );
                    })}

            </div>

            {/* selct users */}
            <div className={styles.usersList}>
                {(usersfetched || []).length == 0
                    ? <div className={styles.NouserSelected}>No users exists...</div>
                    : usersfetched.map((u) => {

                        return (
                            <button
                                type="button"
                                key={u.id}
                                className={`${styles.userItem} ${state.selectedUsers?.some(selected => selected.id === u.id) ? styles.activeUser : ""}`}
                                onClick={() => {
                                    setState(prev => {
                                        return {
                                            ...prev,
                                            selectedUsers: prev.selectedUsers.some(y => y.id === u.id) ? prev.selectedUsers : [...prev.selectedUsers, u]
                                        }
                                    })
                                }}
                            >
                                {u.profile_image ? (
                                    <img
                                        src={`/pics/${u.profile_image}`}
                                        alt=""
                                        className={styles.userAvatar}
                                    />
                                ) : (
                                    <div className={styles.userAvatar}>
                                        <User size={16} />
                                    </div>
                                )}
                                <span>
                                    {u.firstname} {u.lastname}
                                </span>
                            </button>
                        );
                    })}
            </div>
        </div>
    );
}
