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

const greetings: Greeting[] = [
  { id: "2", name: "Ada", message: "Hello from the browser", created_at: "2026-09-03T12:30:00Z" },
  { id: "1", name: "Lin", message: "Hello, World!", created_at: "2026-09-03T12:00:00Z" },
];

export async function fetchHello() {
  return { message: "Hello, World!" };
}

export async function fetchGreetings() {
  return { greetings: [...greetings] };
}

export async function saveGreeting(input: GreetingInput) {
  const next: Greeting = {
    id: String(greetings.length + 1),
    name: input.name.trim(),
    message: input.message.trim(),
    created_at: new Date().toISOString(),
  };

  greetings.unshift(next);
  return next;
}
