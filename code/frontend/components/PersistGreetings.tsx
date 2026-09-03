"use client";

import { FormEvent, useMemo, useState } from "react";
import styles from "./PersistGreetings.module.css";
import { helloMessage, listGreetings, saveGreeting, type Greeting } from "../lib/mock/persist-greetings";

const nameLimit = 80;
const messageLimit = 240;

export default function PersistGreetings() {
  const [greetings, setGreetings] = useState<Greeting[]>(() => listGreetings());
  const [name, setName] = useState("");
  const [message, setMessage] = useState("");
  const [error, setError] = useState("");
  const [apiDown, setApiDown] = useState(false);
  const hello = useMemo(() => (apiDown ? "Hello, World!" : helloMessage), [apiDown]);

  function onSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const trimmedName = name.trim();
    const trimmedMessage = message.trim();

    if (!trimmedName || !trimmedMessage) {
      setError("Name and message must be filled and kept short.");
      return;
    }

    if (trimmedName.length > nameLimit || trimmedMessage.length > messageLimit) {
      setError("Name or message is too long.");
      return;
    }

    setError("");
    const saved = saveGreeting({ name: trimmedName, message: trimmedMessage });
    setGreetings((current) => [saved, ...current]);
    setName("");
    setMessage("");
  }

  return (
    <main className={styles.page}>
      <section className={styles.hero} aria-labelledby="hello-title">
        <div>
          <p className={styles.kicker}>Hello World demo</p>
          <h1 id="hello-title" className={styles.title}>
            Say hello, save greetings, prove DB persistence.
          </h1>
          <p className={styles.lead}>
            Live hello text, greeting form, and stored greetings list.
          </p>
        </div>
        <div className={styles.statusRow}>
          <span className={`${styles.statusPill} ${apiDown ? styles.statusError : styles.statusSuccess}`}>
            {apiDown ? "API unreachable" : "Health: ok"}
          </span>
        </div>
      </section>

      <section className={styles.grid}>
        <article className={styles.panel}>
          <div className={styles.panelHead}>
            <div>
              <h2 className={styles.sectionTitle}>Live hello</h2>
              <p className={styles.sectionLead}>Loaded from GET /api/hello</p>
            </div>
          </div>
          <div className={styles.helloBox}>
            <strong>{hello}</strong>
            <button className={styles.primaryButton} type="button" onClick={() => setApiDown((value) => !value)}>
              Refresh
            </button>
          </div>
          {apiDown ? <div className={styles.errorBox}>API unreachable. Check backend service and try again.</div> : null}
        </article>

        <article className={styles.panel}>
          <div className={styles.panelHead}>
            <div>
              <h2 className={styles.sectionTitle}>Add greeting</h2>
              <p className={styles.sectionLead}>Validates short, non-empty name and message</p>
            </div>
          </div>
          {error ? <div className={styles.errorBox}>{error}</div> : null}
          <form className={styles.form} onSubmit={onSubmit}>
            <label className={styles.field} htmlFor="name">
              Name
              <input id="name" value={name} maxLength={nameLimit} onChange={(event) => setName(event.target.value)} placeholder="Ada" />
            </label>
            <label className={styles.field} htmlFor="message">
              Message
              <textarea id="message" value={message} maxLength={messageLimit} onChange={(event) => setMessage(event.target.value)} placeholder="Hello from the browser" />
            </label>
            <button className={styles.primaryButton} type="submit">Save greeting</button>
          </form>
        </article>
      </section>

      <section className={styles.panel} aria-labelledby="greetings-title">
        <div className={styles.panelHead}>
          <div>
            <h2 id="greetings-title" className={styles.sectionTitle}>Stored greetings</h2>
            <p className={styles.sectionLead}>Newest first from GET /api/greetings</p>
          </div>
        </div>
        {greetings.length ? (
          <div className={styles.list}>
            {greetings.map((greeting) => (
              <article key={greeting.id} className={styles.item}>
                <div className={styles.itemMeta}>
                  <strong>{greeting.name}</strong>
                  <span>{formatDate(greeting.created_at)}</span>
                </div>
                <p>{greeting.message}</p>
              </article>
            ))}
          </div>
        ) : (
          <div className={styles.emptyState}>No greetings yet. Submit first one above.</div>
        )}
      </section>
    </main>
  );
}

function formatDate(value: string) {
  return new Intl.DateTimeFormat("en-GB", {
    year: "numeric",
    month: "2-digit",
    day: "2-digit",
    hour: "2-digit",
    minute: "2-digit",
    hour12: false,
    timeZone: "UTC",
  }).format(new Date(value)).replace(/\//g, "-").replace(",", "");
}
