const $ = (s) => document.querySelector(s);

const screens = {
  loading: $("#loading"),
  login: $("#login-screen"),
  main: $("#main-screen"),
};

function showScreen(name) {
  Object.values(screens).forEach((s) => s.classList.add("hidden"));
  screens[name].classList.remove("hidden");
}

function showToast(text) {
  const toast = document.createElement("div");
  toast.className = "toast";
  toast.textContent = text;
  document.body.appendChild(toast);
  setTimeout(() => toast.remove(), 2000);
}

// ─── Init ─────────────────────────────────────────────────

chrome.runtime.sendMessage({ action: "getState" }, (state) => {
  if (state?.authenticated) {
    renderMain(state);
  } else {
    showScreen("login");
  }
});

// ─── Login ────────────────────────────────────────────────

$("#login-form").addEventListener("submit", async (e) => {
  e.preventDefault();
  const btn = $("#login-btn");
  const errorEl = $("#login-error");
  errorEl.classList.add("hidden");
  btn.disabled = true;
  btn.textContent = "Signing in...";

  chrome.runtime.sendMessage(
    {
      action: "login",
      email: $("#email").value,
      password: $("#password").value,
    },
    (resp) => {
      btn.disabled = false;
      btn.textContent = "Sign in";

      if (resp?.error) {
        errorEl.textContent = resp.error;
        errorEl.classList.remove("hidden");
      } else {
        chrome.runtime.sendMessage({ action: "getState" }, (state) => {
          renderMain(state);
        });
      }
    }
  );
});

// ─── Main Screen ──────────────────────────────────────────

function renderMain(state) {
  showScreen("main");
  $("#user-email").textContent = state.user?.email || "";
  renderAddresses(state.addresses || []);
}

function renderAddresses(addresses) {
  const list = $("#address-list");

  if (addresses.length === 0) {
    list.innerHTML = `
      <div class="no-addresses">
        No addresses yet.<br>
        <a href="https://addrpass.com/dashboard" target="_blank" style="color:#22D3EE">Add one in the dashboard</a>
      </div>
    `;
    return;
  }

  list.innerHTML = addresses
    .map(
      (a) => `
    <div class="address-card" data-id="${a.id}">
      ${a.label ? `<div class="address-label">${a.label}</div>` : ""}
      <div class="address-line">${a.line1}</div>
      <div class="address-detail">${a.city}${a.state ? `, ${a.state}` : ""} ${a.post_code}, ${a.country}</div>
      <div class="address-actions">
        <button class="fill-btn" data-action="fill" data-id="${a.id}">Fill Form</button>
        <button data-action="copy" data-id="${a.id}">Copy</button>
        <button data-action="share" data-id="${a.id}">Share</button>
      </div>
    </div>
  `
    )
    .join("");

  // Bind actions
  list.querySelectorAll("button[data-action]").forEach((btn) => {
    btn.addEventListener("click", () => {
      const id = btn.dataset.id;
      const addr = addresses.find((a) => a.id === id);
      if (!addr) return;

      switch (btn.dataset.action) {
        case "fill":
          chrome.runtime.sendMessage({ action: "fillAddress", address: addr });
          showToast("Address filled!");
          setTimeout(() => window.close(), 500);
          break;

        case "copy": {
          const text = [addr.line1, addr.line2, `${addr.city}${addr.state ? `, ${addr.state}` : ""} ${addr.post_code}`, addr.country, addr.phone]
            .filter(Boolean)
            .join("\n");
          navigator.clipboard.writeText(text);
          showToast("Copied!");
          break;
        }

        case "share":
          chrome.runtime.sendMessage(
            { action: "createShare", addressId: addr.id, scope: "full" },
            (resp) => {
              if (resp?.url) {
                navigator.clipboard.writeText(resp.url);
                showToast("Share link copied!");
              } else {
                showToast("Failed to create share");
              }
            }
          );
          break;
      }
    });
  });
}

// ─── Logout ───────────────────────────────────────────────

$("#logout-btn").addEventListener("click", () => {
  chrome.runtime.sendMessage({ action: "logout" }, () => {
    showScreen("login");
  });
});

// ─── Refresh ──────────────────────────────────────────────

$("#refresh-btn").addEventListener("click", () => {
  chrome.runtime.sendMessage({ action: "refreshAddresses" }, (resp) => {
    if (resp?.addresses) {
      renderAddresses(resp.addresses);
      showToast("Refreshed!");
    }
  });
});
