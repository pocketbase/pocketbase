'use strict';

const fs = require('fs');
const path = require('path');

module.exports = function sortByIsDirectory(dir, paths, cb) {
  // take the listing file names in `dir`
  // returns directory and file array, each entry is
  // of the array a [name, stat] tuple
  let pending = paths.length;
  const errs = [];
  const dirs = [];
  const files = [];

  if (!pending) {
    cb(errs, dirs, files);
    return;
  }

  paths.forEach((file) => {
    fs.stat(path.join(dir, file), (err, s) => {
      if (err) {
        errs.push([file, err]);
      } else if (s.isDirectory()) {
        dirs.push([file, s]);
      } else {
        files.push([file, s]);
      }

      pending -= 1;
      if (pending === 0) {
        cb(errs, dirs, files);
      }
    });
  });
};
