/**
 * Calculator API helpers — kept free of DOM so unit tests can import them.
 */

const OP_SYMBOLS = {
  add: "+",
  subtract: "−",
  multiply: "×",
  divide: "÷",
};

/**
 * @param {typeof fetch} fetchFn
 * @param {string} op
 * @param {number} a
 * @param {number} b
 * @returns {Promise<number>}
 */
export async function callOperation(fetchFn, op, a, b) {
  const res = await fetchFn(`/api/${op}`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ a, b }),
  });

  const data = await res.json();
  if (!res.ok) {
    throw new Error(data.error || `Request failed (${res.status})`);
  }
  return data.result;
}

/**
 * @param {typeof fetch} [fetchFn]
 * @returns {Promise<Array<{op: string, a: number, b: number, result: number, timestamp?: string}>>}
 */
export async function fetchHistory(fetchFn = fetch) {
  const res = await fetchFn("/api/history");
  if (!res.ok) {
    throw new Error(`History request failed (${res.status})`);
  }
  const data = await res.json();
  return Array.isArray(data) ? data : [];
}

/**
 * @param {{op: string, a: number, b: number, result: number}} entry
 * @returns {string}
 */
export function formatHistoryRow(entry) {
  const symbol = OP_SYMBOLS[entry.op] || entry.op;
  return `${entry.a} ${symbol} ${entry.b} = ${entry.result}`;
}

/**
 * @param {{
 *   a: number,
 *   b: number,
 *   op: string,
 *   fetchFn?: typeof fetch,
 *   setBusy?: (busy: boolean) => void,
 *   setResult?: (text: string) => void,
 *   setError?: (message: string | null) => void,
 *   onSuccess?: () => void | Promise<void>,
 *   onFailure?: (info: { op: string, a: number, b: number, message: string }) => void | Promise<void>,
 * }} opts
 */
export async function runOperation({
  a,
  b,
  op,
  fetchFn = fetch,
  setBusy = () => {},
  setResult = () => {},
  setError = () => {},
  onSuccess = async () => {},
  onFailure = async () => {},
}) {
  setError(null);
  setBusy(true);
  try {
    const result = await callOperation(fetchFn, op, a, b);
    setResult(`Result: ${result}`);
    await onSuccess();
  } catch (err) {
    setResult("Result: —");
    const message = err.message || "Request failed";
    setError(message);
    await onFailure({ op, a, b, message });
  } finally {
    setBusy(false);
  }
}
