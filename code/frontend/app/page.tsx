import HelloApiDemo from "../components/HelloApiDemo";
import HelloPage from "../components/HelloPage";
import PersistGreetings from "../components/PersistGreetings";

export default function Home() {
  return (
    <main className="app-shell">
      <HelloPage />
      <HelloApiDemo />
      <PersistGreetings />
    </main>
  );
}
