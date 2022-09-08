/* eslint-disable no-process-env */
/* eslint-disable no-sync */
var https = require('https');
var fs = require('fs');
var core = require('union/lib/core');
var RoutingStream = require('union/lib/routing-stream');

module.exports = function (options) {
  var isArray = Array.isArray(options.after);
  var credentials;

  if (!options) {
    throw new Error('options is required to create a server');
  }

  function requestHandler(req, res) {
    var routingStream = new RoutingStream({
      before: options.before,
      buffer: options.buffer,
      after:
        isArray &&
        options.after.map(function (After) {
          return new After();
        }),
      request: req,
      response: res,
      limit: options.limit,
      headers: options.headers
    });

    routingStream.on('error', function (err) {
      var fn = options.onError || core.errorHandler;
      fn(err, routingStream, routingStream.target, function () {
        routingStream.target.emit('next');
      });
    });

    req.pipe(routingStream);
  }

  var serverOptions;

  serverOptions = options.https;
  if (!serverOptions.key || !serverOptions.cert) {
    throw new Error(
      'Both options key and cert are required.'
    );
  }

  credentials = {
    key: fs.readFileSync(serverOptions.key),
    cert: fs.readFileSync(serverOptions.cert),
    passphrase: process.env.NODE_HTTP_SERVER_SSL_PASSPHRASE
  };

  if (serverOptions.ca) {
    serverOptions.ca = !Array.isArray(serverOptions.ca)
      ? [serverOptions.ca]
      : serverOptions.ca;

    credentials.ca = serverOptions.ca.map(function (ca) {
      return fs.readFileSync(ca);
    });
  }

  return https.createServer(credentials, requestHandler);
};
