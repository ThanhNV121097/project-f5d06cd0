"use client";

import { useEffect, useState, type FormEvent } from "react";
import { createGreeting, getGreetings, type Greeting } from "../lib/hello-api";

const LIMITS = { name: 80, message: 240 };

export default function PersistGreetings() {
  const [name, setName] = useState("");
  const [message, setMessage] = useState("");
  const [greetings, setGreetings] = useState<Greeting[]>([]);
  const [error, setError] = useState("");

  async function load() {
    setGreetings(await getGreetings());
  }

  useEffect(() => {
    load().catch(() => setError("API unreachable. Check backend service and try again."));
  }, []);

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const trimmedName = name.trim();
    const trimmedMessage = message.trim();
    if (!trimmedName || !trimmedMessage || trimmedName.length > LIMITS.name || trimmedMessage.length > LIMITS.message) {
      setError("Name and message must be filled and kept short.");
      return;
    }
    try {
      await createGreeting({ name: trimmedName, message: trimmedMessage });
      await load();
      setName("");
      setMessage("");
      setError("");
    } catch {
      setError("API unreachable. Check backend service and try again.");
    }
  }

  return (
    <section>
      <h2>Persist greetings</h2>
      <p>Live data from GET /api/greetings.</p>
      {error ? <div>{error}</div> : null}
      <form onSubmit={handleSubmit}>
        <input value={name} onChange={(e) => setName(e.target.value)} placeholder="Ada" maxLength={LIMITS.name} />
        <textarea value={message} onChange={(e) => setMessage(e.target.value)} placeholder="Hello from browser" maxLength={LIMITS.message} />
        <button type="submit">Save greeting</button>
      </form>
      <ul>
        {greetings.map((greeting) => (
          <li key={greeting.id}>{greeting.name}: {greeting.message}</li>
        ))}
      </ul>
    </section>
  );
}
