import { fetchHistory, formatHistoryRow, runOperation } from "./calculator.js";

const form = document.getElementById("calc-form");
const aInput = document.getElementById("a");
const bInput = document.getElementById("b");
const resultEl = document.getElementById("result");
const errorEl = document.getElementById("error");
const historyListEl = document.getElementById("history-list");
const historyEmptyEl = document.getElementById("history-empty");
const buttons = form.querySelectorAll("button[data-op]");

function setBusy(busy) {
  buttons.forEach((btn) => {
    btn.disabled = busy;
  });
}

function setResult(text) {
  resultEl.textContent = text;
}

function setError(message) {
  if (message === null) {
    errorEl.hidden = true;
    errorEl.textContent = "";
    return;
  }
  errorEl.textContent = message;
  errorEl.hidden = false;
}

function renderHistory(entries) {
  historyListEl.replaceChildren();
  if (!entries.length) {
    historyEmptyEl.hidden = false;
    return;
  }
  historyEmptyEl.hidden = true;
  for (const entry of entries) {
    const li = document.createElement("li");
    li.textContent = formatHistoryRow(entry);
    historyListEl.appendChild(li);
  }
}

async function loadHistory() {
  try {
    const entries = await fetchHistory();
    renderHistory(entries);
  } catch {
    // Keep whatever is on screen; do not clear the result/error area.
  }
}

function appendOptimisticFailure({ op, a, b }) {
  historyEmptyEl.hidden = true;
  const li = document.createElement("li");
  li.textContent = formatHistoryRow({ op, a, b, result: 0 });
  historyListEl.prepend(li);
}

buttons.forEach((btn) => {
  btn.addEventListener("click", () => {
    runOperation({
      a: Number(aInput.value),
      b: Number(bInput.value),
      op: btn.dataset.op,
      setBusy,
      setResult,
      setError,
      onSuccess: loadHistory,
      onFailure: appendOptimisticFailure,
    });
  });
});

form.addEventListener("submit", (e) => {
  e.preventDefault();
});

loadHistory();
