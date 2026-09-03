export type Greeting = {
  id: string;
  name: string;
  message: string;
  created_at: string;
};

type GreetingInput = {
  name: string;
  message: string;
};

type GreetingsResponse = {
  greetings: Greeting[];
  next_cursor: string | null;
  has_more: boolean;
};

const apiBase = process.env.NEXT_PUBLIC_API_URL ?? "/api";

async function readJSON<T>(response: Response): Promise<T> {
  if (!response.ok) {
    throw new Error("api unreachable");
  }
  return response.json() as Promise<T>;
}

export async function fetchHello(name = "") {
  const url = new URL(`${apiBase}/v1/hello`, window.location.origin);
  if (name) url.searchParams.set("name", name);
  return readJSON<{ message: string }>(await fetch(url.toString()));
}

export async function fetchGreetings(limit = 20) {
  const url = new URL(`${apiBase}/v1/greetings`, window.location.origin);
  url.searchParams.set("limit", String(limit));
  const response = await readJSON<GreetingsResponse>(await fetch(url.toString()));
  return { greetings: response.greetings };
}

export async function saveGreeting(input: GreetingInput) {
  return readJSON<Greeting>(
    await fetch(`${apiBase}/v1/greetings`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(input),
    }),
  );
}
