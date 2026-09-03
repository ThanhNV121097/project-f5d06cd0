"use client";

import { FormEvent, useEffect, useState } from "react";

type Greeting = {
  id: string;
  name: string;
  message: string;
  created_at: string;
};

type ApiGreetingList = {
  greetings: Greeting[];
};

const apiBase = process.env.NEXT_PUBLIC_API_URL ?? "/api";

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const res = await fetch(`${apiBase}${path}`, {
    ...init,
    headers: {
      "Content-Type": "application/json",
      ...(init?.headers ?? {})
    }
  });
  if (!res.ok) throw new Error("API unreachable");
  return res.json() as Promise<T>;
}

export default function Home() {
  const [hello, setHello] = useState("Hello, World!");
  const [helloName, setHelloName] = useState("Ada");
  const [greetings, setGreetings] = useState<Greeting[]>([]);
  const [error, setError] = useState<string | null>(null);
  const [name, setName] = useState("");
  const [message, setMessage] = useState("");

  async function loadHello(nameParam = "") {
    const query = nameParam ? `?name=${encodeURIComponent(nameParam)}` : "";
    const helloRes = await request<{ message: string }>(`/v1/hello${query}`);
    setHello(helloRes.message);
  }

  async function loadGreetings() {
    const listRes = await request<ApiGreetingList>("/v1/greetings");
    setGreetings(listRes.greetings);
  }

  useEffect(() => {
    (async () => {
      try {
        await Promise.all([loadHello(), loadGreetings()]);
        setError(null);
      } catch {
        setError("API unreachable. Check backend is running, then reload.");
      }
    })();
  }, []);

  async function refreshHello(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    try {
      await loadHello(helloName);
      setError(null);
    } catch {
      setError("API unreachable. Check backend is running, then reload.");
    }
  }

  async function submitForm(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    try {
      const saved = await request<Greeting>("/v1/greetings", {
        method: "POST",
        body: JSON.stringify({ name, message })
      });
      await Promise.all([loadHello(), loadGreetings()]);
      setGreetings((current) => [saved, ...current].slice(0, current.length + 1));
      setName("");
      setMessage("");
      setError(null);
    } catch {
      setError("API unreachable. Greeting kept in form. Try again later.");
    }
  }

  return (
    <main className="app-shell">
      <section className="mx-auto flex min-h-[calc(100vh-4rem)] w-full max-w-3xl flex-col gap-6 rounded-[var(--radius-xl)] border border-[var(--color-border)] bg-[var(--color-surface)] p-6 shadow-[var(--shadow-lg)]">
        <div>
          <p className="text-sm font-medium uppercase tracking-[0.18em] text-[var(--color-accent-secondary)]">Hello World</p>
          <h1 className="mt-2 text-3xl">{hello}</h1>
        </div>

        <form onSubmit={refreshHello} className="flex items-end gap-3 rounded-[var(--radius-lg)] border border-[var(--color-border)] p-4" aria-label="Refresh hello panel">
          <div className="grid flex-1 gap-2">
            <label htmlFor="hello-name">Hello name</label>
            <input id="hello-name" value={helloName} onChange={(e) => setHelloName(e.target.value)} className="rounded-[var(--radius-sm)] border border-[var(--color-border)] px-4 py-3" maxLength={80} />
          </div>
          <button className="w-fit rounded-[var(--radius-full)] bg-[var(--color-primary)] px-5 py-3 font-medium text-[var(--color-primary-text)]" type="submit">Refresh hello</button>
        </form>

        {error ? <div className="rounded-[var(--radius-md)] border border-[var(--color-danger)] bg-red-50 p-4 text-[var(--color-danger)]">{error}</div> : null}

        <form onSubmit={submitForm} className="grid gap-4 rounded-[var(--radius-lg)] border border-[var(--color-border)] p-5" aria-label="Add greeting form">
          <div className="grid gap-2">
            <label htmlFor="name">Name</label>
            <input id="name" value={name} onChange={(e) => setName(e.target.value)} className="rounded-[var(--radius-sm)] border border-[var(--color-border)] px-4 py-3" maxLength={80} />
          </div>
          <div className="grid gap-2">
            <label htmlFor="message">Message</label>
            <textarea id="message" value={message} onChange={(e) => setMessage(e.target.value)} className="min-h-28 rounded-[var(--radius-sm)] border border-[var(--color-border)] px-4 py-3" maxLength={240} />
          </div>
          <button className="w-fit rounded-[var(--radius-full)] bg-[var(--color-primary)] px-5 py-3 font-medium text-[var(--color-primary-text)]" type="submit">Save greeting</button>
        </form>

        <section className="grid gap-3">
          <h2 className="text-xl">Stored greetings</h2>
          <div className="grid gap-3">
            {greetings.length === 0 ? <p className="text-[var(--color-text-muted)]">No greetings yet.</p> : greetings.map((g) => (
              <article key={g.id} className="rounded-[var(--radius-lg)] border border-[var(--color-border)] p-4">
                <p className="font-medium">{g.name}</p>
                <p>{g.message}</p>
                <p className="text-sm text-[var(--color-text-muted)]">{new Date(g.created_at).toLocaleString()}</p>
              </article>
            ))}
          </div>
        </section>
      </section>
    </main>
  );
}
