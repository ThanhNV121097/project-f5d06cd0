"use client";

import { FormEvent, useEffect, useState } from "react";
import styles from "./PersistGreetings.module.css";
import {
  fetchHelloMessage,
  fetchGreetings,
  saveGreeting,
  type Greeting,
} from "../lib/mock/persist-greetings";

const nameLimit = 80;
const messageLimit = 240;

export default function PersistGreetings() {
  const [greetings, setGreetings] = useState<Greeting[]>([]);
  const [name, setName] = useState("");
  const [message, setMessage] = useState("");
  const [formError, setFormError] = useState("");
  const [apiError, setApiError] = useState("");
  const [hello, setHello] = useState("Hello, World!");

  useEffect(() => {
    void loadData();
  }, []);

  async function loadData() {
    try {
      const [helloResponse, greetingResponse] = await Promise.all([fetchHelloMessage(), fetchGreetings()]);
      setHello(helloResponse.message);
      setGreetings(greetingResponse.greetings);
      setApiError("");
    } catch {
      setApiError("API unreachable. Check backend service and try again.");
    }
  }

  async function onSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const trimmedName = name.trim();
    const trimmedMessage = message.trim();

    if (!trimmedName || !trimmedMessage) {
      setFormError("Name and message must be filled and kept short.");
      return;
    }

    if (trimmedName.length > nameLimit || trimmedMessage.length > messageLimit) {
      setFormError("Name or message is too long.");
      return;
    }

    try {
      await saveGreeting({ name: trimmedName, message: trimmedMessage });
      await loadData();
      setName("");
      setMessage("");
      setFormError("");
      setApiError("");
    } catch {
      setApiError("API unreachable. Check backend service and try again.");
    }
  }

  return (
    <main className={styles.page}>
      <section className={styles.hero} aria-labelledby="hello-title">
        <div>
          <p className={styles.kicker}>Hello World demo</p>
          <h1 id="hello-title" className={styles.title}>
            Say hello, save greetings, prove DB persistence.
          </h1>
          <p className={styles.lead}>Live hello text, greeting form, and stored greetings list.</p>
        </div>
        <div className={styles.statusRow}>
          <span className={`${styles.statusPill} ${apiError ? styles.statusError : styles.statusSuccess}`}>
            {apiError ? "API unreachable" : "Health: ok"}
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
            <button className={styles.primaryButton} type="button" onClick={() => void loadData()}>
              Refresh
            </button>
          </div>
          {apiError ? <div className={styles.errorBox}>{apiError}</div> : null}
        </article>

        <article className={styles.panel}>
          <div className={styles.panelHead}>
            <div>
              <h2 className={styles.sectionTitle}>Add greeting</h2>
              <p className={styles.sectionLead}>Validates short, non-empty name and message</p>
            </div>
          </div>
          {formError ? <div className={styles.errorBox}>{formError}</div> : null}
          <form className={styles.form} onSubmit={onSubmit}>
            <label className={styles.field} htmlFor="name">
              Name
              <input id="name" value={name} maxLength={nameLimit} onChange={(event) => setName(event.target.value)} placeholder="Ada" />
            </label>
            <label className={styles.field} htmlFor="message">
              Message
              <textarea id="message" value={message} maxLength={messageLimit} onChange={(event) => setMessage(event.target.value)} placeholder="Hello from the browser" />
            </label>
            <button className={styles.primaryButton} type="submit">
              Save greeting
            </button>
          </form>
        </article>
      </section>

      <section className={styles.panel} aria-labelledby="greetings-title">
        <div className={styles.panelHead}>
          <div>
            <h2 id="greetings-title" className={styles.sectionTitle}>
              Stored greetings
            </h2>
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
  })
    .format(new Date(value))
    .replace(/\//g, "-")
    .replace(",", "");
}

