export type Greeting = {
  id: number;
  name: string;
  message: string;
  created_at: string;
};

const reachable = true;

const helloResponse = {
  message: "Hello, World!",
};

const greetings: Greeting[] = [];

export function apiIsReachable() {
  return reachable;
}

export async function getHelloMessage() {
  return helloResponse;
}

export async function listGreetings() {
  return [...greetings].sort((left, right) => right.id - left.id);
}

export async function createGreeting(input: Omit<Greeting, "id" | "created_at">) {
  const greeting: Greeting = {
    id: greetings.length + 1,
    created_at: new Date().toISOString(),
    name: input.name,
    message: input.message,
  };
  greetings.unshift(greeting);
  return greeting;
}
