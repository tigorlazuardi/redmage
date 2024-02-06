{
  const eventSource = new EventSource("/hot_reload");
  eventSource.onmessage = function () {
    window.location.reload();
  };
}
