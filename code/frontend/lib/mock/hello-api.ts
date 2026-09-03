export type Greeting = {
  id: string;
  name: string;
  message: string;
  created_at: string;
};

const apiBase = process.env.NEXT_PUBLIC_API_URL ?? "/api";

type GreetingsEnvelope = {
  greetings: Greeting[];
};

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const response = await fetch(`${apiBase}/v1${path}`, {
    ...init,
    headers: {
      "Content-Type": "application/json",
      ...(init?.headers ?? {})
    }
  });

  if (!response.ok) {
    throw new Error("API request failed");
  }

  return response.json() as Promise<T>;
}

export function getHello(name = "") {
  const query = name.trim() ? `?name=${encodeURIComponent(name.trim())}` : "";
  return request<{ message: string }>(`/hello${query}`);
}

export function getGreetings() {
  return request<GreetingsEnvelope>("/greetings").then((data) => data.greetings);
}

export function createGreeting(input: Pick<Greeting, "name" | "message">) {
  return request<Greeting>("/greetings", {
    method: "POST",
    body: JSON.stringify(input)
  });
}
