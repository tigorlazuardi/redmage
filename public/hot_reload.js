{
  let eventSource = null;
  let id;

  const sessions = JSON.parse(
    localStorage.getItem("hot_reload_sessions") || "{}",
  );

  for (const key in sessions) {
    if (!sessions[key].connected) {
      id = key;
      sessions[id].connected = true;
      break;
    }
  }

  if (!id) {
    id = Math.random().toString(36).substring(2);
    sessions[id] = { connected: true };
  }

  localStorage.setItem("hot_reload_sessions", JSON.stringify(sessions));

  const connect = function () {
    eventSource = new EventSource(`/hot_reload?id=${id}`);
    eventSource.onmessage = function () {
      window.location.reload();
    };
    eventSource.onerror = function () {
      eventSource.close();
      sessions[id].connected = false;
      localStorage.setItem("hot_reload_sessions", JSON.stringify(sessions));
      setTimeout(connect, 50);
    };
  };

  connect();
}
