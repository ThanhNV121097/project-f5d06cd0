"use client";

import { useEffect, useMemo, useState, type FormEvent } from "react";
import styles from "./HelloPage.module.css";
import {
  createGreeting,
  getGreetings,
  getHello,
  type Greeting,
} from "../lib/hello-api";

type Status = "loading" | "ready" | "error";

const NAME_LIMIT = 80;
const MESSAGE_LIMIT = 240;

export default function HelloPage() {
  const [hello, setHello] = useState("Loading hello message...");
  const [status, setStatus] = useState<Status>("loading");
  const [error, setError] = useState("");
  const [name, setName] = useState("");
  const [message, setMessage] = useState("");
  const [greetings, setGreetings] = useState<Greeting[]>([]);
  const [saving, setSaving] = useState(false);
  const [helloName, setHelloName] = useState("");

  const statusLabel = useMemo(() => {
    if (status === "error") return "API unreachable";
    if (status === "loading") return "Connecting";
    return "Connected";
  }, [status]);

  async function load(nameOverride = helloName) {
    try {
      const [helloResponse, listResponse] = await Promise.all([
        getHello(nameOverride),
        getGreetings(),
      ]);
      setHello(helloResponse.message);
      setGreetings(listResponse);
      setStatus("ready");
      setError("");
    } catch {
      setStatus("error");
      setError("API unreachable. Check backend service and try again.");
    }
  }

  useEffect(() => {
    void load();
  }, []);

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const trimmedName = name.trim();
    const trimmedMessage = message.trim();

    if (!trimmedName || !trimmedMessage || trimmedName.length > NAME_LIMIT || trimmedMessage.length > MESSAGE_LIMIT) {
      setError("Name and message must be filled and kept short.");
      return;
    }

    setSaving(true);
    setError("");
    try {
      await createGreeting({ name: trimmedName, message: trimmedMessage });
      await load();
      setName("");
      setMessage("");
    } catch {
      setError("API unreachable. Check backend service and try again.");
      setStatus("error");
    } finally {
      setSaving(false);
    }
  }

  return (
    <main className={styles.page}>
      <section className={styles.shell}>
        <header className={styles.header}>
          <div>
            <p className={styles.kicker}>Hello World demo</p>
            <h1 className={styles.title}>Live hello, saved greetings, real stack proof.</h1>
            <p className={styles.lead}>
              Page reads API hello text, saves greeting, and refreshes stored greetings from backend data.
            </p>
          </div>
          <span className={status === "error" ? styles.pillError : styles.pill} aria-live="polite">
            <span className={status === "error" ? styles.dotError : styles.dot} aria-hidden="true" />
            {statusLabel}
          </span>
        </header>

        <div className={styles.grid}>
          <section className={styles.card} aria-labelledby="hello-heading">
            <h2 id="hello-heading" className={styles.sectionTitle}>Live hello</h2>
            <p className={styles.sectionLead}>Loaded from GET /api/hello.</p>
            <label className={styles.field} htmlFor="hello-name">
              <span>Name for GET /api/hello</span>
              <input id="hello-name" value={helloName} maxLength={NAME_LIMIT} onChange={(event) => setHelloName(event.target.value)} placeholder="Ada" />
            </label>
            <button className={styles.button} type="button" onClick={() => void load(helloName)}>
              Refresh hello
            </button>
            <div className={styles.helloBox}>
              <p className={styles.helloText}>{hello}</p>
            </div>
          </section>

          <section className={styles.card} aria-labelledby="form-heading">
            <h2 id="form-heading" className={styles.sectionTitle}>Add greeting</h2>
            <p className={styles.sectionLead}>Short, non-empty name and message.</p>
            {error ? <div className={styles.errorBox}>{error}</div> : null}
            <form className={styles.form} onSubmit={handleSubmit}>
              <label className={styles.field} htmlFor="name">
                <span>Name</span>
                <input id="name" value={name} maxLength={NAME_LIMIT} onChange={(event) => setName(event.target.value)} placeholder="Ada" />
              </label>
              <label className={styles.field} htmlFor="message">
                <span>Message</span>
                <textarea id="message" value={message} maxLength={MESSAGE_LIMIT} onChange={(event) => setMessage(event.target.value)} placeholder="Hello from browser" />
              </label>
              <button className={styles.button} type="submit" disabled={saving}>{saving ? "Saving..." : "Save greeting"}</button>
            </form>
          </section>
        </div>

        <section className={styles.cardWide} aria-labelledby="list-heading">
          <h2 id="list-heading" className={styles.sectionTitle}>Stored greetings</h2>
          <p className={styles.sectionLead}>Newest first from GET /api/greetings.</p>
          {greetings.length ? (
            <ul className={styles.list}>
              {greetings.map((greeting) => (
                <li key={greeting.id} className={styles.item}>
                  <div className={styles.itemTop}>
                    <strong>{greeting.name}</strong>
                    <time dateTime={greeting.created_at}>{greeting.created_at}</time>
                  </div>
                  <p>{greeting.message}</p>
                </li>
              ))}
            </ul>
          ) : (
            <div className={styles.emptyState}>No greetings yet. Submit first one above.</div>
          )}
        </section>
      </section>
    </main>
  );
}
