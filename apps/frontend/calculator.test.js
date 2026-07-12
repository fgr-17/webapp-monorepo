import assert from "node:assert/strict";
import { mock, test } from "node:test";

import { callOperation, runOperation } from "./calculator.js";

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
  });

  assert.deepEqual(events, [
    ["error", null],
    ["busy", true],
    ["result", "Result: 7"],
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
  });

  assert.deepEqual(events, [
    ["error", null],
    ["busy", true],
    ["result", "Result: —"],
    ["error", "division by zero"],
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
