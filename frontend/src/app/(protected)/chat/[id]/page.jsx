"use client";

import { useParams } from "next/navigation";
import styles from "./chat.module.css";
import { useState } from "react";
import { useWebSocket } from "@/lib/UseWebsocket";

export default function Page() {
  const params = useParams();
  const [inputText, setInputText] = useState("");

  const { messages, sendMessage } = useWebSocket();

  return (
    <div className={styles.chatContainer}>
      <div className={styles.chatHeader}>chat with User ID: {params.id}</div>

      <div className={styles.messagesArea}>
        {messages.map((msg) => (
          <div
            key={msg.id}
            className={`${styles.messageRow} ${msg.sender_id !== params.id ? styles.sentRow : styles.receivedRow}`}
          >
            <div
              className={`${styles.messageBubble} ${msg.sender_id !== params.id ? styles.sentBubble : styles.receivedBubble}`}
            >
              {msg.content}
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
        <button
          className={styles.sendButton}
          onClick={(e) => {
            e.preventDefault();
            const content = inputText.trim();

            if (content === "") return;

            sendMessage({
              type: "send",
              payload: {
                type: "chat",
                receiver_Id: params.id,
                content: content,
              },
            });

            setInputText("");
          }}
        >
          Send
        </button>
      </div>
    </div>
  );
}
