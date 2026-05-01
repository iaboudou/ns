import styles from './UsersList.module.css';

export default function UsersList() {
    const users = [
        { id: 1, name: "Alex" },
        { id: 2, name: "Sarah" },
        { id: 3, name: "John" },
        { id: 4, name: "John" },
        { id: 5, name: "John" },
    ];

    return (
        <div className={styles.usersList}>
            <h3>Users</h3>
            {users.map((user) => (
                <a key={user.id} href={`/chat/${user.id}`}>
                    <div key={user.id} className={styles.userItem}>
                        {user.name}
                    </div>
                </a>
            ))}
        </div>
    );
}