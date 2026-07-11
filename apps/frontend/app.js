const form = document.getElementById("calc-form");
const aInput = document.getElementById("a");
const bInput = document.getElementById("b");
const resultEl = document.getElementById("result");
const errorEl = document.getElementById("error");
const buttons = form.querySelectorAll("button[data-op]");

async function runOperation(op) {
  const a = Number(aInput.value);
  const b = Number(bInput.value);

  errorEl.hidden = true;
  errorEl.textContent = "";
  buttons.forEach((btn) => {
    btn.disabled = true;
  });

  try {
    const res = await fetch(`/api/${op}`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ a, b }),
    });

    const data = await res.json();
    if (!res.ok) {
      throw new Error(data.error || `Request failed (${res.status})`);
    }

    resultEl.textContent = `Result: ${data.result}`;
  } catch (err) {
    resultEl.textContent = "Result: —";
    errorEl.textContent = err.message || "Request failed";
    errorEl.hidden = false;
  } finally {
    buttons.forEach((btn) => {
      btn.disabled = false;
    });
  }
}

buttons.forEach((btn) => {
  btn.addEventListener("click", () => {
    runOperation(btn.dataset.op);
  });
});

form.addEventListener("submit", (e) => {
  e.preventDefault();
});
