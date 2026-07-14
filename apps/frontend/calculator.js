/**
 * Calculator API helpers — kept free of DOM so unit tests can import them.
 */

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
 * @param {{
 *   a: number,
 *   b: number,
 *   op: string,
 *   fetchFn?: typeof fetch,
 *   setBusy?: (busy: boolean) => void,
 *   setResult?: (text: string) => void,
 *   setError?: (message: string | null) => void,
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
}) {
  setError(null);
  setBusy(true);
  try {
    const result = await callOperation(fetchFn, op, a, b);
    setResult(`Result: ${result}`);
  } catch (err) {
    setResult("Result: —");
    setError(err.message || "Request failed");
  } finally {
    setBusy(false);
  }
}
