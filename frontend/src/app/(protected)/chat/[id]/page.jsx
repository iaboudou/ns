"use client";

import { useParams } from "next/navigation";
import styles from "./chat.module.css";
import { useState, useEffect, useRef } from "react";
import { useWebSocket } from "@/lib/UseWebsocket";

export default function Page() {
  const params = useParams();
  const [inputText, setInputText] = useState("");
  const [showEmojis, setShowEmojis] = useState(false);
  const [chatUser, setChatUser] = useState(null);
  const scrollRef = useRef(null);

  const { messages, sendMessage, port, hasMoreMap } = useWebSocket();
  const hasMore = hasMoreMap[params.id] !== false; // default true

  // fetch user details for the header
  useEffect(() => {
    const fetchUser = async () => {
      try {
        const resp = await fetch(`http://localhost:4001/api/getUsers?search=${params.id}`, {
          credentials: "include"
        });
        const res = await resp.json();
        if (res.code === 200 && res.data.length > 0) {
          // extra filter
          const user = res.data.find(u => u.id === params.id);
          if (user) setChatUser(user);
        }
      } catch (err) {
        console.error("Failed to fetch chat user", err);
      }
    };
    if (params.id) fetchUser();
  }, [params.id]);

  useEffect(() => {
    if (port && params.id) {
      // mark messages as read in DB
      port.postMessage({
        type: "send",
        payload: {
          type: "mark_read",
          receiver_Id: params.id,
        },
      });

      // load history
      port.postMessage({
        type: "send",
        payload: {
          type: "load_history",
          receiver_Id: params.id,
          last_read_time: 0,
        },
      });
    }
  }, [port, params.id]);

  // filter messages for this specific conversation
  const conversationMessages = messages.filter(
    (msg) => msg.sender_id === params.id || msg.receiver_id === params.id
  ).sort((a, b) => a.created_at - b.created_at);

  // scroll to bottom when messages change
  useEffect(() => {
    if (scrollRef.current) {
      // if the last message is new, we definitely scroll to bottom
      // if it's history, we might want to be careful, but for now let's scroll
      const lastMsg = conversationMessages[conversationMessages.length - 1];
      if (lastMsg?.isNew || conversationMessages.length <= 10) {
        scrollRef.current.scrollTop = scrollRef.current.scrollHeight;
      }
    }
  }, [conversationMessages]);

  const handleLoadMore = () => {
    if (conversationMessages.length === 0) return;
    const oldestTime = conversationMessages[0].created_at;
    port.postMessage({
      type: "send",
      payload: {
        type: "load_history",
        receiver_Id: params.id,
        last_read_time: oldestTime,
      },
    });
  };

  const commonEmojis = ["👍", "😀", "😂", "🥰", "😎", "🤔", "😅", "🔥", "❤️", "🙏", "✨", "🎉"];

  return (
    <div className={styles.chatContainer}>
      <div className={styles.chatHeader}>
        {chatUser ? (
          <>
            <span className={styles.headerName}>{chatUser.firstname} {chatUser.lastname}</span>
            {chatUser.nickname && <span className={styles.headerNick}> @{chatUser.nickname}</span>}
          </>
        ) : (
          `Loading chat...`
        )}
      </div>

      <div className={styles.messagesArea} ref={scrollRef}>
        {hasMore && conversationMessages.length >= 10 && (
          <button className={styles.loadMore} onClick={handleLoadMore}>
            Load older messages
          </button>
        )}
        {conversationMessages.map((msg) => (
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
        <div className={styles.emojiWrapper}>
          <button 
            className={styles.emojiToggle}
            onClick={() => setShowEmojis(!showEmojis)}
          >
            😀
          </button>
          {showEmojis && (
            <div className={styles.emojiPicker}>
              {commonEmojis.map(emoji => (
                <span 
                  key={emoji} 
                  className={styles.emojiItem}
                  onClick={() => {
                    setInputText(prev => prev + emoji);
                    setShowEmojis(false);
                  }}
                >
                  {emoji}
                </span>
              ))}
            </div>
          )}
        </div>

        <input
          type="text"
          className={styles.inputField}
          placeholder="Type a message..."
          value={inputText}
          onChange={(e) => setInputText(e.target.value)}
          onKeyDown={(e) => {
            if (e.key === "Enter") {
               const content = inputText.trim();
               if (content === "") return;
               sendMessage({
                 type: "send",
                 payload: { type: "chat", receiver_Id: params.id, content: content },
               });
               setInputText("");
            }
          }}
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
