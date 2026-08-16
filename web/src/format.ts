const E18 = 1_000_000_000_000_000_000n;

export function groupInteger(value: string): string {
  const normalized = BigInt(value).toString();
  return normalized.replace(/\B(?=(\d{3})+(?!\d))/g, "\u2009");
}

export function formatE18(value: string): string {
  const amount = BigInt(value);
  const sign = amount < 0n ? "-" : "";
  const absolute = amount < 0n ? -amount : amount;
  const whole = absolute / E18;
  const fraction = (absolute % E18).toString().padStart(18, "0").replace(/0+$/, "");
  return `${sign}${whole}${fraction ? `.${fraction}` : ""}e18`;
}

export function shortAddress(address: string): string {
  if (address.length < 14) return address;
  return `${address.slice(0, 8)}…${address.slice(-6)}`;
}

export function percentFromRatio(numerator: string, denominator: string): string {
  const bottom = BigInt(denominator);
  if (bottom === 0n) return "0.00%";
  const hundredths = (BigInt(numerator) * 10_000n) / bottom;
  return formatHundredths(hundredths);
}

function formatHundredths(value: bigint): string {
  const whole = value / 100n;
  const fraction = (value % 100n).toString().padStart(2, "0");
  return `${whole}.${fraction}%`;
}

export function formatSyncTime(value: string | null): string {
  if (!value) return "Not yet";
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return "Unknown";
  return new Intl.RelativeTimeFormat(undefined, { numeric: "auto" }).format(
    -Math.max(0, Math.round((Date.now() - date.getTime()) / 1_000)),
    "second",
  );
}
