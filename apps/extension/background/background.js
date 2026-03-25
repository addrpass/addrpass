const API_URL = "https://api.addrpass.com";

// ─── Auth ─────────────────────────────────────────────────

async function getToken() {
  const { addrpass_token } = await chrome.storage.local.get("addrpass_token");
  return addrpass_token || null;
}

async function setToken(token) {
  await chrome.storage.local.set({ addrpass_token: token });
}

async function clearToken() {
  await chrome.storage.local.remove(["addrpass_token", "addrpass_user", "addrpass_addresses"]);
}

async function apiFetch(path, opts = {}) {
  const headers = { "Content-Type": "application/json" };
  const token = await getToken();
  if (token) headers["Authorization"] = `Bearer ${token}`;

  const res = await fetch(`${API_URL}${path}`, {
    method: opts.method || "GET",
    headers,
    body: opts.body ? JSON.stringify(opts.body) : undefined,
  });

  if (!res.ok) {
    const err = await res.json().catch(() => ({ error: res.statusText }));
    throw new Error(err.error || "Request failed");
  }
  if (res.status === 204) return {};
  return res.json();
}

// ─── Message Handler ──────────────────────────────────────

chrome.runtime.onMessage.addListener((msg, sender, sendResponse) => {
  handleMessage(msg).then(sendResponse).catch((e) => sendResponse({ error: e.message }));
  return true; // Keep channel open for async response
});

async function handleMessage(msg) {
  switch (msg.action) {
    case "login": {
      const resp = await apiFetch("/api/v1/auth/login", {
        method: "POST",
        body: { email: msg.email, password: msg.password },
      });
      await setToken(resp.token);
      await chrome.storage.local.set({ addrpass_user: resp.user });
      // Pre-fetch addresses
      const addresses = await apiFetch("/api/v1/addresses");
      await chrome.storage.local.set({ addrpass_addresses: addresses });
      return { success: true, user: resp.user };
    }

    case "logout": {
      await clearToken();
      return { success: true };
    }

    case "getState": {
      const token = await getToken();
      if (!token) return { authenticated: false };
      const { addrpass_user, addrpass_addresses } = await chrome.storage.local.get([
        "addrpass_user",
        "addrpass_addresses",
      ]);
      return { authenticated: true, user: addrpass_user, addresses: addrpass_addresses || [] };
    }

    case "refreshAddresses": {
      const addresses = await apiFetch("/api/v1/addresses");
      await chrome.storage.local.set({ addrpass_addresses: addresses });
      return { addresses };
    }

    case "createShare": {
      const resp = await apiFetch("/api/v1/shares", {
        method: "POST",
        body: {
          address_id: msg.addressId,
          access_type: "public",
          scope: msg.scope || "full",
          max_accesses: msg.maxAccesses || null,
        },
      });
      return resp;
    }

    case "fillAddress": {
      // Send fill command to content script
      const [tab] = await chrome.tabs.query({ active: true, currentWindow: true });
      if (tab) {
        chrome.tabs.sendMessage(tab.id, {
          action: "fill",
          address: msg.address,
        });
      }
      return { success: true };
    }

    default:
      return { error: "Unknown action" };
  }
}
