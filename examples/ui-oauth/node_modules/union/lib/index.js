/*
 * index.js: Top-level plugin exposing HTTP features in flatiron
 *
 * (C) 2011, Charlie Robbins & the Contributors
 * MIT LICENSE
 *
 */

var union = exports;

//
// Expose version information
//
exports.version = require('../package.json').version;

//
// Expose core union components
//
union.BufferedStream = require('./buffered-stream');
union.HttpStream     = require('./http-stream');
union.ResponseStream = require('./response-stream');
union.RoutingStream  = require('./routing-stream');
union.createServer   = require('./core').createServer;
union.errorHandler   = require('./core').errorHandler;
