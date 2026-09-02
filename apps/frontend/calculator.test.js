import assert from "node:assert/strict";
import { mock, test } from "node:test";

import {
  callOperation,
  fetchHistory,
  formatHistoryRow,
  runOperation,
} from "./calculator.js";

function jsonResponse(status, body) {
  return {
    ok: status >= 200 && status < 300,
    status,
    async json() {
      return body;
    },
  };
}

test("callOperation posts to /api/{op} and returns result", async () => {
  const fetchFn = mock.fn(async (url, init) => {
    assert.equal(url, "/api/add");
    assert.equal(init.method, "POST");
    assert.equal(init.headers["Content-Type"], "application/json");
    assert.deepEqual(JSON.parse(init.body), { a: 10, b: 2 });
    return jsonResponse(200, { result: 12 });
  });

  const result = await callOperation(fetchFn, "add", 10, 2);
  assert.equal(result, 12);
  assert.equal(fetchFn.mock.callCount(), 1);
});

test("callOperation throws API error message", async () => {
  const fetchFn = mock.fn(async () =>
    jsonResponse(400, { error: "division by zero" }),
  );

  await assert.rejects(
    () => callOperation(fetchFn, "divide", 1, 0),
    /division by zero/,
  );
});

test("callOperation throws fallback message when error field missing", async () => {
  const fetchFn = mock.fn(async () => jsonResponse(500, {}));

  await assert.rejects(
    () => callOperation(fetchFn, "add", 1, 2),
    /Request failed \(500\)/,
  );
});

test("runOperation success updates result and busy flags", async () => {
  const events = [];
  const fetchFn = mock.fn(async () => jsonResponse(200, { result: 7 }));

  await runOperation({
    a: 3,
    b: 4,
    op: "add",
    fetchFn,
    setBusy: (busy) => events.push(["busy", busy]),
    setResult: (text) => events.push(["result", text]),
    setError: (msg) => events.push(["error", msg]),
    onSuccess: async () => events.push(["success"]),
  });

  assert.deepEqual(events, [
    ["error", null],
    ["busy", true],
    ["result", "Result: 7"],
    ["success"],
    ["busy", false],
  ]);
});

test("runOperation failure clears result and shows error", async () => {
  const events = [];
  const fetchFn = mock.fn(async () =>
    jsonResponse(400, { error: "division by zero" }),
  );

  await runOperation({
    a: 1,
    b: 0,
    op: "divide",
    fetchFn,
    setBusy: (busy) => events.push(["busy", busy]),
    setResult: (text) => events.push(["result", text]),
    setError: (msg) => events.push(["error", msg]),
    onFailure: async (info) =>
      events.push(["failure", info.op, info.a, info.b]),
  });

  assert.deepEqual(events, [
    ["error", null],
    ["busy", true],
    ["result", "Result: —"],
    ["error", "division by zero"],
    ["failure", "divide", 1, 0],
    ["busy", false],
  ]);
});

test("runOperation covers subtract multiply divide paths", async () => {
  for (const [op, a, b, result] of [
    ["subtract", 10, 2, 8],
    ["multiply", 10, 2, 20],
    ["divide", 10, 2, 5],
  ]) {
    const fetchFn = mock.fn(async (url) => {
      assert.equal(url, `/api/${op}`);
      return jsonResponse(200, { result });
    });
    const got = await callOperation(fetchFn, op, a, b);
    assert.equal(got, result);
  }
});

test("fetchHistory returns array from /api/history", async () => {
  const fetchFn = mock.fn(async (url) => {
    assert.equal(url, "/api/history");
    return jsonResponse(200, [{ op: "add", a: 10, b: 2, result: 12 }]);
  });
  const entries = await fetchHistory(fetchFn);
  assert.equal(entries.length, 1);
  assert.equal(entries[0].op, "add");
});

test("formatHistoryRow uses operation symbols", () => {
  assert.equal(
    formatHistoryRow({ op: "add", a: 10, b: 2, result: 12 }),
    "10 + 2 = 12",
  );
  assert.equal(
    formatHistoryRow({ op: "divide", a: 10, b: 2, result: 5 }),
    "10 ÷ 2 = 5",
  );
});
