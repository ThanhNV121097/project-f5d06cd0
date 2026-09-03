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
};

const apiBase = process.env.NEXT_PUBLIC_API_URL ?? "/api";

async function readJSON<T>(response: Response): Promise<T> {
  if (!response.ok) {
    throw new Error("api unreachable");
  }
  return response.json() as Promise<T>;
}

export function fetchHello() {
  return readJSON<{ message: string }>(await fetch(`${apiBase}/v1/hello`));
}

export async function fetchGreetings() {
  return readJSON<GreetingsResponse>(fetch(`${apiBase}/v1/greetings`));
}

export async function saveGreeting(input: GreetingInput) {
  return readJSON<Greeting>(
    fetch(`${apiBase}/v1/greetings`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(input),
    }),
  );
}
