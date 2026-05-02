'use client';

import styles from '@/components/groups/styles/groups.module.css';

import { CreateGroup } from '@/_lib/group';
import { useState } from 'react';

export default function HandleCreateGroup({ setGroups }) {
  const [groupTitle, setGroupTitle] = useState('');
  const [groupDescription, setGroupDescription] = useState('');
  const [error, setError] = useState('');

  return (
    <form
      className={styles.form}
      onSubmit={(e) => {
        e.preventDefault();

        if (groupTitle.trim() === '' || groupDescription.trim() === '') {
          setError('please fill all the fields');
          return;
        }

        CreateGroup(groupTitle, groupDescription)
          .then((newGroup) => {
            setGroupTitle('');
            setGroupDescription('');
            setError('');
            setGroups((prev) => [newGroup, ...prev]);
          })
          .catch((err) => {
            if (err.message === 'Something wrong happened. Please try later') alert(err.message);
            else setError(err.message);
          });
      }}
    >
      <div className={styles.formGroup}>
        {error !== '' && <p className={styles.groupError}>{error}</p>}
        <input className={styles.input} placeholder="Title" maxLength={60} minLength={3} value={groupTitle} onChange={(e) => setGroupTitle(e.target.value)} />
      </div>

      <div className={styles.formGroup}>
        <input className={styles.input} placeholder="Description" maxLength={300} minLength={20} value={groupDescription} onChange={(e) => setGroupDescription(e.target.value)} />
      </div>

      <button className={styles.primaryButton}>Create</button>
    </form>
  );
}
