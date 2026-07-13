import { runOperation } from "./calculator.js";

const form = document.getElementById("calc-form");
const aInput = document.getElementById("a");
const bInput = document.getElementById("b");
const resultEl = document.getElementById("result");
const errorEl = document.getElementById("error");
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

buttons.forEach((btn) => {
  btn.addEventListener("click", () => {
    runOperation({
      a: Number(aInput.value),
      b: Number(bInput.value),
      op: btn.dataset.op,
      setBusy,
      setResult,
      setError,
    });
  });
});

form.addEventListener("submit", (e) => {
  e.preventDefault();
});
