export default function ChatWindow({ isSelected, userId }) {
    if (isSelected) {
        return (
            <div>
                <h1>User id: {userId}</h1>
            </div>
        );
    } else {
        return (
            <div>
                <h1>you did not select the chat yet</h1>
            </div>
        );
    }
}