// UI helpers for the products page (group filter multiselect + per-row single select).
(function () {
  function updateGroupSummary(root) {
    if (!root) return;
    var countEl = root.querySelector(".ms-count");
    if (!countEl) return;
    var checked = root.querySelectorAll("input[type=checkbox][name=group_id]:checked");
    countEl.textContent = String(checked.length);
  }

  document.addEventListener("input", function (e) {
    var t = e.target;
    if (!t || t.getAttribute("data-ms-filter") !== "groups") return;
    var root = t.closest(".multiselect");
    if (!root) return;
    var q = (t.value || "").toLowerCase().trim();
    var opts = root.querySelectorAll(".ms-option");
    for (var i = 0; i < opts.length; i++) {
      var text = (opts[i].textContent || "").toLowerCase();
      opts[i].style.display = q === "" || text.indexOf(q) !== -1 ? "" : "none";
    }
  });

  document.addEventListener("change", function (e) {
    var t = e.target;
    if (!t) return;

    // Single-select dropdown in the table row (radio buttons).
    if (t.name === "group_id" && t.type === "radio") {
      var root = t.closest("details.singleselect");
      if (!root) return;
      var current = root.querySelector(".ss-current");
      if (current) {
        var lbl = t.closest("label");
        current.textContent = (lbl ? lbl.textContent || "" : "").trim();
      }
      root.removeAttribute("open");
      return;
    }

    if (t.name !== "group_id" || t.type !== "checkbox") return;
    var ms = t.closest(".multiselect");
    updateGroupSummary(ms);
  });

  document.addEventListener("DOMContentLoaded", function () {
    var roots = document.querySelectorAll("details.multiselect");
    for (var i = 0; i < roots.length; i++) updateGroupSummary(roots[i]);
  });

  // Update group summary after HTMX content swap
  document.addEventListener("htmx:afterSwap", function (e) {
    var root = e.detail.target;
    if (!root) return;
    var ms = root.querySelector("details.multiselect");
    if (ms) updateGroupSummary(ms);
  });

  document.addEventListener("click", function (e) {
    var target = e.target;
    var open = document.querySelectorAll("details.multiselect[open]");
    for (var i = 0; i < open.length; i++) {
      if (open[i].contains(target)) continue;
      open[i].removeAttribute("open");
    }
  });
})();

