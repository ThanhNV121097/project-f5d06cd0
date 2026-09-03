export type Greeting = {
  id: string;
  name: string;
  message: string;
  created_at: string;
};

const greetings: Greeting[] = [
  {
    id: "g_2",
    name: "Lin",
    message: "Hello, World!",
    created_at: "2025-01-15T09:58:00Z",
  },
  {
    id: "g_1",
    name: "Ada",
    message: "Hello from the database",
    created_at: "2025-01-15T10:24:00Z",
  },
];

export async function fetchHelloMessage() {
  return { message: "Hello, World!" };
}

export async function fetchGreetings() {
  return { greetings: [...greetings] };
}

export async function saveGreeting(input: { name: string; message: string }) {
  const greeting: Greeting = {
    id: `g_${Date.now()}`,
    name: input.name.trim(),
    message: input.message.trim(),
    created_at: new Date().toISOString(),
  };

  greetings.unshift(greeting);
  return greeting;
}

