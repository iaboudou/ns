export default function UsersList() {
    const users = [
        { id: 1, name: "Alex" },
        { id: 2, name: "Sarah" },
        { id: 3, name: "John" },
        { id: 4, name: "John" },
        { id: 5, name: "John" },
    ];

    return (
        <div>
            <h3>Users</h3>
            {users.map((user) => (
                <a key={user.id} href={`/chat/${user.id}`}>
                    <div key={user.id} style={{ padding: "10px", cursor: "pointer" }}>
                        {user.name}
                    </div>
                </a>
            ))}
        </div>
    );
}