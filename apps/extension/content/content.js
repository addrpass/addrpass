// ─── Address Field Detection ──────────────────────────────
//
// Strategy: Match fields using multiple signals (prioritized):
// 1. HTML autocomplete attribute (strongest signal)
// 2. name/id attribute patterns
// 3. placeholder text
// 4. associated label text

const FIELD_MAP = {
  line1: {
    autocomplete: ["street-address", "address-line1", "shipping street-address", "billing street-address"],
    patterns: /address.?1|street.?addr|addr.?line.?1|address.?line|street|sokak|cadde|adres/i,
  },
  line2: {
    autocomplete: ["address-line2"],
    patterns: /address.?2|addr.?line.?2|apt|suite|unit|apartment|flat|daire/i,
  },
  city: {
    autocomplete: ["address-level2"],
    patterns: /city|town|locality|il[cç]e|sehir|şehir|stadt|ville|ciudad/i,
  },
  state: {
    autocomplete: ["address-level1"],
    patterns: /state|province|region|county|prefecture|il\b|eyalet|bundesland/i,
  },
  postCode: {
    autocomplete: ["postal-code"],
    patterns: /zip|postal|postcode|post.?code|posta.?kodu|plz/i,
  },
  country: {
    autocomplete: ["country", "country-name"],
    patterns: /country|nation|ulke|ülke|land|pays|país/i,
  },
  phone: {
    autocomplete: ["tel", "tel-national"],
    patterns: /phone|tel|mobile|telefon|handy|t[ée]l[ée]phone/i,
  },
};

function getFieldSignals(el) {
  const autocomplete = (el.getAttribute("autocomplete") || "").toLowerCase().trim();
  const name = (el.getAttribute("name") || "").toLowerCase();
  const id = (el.getAttribute("id") || "").toLowerCase();
  const placeholder = (el.getAttribute("placeholder") || "").toLowerCase();

  // Check associated label
  let labelText = "";
  if (el.id) {
    const label = document.querySelector(`label[for="${el.id}"]`);
    if (label) labelText = label.textContent.toLowerCase().trim();
  }
  // Also check parent label
  const parentLabel = el.closest("label");
  if (parentLabel) labelText = parentLabel.textContent.toLowerCase().trim();

  return { autocomplete, name, id, placeholder, labelText };
}

function classifyField(el) {
  const signals = getFieldSignals(el);

  for (const [fieldName, config] of Object.entries(FIELD_MAP)) {
    // Priority 1: autocomplete attribute
    if (config.autocomplete.some((ac) => signals.autocomplete.includes(ac))) {
      return fieldName;
    }

    // Priority 2: name/id pattern
    if (config.patterns.test(signals.name) || config.patterns.test(signals.id)) {
      return fieldName;
    }

    // Priority 3: placeholder
    if (config.patterns.test(signals.placeholder)) {
      return fieldName;
    }

    // Priority 4: label text
    if (config.patterns.test(signals.labelText)) {
      return fieldName;
    }
  }

  return null;
}

function detectAddressFields() {
  const inputs = document.querySelectorAll(
    'input[type="text"], input[type="tel"], input[type="search"], input:not([type]), select, textarea'
  );

  const fields = {};
  let count = 0;

  for (const el of inputs) {
    // Skip hidden, disabled, readonly
    if (el.offsetParent === null || el.disabled || el.readOnly) continue;
    // Skip very small inputs (likely search bars)
    if (el.offsetWidth < 50) continue;

    const fieldName = classifyField(el);
    if (fieldName && !fields[fieldName]) {
      fields[fieldName] = el;
      count++;
    }
  }

  // Need at least 2 address fields to consider it a form
  return count >= 2 ? fields : null;
}

// ─── Fill Logic ───────────────────────────────────────────

function fillField(el, value) {
  if (!value) return;

  if (el.tagName === "SELECT") {
    // Try to match option by value or text
    for (const opt of el.options) {
      if (
        opt.value.toLowerCase() === value.toLowerCase() ||
        opt.text.toLowerCase().includes(value.toLowerCase()) ||
        value.toLowerCase().includes(opt.text.toLowerCase())
      ) {
        el.value = opt.value;
        el.dispatchEvent(new Event("change", { bubbles: true }));
        return;
      }
    }
    // Try 2-letter country code match
    if (value.length === 2) {
      for (const opt of el.options) {
        if (opt.value.toUpperCase() === value.toUpperCase()) {
          el.value = opt.value;
          el.dispatchEvent(new Event("change", { bubbles: true }));
          return;
        }
      }
    }
  } else {
    // Simulate realistic typing for React/Vue/Angular forms
    const nativeSetter = Object.getOwnPropertyDescriptor(
      el.tagName === "TEXTAREA" ? HTMLTextAreaElement.prototype : HTMLInputElement.prototype,
      "value"
    )?.set;

    if (nativeSetter) {
      nativeSetter.call(el, value);
    } else {
      el.value = value;
    }

    el.dispatchEvent(new Event("input", { bubbles: true }));
    el.dispatchEvent(new Event("change", { bubbles: true }));
    el.dispatchEvent(new Event("blur", { bubbles: true }));
  }
}

function fillAddressFields(fields, address) {
  const mapping = {
    line1: address.line1,
    line2: address.line2,
    city: address.city,
    state: address.state,
    postCode: address.post_code,
    country: address.country,
    phone: address.phone,
  };

  for (const [fieldName, el] of Object.entries(fields)) {
    fillField(el, mapping[fieldName]);
  }
}

// ─── Badge Overlay ────────────────────────────────────────

let badge = null;

function showBadge(fields) {
  if (badge) return;

  // Find the first address field to position near
  const firstField = Object.values(fields)[0];
  if (!firstField) return;

  badge = document.createElement("div");
  badge.id = "addrpass-badge";
  badge.innerHTML = `
    <svg width="16" height="16" viewBox="0 0 32 32" fill="none">
      <path d="M16 2L4 8v8c0 8.4 5.12 16.24 12 18 6.88-1.76 12-9.6 12-18V8L16 2z" fill="#0F172A"/>
      <circle cx="16" cy="13" r="3" fill="#22D3EE"/>
      <path d="M14 15.5L13 22h6l-1-6.5" fill="#22D3EE" opacity="0.7"/>
    </svg>
    <span>Fill with AddrPass</span>
  `;
  badge.addEventListener("click", () => {
    chrome.runtime.sendMessage({ action: "getState" }, (state) => {
      if (state?.authenticated && state.addresses?.length > 0) {
        // Fill with first address (or show picker if multiple)
        if (state.addresses.length === 1) {
          fillAddressFields(fields, state.addresses[0]);
          removeBadge();
        } else {
          showAddressPicker(fields, state.addresses);
        }
      } else {
        // Not logged in — open popup
        chrome.runtime.sendMessage({ action: "openPopup" });
      }
    });
  });

  document.body.appendChild(badge);
}

function removeBadge() {
  if (badge) {
    badge.remove();
    badge = null;
  }
}

function showAddressPicker(fields, addresses) {
  removeBadge();

  const picker = document.createElement("div");
  picker.id = "addrpass-picker";
  picker.innerHTML = `
    <div class="addrpass-picker-header">
      <svg width="14" height="14" viewBox="0 0 32 32" fill="none">
        <path d="M16 2L4 8v8c0 8.4 5.12 16.24 12 18 6.88-1.76 12-9.6 12-18V8L16 2z" fill="#0F172A"/>
        <circle cx="16" cy="13" r="3" fill="#22D3EE"/>
        <path d="M14 15.5L13 22h6l-1-6.5" fill="#22D3EE" opacity="0.7"/>
      </svg>
      <span>Select address</span>
      <button id="addrpass-picker-close">&times;</button>
    </div>
    <div class="addrpass-picker-list">
      ${addresses
        .map(
          (a, i) => `
        <div class="addrpass-picker-item" data-index="${i}">
          ${a.label ? `<span class="addrpass-picker-label">${a.label}</span>` : ""}
          <span class="addrpass-picker-line">${a.line1}</span>
          <span class="addrpass-picker-detail">${a.city}, ${a.country}</span>
        </div>
      `
        )
        .join("")}
    </div>
  `;

  picker.querySelector("#addrpass-picker-close").addEventListener("click", () => picker.remove());

  picker.querySelectorAll(".addrpass-picker-item").forEach((item) => {
    item.addEventListener("click", () => {
      const idx = parseInt(item.dataset.index);
      fillAddressFields(fields, addresses[idx]);
      picker.remove();
    });
  });

  document.body.appendChild(picker);
}

// ─── Listen for Fill Commands from Background ─────────────

chrome.runtime.onMessage.addListener((msg) => {
  if (msg.action === "fill") {
    const fields = detectAddressFields();
    if (fields) {
      fillAddressFields(fields, msg.address);
    }
  }
});

// ─── Init: Detect Forms on Page Load ──────────────────────

function init() {
  // Wait a bit for SPAs to render
  setTimeout(() => {
    const fields = detectAddressFields();
    if (fields) {
      showBadge(fields);
    }
  }, 1500);
}

// Also watch for dynamic content (SPAs)
const observer = new MutationObserver(() => {
  if (!badge && !document.getElementById("addrpass-picker")) {
    const fields = detectAddressFields();
    if (fields) {
      showBadge(fields);
    }
  }
});

observer.observe(document.body, { childList: true, subtree: true });
init();
