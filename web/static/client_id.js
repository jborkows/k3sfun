// Sets a stable per-tab client id for echo-suppression and sends it with all HTMX requests.
(function () {
  var key = "shopping_client_id";
  var id = sessionStorage.getItem(key);
  if (!id) {
    id =
      window.crypto && crypto.randomUUID
        ? crypto.randomUUID()
        : String(Date.now()) + "-" + String(Math.random()).slice(2);
    sessionStorage.setItem(key, id);
  }

  if (document.body) {
    document.body.setAttribute("hx-headers", JSON.stringify({ "X-Client-ID": id }));
    var connect = document.body.getAttribute("sse-connect") || "";
    if (connect && connect.indexOf("client=") === -1) {
      document.body.setAttribute(
        "sse-connect",
        connect + (connect.indexOf("?") === -1 ? "?" : "&") + "client=" + encodeURIComponent(id),
      );
    }
  }
})();

