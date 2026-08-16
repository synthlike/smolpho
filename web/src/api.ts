export interface IndexerStatus {
  chainHead: string | null;
  lastSuccessfulSync: string | null;
}

export interface Market {
  totalSupplyAssets: string;
  totalBorrowAssets: string;
}

export interface PositionSummary {
  address: string;
  supplyAssets: string;
  borrowAssets: string;
  collateral: string;
}

export interface PositionsPage {
  checkpoint: { blockNumber: string };
  positions: PositionSummary[];
  nextCursor: string | null;
  total: string;
}

const configuredBase = import.meta.env.VITE_API_URL?.trim() ?? "";
const apiBase = configuredBase.endsWith("/")
  ? configuredBase.slice(0, -1)
  : configuredBase;

export async function getJSON<T>(path: string, signal?: AbortSignal): Promise<T> {
  const response = await fetch(`${apiBase}${path}`, {
    headers: { Accept: "application/json" },
    signal,
  });

  if (!response.ok) {
    let message = `API request failed with status ${response.status}`;
    try {
      const body = (await response.json()) as { error?: string };
      if (body.error) message = body.error;
    } catch {
      // Keep the status-based fallback for a non-JSON error response.
    }
    throw new Error(message);
  }

  return (await response.json()) as T;
}
