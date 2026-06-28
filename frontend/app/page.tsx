import { AeolyzerApp } from "@/components/aeolyzer-app"

export default function Home() {
  // Server Component wrapper isolates the client-side state of AeolyzerApp from the root layout, minimizing initial JS payload size.
  return <AeolyzerApp />
}
