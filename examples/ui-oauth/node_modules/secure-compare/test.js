/**
 * Dependencies
 */

var compare = require('./');

require('chai').should();


/**
 * Tests
 */

describe ('secure-compare', function () {
  it ('compare', function () {
    compare('abc', 'abc').should.equal(true);
    compare('abc', 'ab').should.equal(false);
  });
});
