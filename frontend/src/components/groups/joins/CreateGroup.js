"use client";

import styles from "@/components/groups/styles/groups.module.css";

import { CreateGroup } from "@/_lib/group";
import { useState } from "react";

export default function HandleCreateGroup({
  setGroups,
  setShowCreateGroupForm,
}) {
  const [groupTitle, setGroupTitle] = useState("");
  const [groupDescription, setGroupDescription] = useState("");
  const [groupImage, setGroupImage] = useState(null);

  const [error, setError] = useState("");

  return (
    <form
      className={styles.form}
      onSubmit={(e) => {
        e.preventDefault();

        if (groupTitle.trim() === "" || groupDescription.trim() === "") {
          setError("please fill all the fields");
          return;
        }

        const formData = new FormData();

        formData.append("title", groupTitle);
        formData.append("description", groupDescription);

        if (groupImage) {
          formData.append("image", groupImage);
        }

        CreateGroup(formData)
          .then((newGroup) => {
            setGroupTitle("");
            setGroupDescription("");
            setGroupImage(null);
            setError("");

            setGroups((prev) => [newGroup, ...prev]);
            setShowCreateGroupForm((v) => !v);
          })
          .catch((err) => {
            if (err.message === "Something wrong happened. Please try later") {
              alert(err.message);
            } else {
              setError(err.message);
            }
          });
      }}
    >
      <div className={styles.formGroup}>
        {error !== "" && <p className={styles.groupError}>{error}</p>}

        <input
          className={styles.input}
          placeholder="Title"
          maxLength={60}
          minLength={3}
          value={groupTitle}
          onChange={(e) => setGroupTitle(e.target.value)}
        />
      </div>

      <div className={styles.formGroup}>
        <input
          className={styles.input}
          placeholder="Description"
          maxLength={300}
          minLength={20}
          value={groupDescription}
          onChange={(e) => setGroupDescription(e.target.value)}
        />
      </div>

      <div>
        <input
          type="file"
          accept="image/*"
          onChange={(e) => setGroupImage(e.target.files?.[0] || null)}
        />
      </div>

      <button className={styles.primaryButton}>Create</button>
    </form>
  );
}
