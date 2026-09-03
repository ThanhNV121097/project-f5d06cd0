export type Greeting = {
  id: number;
  name: string;
  message: string;
  created_at: string;
};

const greetings: Greeting[] = [
  { id: 2, name: "Lin", message: "Hello, World!", created_at: "2025-01-15 09:58" },
  { id: 1, name: "Ada", message: "Hello from the database", created_at: "2025-01-15 10:24" },
];

export async function getHello(name = "") {
  const clean = name.trim();
  const label = clean || "World";
  return { message: `Hello, ${label}!` };
}

export async function getGreetings() {
  return [...greetings].sort((a, b) => b.id - a.id);
}

export async function createGreeting(input: Pick<Greeting, "name" | "message">) {
  const next = { id: greetings[0].id + 1, created_at: "2025-01-15 11:05", ...input };
  greetings.unshift(next);
  return next;
}
