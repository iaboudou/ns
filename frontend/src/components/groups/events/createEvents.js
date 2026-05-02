'use client';

import { CreateEvent } from '@/_lib/group';
import styles from '@/components/groups/styles/groups.module.css';
import { useState } from 'react';

export default function CreateEvents({ groupId, setEvents }) {
  const [eventTitle, setEventTitle] = useState('');
  const [eventDescription, setEventDescription] = useState('');
  const [eventDateTime, setEventDateTime] = useState('');
  const [error, setError] = useState('');

  const handleCreateEvent = async (e) => {
    e.preventDefault();

    if (eventTitle.trim() === '' || eventDescription.trim() === '' || eventDateTime.trim() === '') {
      setError('please fill all the fields');
      return;
    }

    const date = new Date(eventDateTime);

    if (date.getTime() <= Date.now() + 5 * 60 * 1000) {
      setError('an event must be at least in 5 minutes');
      return;
    }

    const isoDate = date.toISOString();

    CreateEvent(eventTitle.trim(), eventDescription.trim(), isoDate, groupId)
      .then((newEvents) => {
        setEventTitle('');
        setEventDescription('');
        setEventDateTime('');
        setError('');
        setEvents((prev) => [newEvents, ...prev]);
      })
      .catch((err) => {
        if (err.message === 'Something wrong happened. Please try later') alert(err.message);
        else setError(err.message);
      });
  };

  return (
    <form className={styles.form} onSubmit={handleCreateEvent}>
      <div className={styles.formGroup}>
        {error && <p className={styles.groupError}>{error}</p>}
        <input minLength={3} maxLength={80} className={styles.input} placeholder="Title" value={eventTitle} onChange={(e) => setEventTitle(e.target.value)} />
      </div>

      <div className={styles.formGroup}>
        <input minLength={10} maxLength={500} className={styles.input} placeholder="Description" value={eventDescription} onChange={(e) => setEventDescription(e.target.value)} />
      </div>

      <div className={styles.formGroup}>
        <input className={styles.input} type="datetime-local" value={eventDateTime} onChange={(e) => setEventDateTime(e.target.value)} />
      </div>

      <button className={styles.primaryButton}>Create</button>
    </form>
  );
}
