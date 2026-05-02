"use client";

import { useParams } from "next/navigation";
import styles from "./chat.module.css";
import { useState } from "react";

export default function Page() {
  const params = useParams();
  const [inputText, setInputText] = useState("");

  // example messages
  const messages = [
    { id: 1, text: "Hello! How are you?", isSentByMe: false },
    { id: 2, text: "I'm doing well, thanks! How about you?", isSentByMe: true },
    { id: 3, text: "Great! Just working on this chat UI.", isSentByMe: false },
  ];

  return (
    <div className={styles.chatContainer}>
      <div className={styles.chatHeader}>
        chat with User ID: {params.id}
      </div>

      <div className={styles.messagesArea}>
        {messages.map((msg) => (
          <div
            key={msg.id}
            className={`${styles.messageRow} ${msg.isSentByMe ? styles.sentRow : styles.receivedRow}`}
          >
            <div className={`${styles.messageBubble} ${msg.isSentByMe ? styles.sentBubble : styles.receivedBubble}`}>
              {msg.text}
            </div>
          </div>
        ))}
      </div>

      <div className={styles.inputArea}>
        <input
          type="text"
          className={styles.inputField}
          placeholder="Type a message..."
          value={inputText}
          onChange={(e) => setInputText(e.target.value)}
        />
        <button className={styles.sendButton}>Send</button>
      </div>
    </div>
  );
}