type AnalyzeRequest = {
  owner: string;
  repo: string;
};

type AnalyzeResponse = {
  result: string;
};

// AnalyzeForm.tsx から named import するため、
// 必ず export を付ける。
export async function analyzeRepository(
  request: AnalyzeRequest,
): Promise<AnalyzeResponse> {
  const apiUrl = process.env.NEXT_PUBLIC_API_URL;

  if (!apiUrl) {
    throw new Error(
      "NEXT_PUBLIC_API_URL is not defined",
    );
  }

  const response = await fetch(
    `${apiUrl}/api/analyze`,
    {
      method: "POST",

      headers: {
        "Content-Type": "application/json",
      },

      body: JSON.stringify({
        owner: request.owner,
        repo: request.repo,
      }),
    },
  );

  if (!response.ok) {
    throw new Error(
      `analyze repository failed: ${response.status}`,
    );
  }

  return response.json();
}