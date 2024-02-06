{
  let eventSource = null;
  function connect() {
    eventSource = new EventSource("/hot_reload");
    eventSource.onmessage = function () {
      window.location.reload();
    };
    eventSource.onerror = function () {
      eventSource.close();
      setTimeout(connect, 50);
    };
  }

  connect();
}
