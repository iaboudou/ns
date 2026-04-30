
import styles from "./about.module.css"

export function about(user) {
    return <div className={styles.aboutContainer}>
        <h2 className={styles.cardTitle}>Personal Information</h2>

        <div className={styles.infoRow}>
            <span className={styles.infoLabel}>First Name</span>
            <span className={styles.infoValue}>{user.firstname}</span>
        </div>
        <div className={styles.infoRow}>
            <span className={styles.infoLabel}>Last Name</span>
            <span className={styles.infoValue}>{user.lastname}</span>
        </div>
        <div className={styles.infoRow}>
            <span className={styles.infoLabel}>Nickname</span>
            <span className={styles.infoValue}>{user.nickname}</span>
        </div>
        <div className={styles.infoRow}>
            <span className={styles.infoLabel}>Email</span>
            <span className={styles.infoValue}>{user.email}</span>
        </div>
        <div className={styles.infoRow}>
            <span className={styles.infoLabel}>Gender</span>
            <span className={styles.infoValue}>{user.gender}</span>
        </div>
        <div className={styles.infoRow}>
            <span className={styles.infoLabel}>Birthday</span>
            <span className={styles.infoValue}>{user.birthday}</span>
        </div>
    </div>
}