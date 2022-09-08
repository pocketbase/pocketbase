/**
 * Specification: http://www.w3.org/TR/2012/WD-cors-20120403/
 * W3C Working Draft 3 April 2012
 */
"use strict";

/*jshint node:true */

var simpleMethods, simpleRequestHeaders, simpleResponseHeaders, toLowerCase, checkOriginMatch;

// A method is said to be a simple method if it is a case-sensitive match for one of the following:
Object.defineProperty(exports, "simpleMethods", {
    get: function () {
        return [
            "GET",
            "HEAD",
            "POST"
        ];
    }
});
simpleMethods = exports.simpleMethods;

// A header is said to be a simple header if the header field name is an ASCII case-insensitive match for one of
// the following:
Object.defineProperty(exports, "simpleRequestHeaders", {
    get: function () {
        return [
            "accept",
            "accept-language",
            "content-language",
            "content-type"
        ];
    }
});
simpleRequestHeaders = exports.simpleRequestHeaders;

// A header is said to be a simple response header if the header field name is an ASCII case-insensitive
// match for one of the following:
Object.defineProperty(exports, "simpleResponseHeaders", {
    get: function () {
        return [
            "cache-control",
            "content-language",
            "content-type",
            "expires",
            "last-modified",
            "pragma"
        ];
    }
});
simpleResponseHeaders = exports.simpleResponseHeaders;

toLowerCase = function (array) {
    return array.map(function (el) {
        return el.toLowerCase();
    });
};

checkOriginMatch = function (originHeader, origins, callback) {
    if (typeof origins === "function") {
        origins(originHeader, function (err, allow) {
            callback(err, allow);
        });
    } else if (origins.length > 0) {
        callback(null, origins.some(function (origin) {
            return origin === originHeader;
        }));
    } else {
        // Always matching is acceptable since the list of origins can be unbounded.
        callback(null, true);
    }
};

exports.create = function (options) {
    options = options || {};
    options.origins = options.origins || [];
    options.methods = options.methods || simpleMethods;
    if (options.hasOwnProperty("requestHeaders") === true) {
        options.requestHeaders = toLowerCase(options.requestHeaders);
    } else {
        options.requestHeaders = simpleRequestHeaders;
    }
    if (options.hasOwnProperty("responseHeaders") === true) {
        options.responseHeaders = toLowerCase(options.responseHeaders);
    } else {
        options.responseHeaders = simpleResponseHeaders;
    }
    options.maxAge = options.maxAge || null;
    options.supportsCredentials = options.supportsCredentials || false;
    if (options.hasOwnProperty("endPreflightRequests") === false) {
        options.endPreflightRequests = true;
    }
    return function (req, res, next) {
        var methodMatches, headersMatch, requestMethod, requestHeaders, exposedHeaders, endPreflight;
        // If the Origin header is not present terminate this set of steps.
        if (!req.headers.hasOwnProperty("origin")) {
            // The request is outside the scope of the CORS specification. If there is no Origin header,
            // it could be a same-origin request. Let's let the user-agent handle this situation.
            next();
        } else {
            // If the value of the Origin header is not a case-sensitive match for any of the values in
            // list of origins, do not set any additional headers and terminate this set of steps.
            checkOriginMatch(req.headers.origin, options.origins, function (err, originMatches) {
                if (err !== null) {
                    next(err);
                } else {
                    if (typeof originMatches !== "boolean" || originMatches === false) {
                        next();
                    } else {
                        // Respond to preflight request.
                        if (req.method === "OPTIONS") {
                            endPreflight = function () {
                                if (options.endPreflightRequests === true) {
                                    res.writeHead(204);
                                    res.end();
                                } else {
                                    next();
                                }
                            };
                            // If there is no Access-Control-Request-Method header or if parsing failed, do not set
                            // any additional headers and terminate this set of steps.
                            if (!req.headers.hasOwnProperty("access-control-request-method")) {
                                endPreflight();
                            } else {
                                requestMethod = req.headers["access-control-request-method"];
                                // If there are no Access-Control-Request-Headers headers let header field-names be the
                                // empty list. If parsing failed do not set any additional headers and terminate this set
                                // of steps.
                                // Checking for an empty header is a workaround for a bug Chrome 52:
                                // https://bugs.chromium.org/p/chromium/issues/detail?id=633729
                                if (req.headers.hasOwnProperty("access-control-request-headers") && req.headers["access-control-request-headers"] !== "") {
                                    requestHeaders = toLowerCase(req.headers["access-control-request-headers"].split(/,\s*/));
                                } else {
                                    requestHeaders = [];
                                }
                                // If method is not a case-sensitive match for any of the values in list of methods do not
                                // set any additional headers and terminate this set of steps.
                                methodMatches = options.methods.indexOf(requestMethod) !== -1;
                                if (methodMatches === false) {
                                    endPreflight();
                                } else {
                                    // If any of the header field-names is not a ASCII case-insensitive match for any of
                                    // the values in list of headers do not set any additional headers and terminate this
                                    // set of steps.
                                    headersMatch = requestHeaders.every(function (requestHeader) {
                                        // Browsers automatically add Origin to Access-Control-Request-Headers. However,
                                        // Origin is not one of the simple request headers. Therefore, the header is
                                        // accepted even if it is not in the list of request headers because CORS would
                                        // not work without it.
                                        if (requestHeader === "origin") {
                                            return true;
                                        } else {
                                            if (options.requestHeaders.indexOf(requestHeader) !== -1) {
                                                return true;
                                            } else {
                                                return false;
                                            }
                                        }
                                    });
                                    if (headersMatch === false) {
                                        endPreflight();
                                    } else {
                                        if (options.supportsCredentials === true) {
                                            // If the resource supports credentials add a single Access-Control-Allow-Origin
                                            // header, with the value of the Origin header as value, and add a single
                                            // Access-Control-Allow-Credentials header with the literal string "true"
                                            // as value.
                                            res.setHeader("Access-Control-Allow-Origin", req.headers.origin);
                                            res.setHeader("Access-Control-Allow-Credentials", "true");
                                        } else {
                                            // Otherwise, add a single Access-Control-Allow-Origin header, with either the
                                            // value of the Origin header or the string "*" as value.
                                            if (options.origins.length > 0 || typeof options.origins === "function") {
                                                res.setHeader("Access-Control-Allow-Origin", req.headers.origin);
                                            } else {
                                                res.setHeader("Access-Control-Allow-Origin", "*");
                                            }
                                        }
                                        // Optionally add a single Access-Control-Max-Age header with as value the amount
                                        // of seconds the user agent is allowed to cache the result of the request.
                                        if (options.maxAge !== null) {
                                            res.setHeader("Access-Control-Max-Age", options.maxAge);
                                        }
                                        // Add one or more Access-Control-Allow-Methods headers consisting of (a subset
                                        // of) the list of methods.
                                        res.setHeader("Access-Control-Allow-Methods", options.methods.join(","));
                                        // Add one or more Access-Control-Allow-Headers headers consisting of (a subset
                                        // of) the list of headers.
                                        res.setHeader("Access-Control-Allow-Headers", options.requestHeaders.join(","));
                                        // And out.
                                        endPreflight();
                                    }
                                }
                            }
                        } else {
                            if (options.supportsCredentials === true) {
                                // If the resource supports credentials add a single Access-Control-Allow-Origin header,
                                // with the value of the Origin header as value, and add a single
                                // Access-Control-Allow-Credentials header with the literal string "true" as value.
                                res.setHeader("Access-Control-Allow-Origin", req.headers.origin);
                                res.setHeader("Access-Control-Allow-Credentials", "true");
                            } else {
                                // Otherwise, add a single Access-Control-Allow-Origin header, with either the value of
                                // the Origin header or the literal string "*" as value.
                                // If the list of origins is empty, use "*" as value.
                                if (options.origins.length > 0 || typeof options.origins === "function") {
                                    res.setHeader("Access-Control-Allow-Origin", req.headers.origin);
                                } else {
                                    res.setHeader("Access-Control-Allow-Origin", "*");
                                }
                            }
                            // If the list of exposed headers is not empty add one or more Access-Control-Expose-Headers
                            // headers, with as values the header field names given in the list of exposed headers.
                            exposedHeaders = options.responseHeaders.filter(function (optionsResponseHeader) {
                                return simpleResponseHeaders.indexOf(optionsResponseHeader) === -1;
                            });
                            if (exposedHeaders.length > 0) {
                                res.setHeader("Access-Control-Expose-Headers", exposedHeaders.join(","));
                            }
                            // And out.
                            next();
                        }
                    }
                }
            });
        }
    };
};
