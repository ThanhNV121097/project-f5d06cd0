"use client";

import { FormEvent, useEffect, useState } from "react";
import styles from "./HelloApiDemo.module.css";
import { createGreeting, getHello, getGreetings, type Greeting } from "../lib/hello-api";

type ApiStatus = "ready" | "error";

const LIMITS = { name: 80, message: 240 };

export default function HelloApiDemo() {
  const [apiError, setApiError] = useState<string | null>(null);
  const [status, setStatus] = useState<ApiStatus | null>(null);
  const [helloMessage, setHelloMessage] = useState("Hello, World!");
  const [greetings, setGreetings] = useState<Greeting[]>([]);
  const [saving, setSaving] = useState(false);
  const [formError, setFormError] = useState<string | null>(null);

  const empty = greetings.length === 0;

  async function load() {
    const [hello, stored] = await Promise.all([getHello(), getGreetings()]);
    setHelloMessage(hello.message);
    setGreetings(stored);
    setStatus("ready");
    setApiError(null);
  }

  useEffect(() => {
    let live = true;
    load().catch(() => {
      if (!live) return;
      setStatus("error");
      setApiError("API unreachable. Check backend service and try again.");
    });
    return () => {
      live = false;
    };
  }, []);

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const form = event.currentTarget;
    const data = new FormData(form);
    const name = String(data.get("name") ?? "").trim();
    const message = String(data.get("message") ?? "").trim();

    if (!name || !message || name.length > LIMITS.name || message.length > LIMITS.message) {
      setFormError("Name and message must be filled and kept short.");
      return;
    }

    setSaving(true);
    setFormError(null);
    try {
      await createGreeting({ name, message });
      await load();
      form.reset();
    } catch {
      setFormError("Save failed. Try again.");
    } finally {
      setSaving(false);
    }
  }

  return (
    <main className={styles.page}>
      <section className={styles.hero}>
        <div>
          <p className={styles.kicker}>Hello World demo</p>
          <h1 className={styles.title}>Say hello, save greetings, prove DB persistence.</h1>
          <p className={styles.lead}>Clean demo page pulls live hello text from API, submits greetings to backend, and refreshes stored greetings list after save.</p>
        </div>
        {status ? (
          <div className={styles.pillRow} aria-label="API status">
            <span className={`${styles.pill} ${status === "error" ? styles.pillError : styles.pillOk}`}>
              <span className={styles.dot} />{status === "error" ? "API unreachable" : "Health: ok"}
            </span>
          </div>
        ) : null}
      </section>

      <section className={styles.grid}>
        <article className={styles.panel}>
          <h2 className={styles.sectionTitle}>Live hello</h2>
          <p className={styles.sectionLead}>Loaded from GET /v1/hello</p>
          <div className={styles.helloBox}>
            <strong>{helloMessage}</strong>
          </div>
          {apiError ? <div className={styles.errorBox}>{apiError}</div> : null}
        </article>

        <article className={styles.panel}>
          <h2 className={styles.sectionTitle}>Add greeting</h2>
          <p className={styles.sectionLead}>Validates short, non-empty name and message</p>
          {formError ? <div className={styles.errorBox}>{formError}</div> : null}
          <form className={styles.form} onSubmit={handleSubmit}>
            <label className={styles.field} htmlFor="name">
              <span>Name</span>
              <input id="name" name="name" maxLength={LIMITS.name} placeholder="Ada" />
            </label>
            <label className={styles.field} htmlFor="message">
              <span>Message</span>
              <textarea id="message" name="message" maxLength={LIMITS.message} placeholder="Hello from the browser" />
            </label>
            <button className={styles.button} type="submit" disabled={saving}>{saving ? "Saving..." : "Save greeting"}</button>
          </form>
        </article>
      </section>

      <section className={styles.panel}>
        <h2 className={styles.sectionTitle}>Stored greetings</h2>
        <p className={styles.sectionLead}>Newest first from GET /api/greetings</p>
        {empty ? <div className={styles.empty}>No greetings yet. Submit first one above.</div> : (
          <div className={styles.list}>
            {greetings.map((greeting) => (
              <article key={greeting.id} className={styles.item}>
                <div className={styles.meta}><strong>{greeting.name}</strong><span>{greeting.created_at}</span></div>
                <p>{greeting.message}</p>
              </article>
            ))}
          </div>
        )}
      </section>
    </main>
  );
}

