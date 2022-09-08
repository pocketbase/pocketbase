'use strict';

module.exports = (stat, weakEtag) => {
  let etag = `"${[stat.ino, stat.size, stat.mtime.toISOString()].join('-')}"`;
  if (weakEtag) {
    etag = `W/${etag}`;
  }
  return etag;
};
