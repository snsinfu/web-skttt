(function() {
  'use strict';

  document.addEventListener('DOMContentLoaded', function() {
    setup();
  });

  function setup() {
    u('#publisher').ajax(function(err, res) {
    });

    u('#request-socket').ajax(function(err, res) {
      if (err || res.error) {
        return;
      }

      const {key, topic} = res.data;

      const http = window.location.protocol;
      const ws = 'ws' + http.substr(4);
      const host = window.location.host;

      const publisherURI = http + '//' + host + '/' + topic + '/' + key;
      const subscriberURI = ws + '//' + host + '/' + topic;

      console.log(publisherURI);
      console.log(subscriberURI);

      u('#publisher').attr('action', publisherURI);
      u('#subscriber').empty();

      const socket = new WebSocket(subscriberURI);
      socket.addEventListener('message', e => {
        u('#subscriber').prepend(
          u('<li>').text(e.data)
        );
      })
    });
  }
})();
