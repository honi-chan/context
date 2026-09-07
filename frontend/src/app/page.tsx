"use client";

import { useState } from "react";

export default function Home() {
  const [status, setStatus] = useState("");

  const checkBackend = async () => {
    const response = await fetch("http://localhost:8080/health");
    const data = await response.json();

    setStatus(data.status);
  };

  return (
    <main>
      <button onClick={checkBackend}>Backend接続確認</button>
      <p>{status}</p>
    </main>
  );
}