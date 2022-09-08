#!/usr/bin/env node

'use strict';

var chalk     = require('chalk'),
    os         = require('os'),
    httpServer = require('../lib/http-server'),
    portfinder = require('portfinder'),
    opener     = require('opener'),

    fs         = require('fs'),
    url        = require('url');
var argv = require('minimist')(process.argv.slice(2), {
  alias: {
    tls: 'ssl'
  }
});
var ifaces = os.networkInterfaces();

process.title = 'http-server';

if (argv.h || argv.help) {
  console.log([
    'usage: http-server [path] [options]',
    '',
    'options:',
    '  -p --port    Port to use. If 0, look for open port. [8080]',
    '  -a           Address to use [0.0.0.0]',
    '  -d           Show directory listings [true]',
    '  -i           Display autoIndex [true]',
    '  -g --gzip    Serve gzip files when possible [false]',
    '  -b --brotli  Serve brotli files when possible [false]',
    '               If both brotli and gzip are enabled, brotli takes precedence',
    '  -e --ext     Default file extension if none supplied [none]',
    '  -s --silent  Suppress log messages from output',
    '  --cors[=headers]   Enable CORS via the "Access-Control-Allow-Origin" header',
    '                     Optionally provide CORS headers list separated by commas',
    '  -o [path]    Open browser window after starting the server.',
    '               Optionally provide a URL path to open the browser window to.',
    '  -c           Cache time (max-age) in seconds [3600], e.g. -c10 for 10 seconds.',
    '               To disable caching, use -c-1.',
    '  -t           Connections timeout in seconds [120], e.g. -t60 for 1 minute.',
    '               To disable timeout, use -t0',
    '  -U --utc     Use UTC time format in log messages.',
    '  --log-ip     Enable logging of the client\'s IP address',
    '',
    '  -P --proxy       Fallback proxy if the request cannot be resolved. e.g.: http://someurl.com',
    '  --proxy-options  Pass options to proxy using nested dotted objects. e.g.: --proxy-options.secure false',
    '',
    '  --username   Username for basic authentication [none]',
    '               Can also be specified with the env variable NODE_HTTP_SERVER_USERNAME',
    '  --password   Password for basic authentication [none]',
    '               Can also be specified with the env variable NODE_HTTP_SERVER_PASSWORD',
    '',
    '  -S --tls --ssl   Enable secure request serving with TLS/SSL (HTTPS)',
    '  -C --cert    Path to TLS cert file (default: cert.pem)',
    '  -K --key     Path to TLS key file (default: key.pem)',
    '',
    '  -r --robots        Respond to /robots.txt [User-agent: *\\nDisallow: /]',
    '  --no-dotfiles      Do not show dotfiles',
    '  --mimetypes        Path to a .types file for custom mimetype definition',
    '  -h --help          Print this list and exit.',
    '  -v --version       Print the version and exit.'
  ].join('\n'));
  process.exit();
}

var port = argv.p || argv.port || parseInt(process.env.PORT, 10),
    host = argv.a || '0.0.0.0',
    tls = argv.S || argv.tls,
    sslPassphrase = process.env.NODE_HTTP_SERVER_SSL_PASSPHRASE,
    proxy = argv.P || argv.proxy,
    proxyOptions = argv['proxy-options'],
    utc = argv.U || argv.utc,
    version = argv.v || argv.version,
    logger;

var proxyOptionsBooleanProps = [
  'ws', 'xfwd', 'secure', 'toProxy', 'prependPath', 'ignorePath', 'changeOrigin',
  'preserveHeaderKeyCase', 'followRedirects', 'selfHandleResponse'
];

if (proxyOptions) {
  Object.keys(proxyOptions).forEach(function (key) {
    if (proxyOptionsBooleanProps.indexOf(key) > -1) {
      proxyOptions[key] = proxyOptions[key].toLowerCase() === 'true';
    }
  });
}

if (!argv.s && !argv.silent) {
  logger = {
    info: console.log,
    request: function (req, res, error) {
      var date = utc ? new Date().toUTCString() : new Date();
      var ip = argv['log-ip']
          ? req.headers['x-forwarded-for'] || '' +  req.connection.remoteAddress
          : '';
      if (error) {
        logger.info(
          '[%s] %s "%s %s" Error (%s): "%s"',
          date, ip, chalk.red(req.method), chalk.red(req.url),
          chalk.red(error.status.toString()), chalk.red(error.message)
        );
      }
      else {
        logger.info(
          '[%s] %s "%s %s" "%s"',
          date, ip, chalk.cyan(req.method), chalk.cyan(req.url),
          req.headers['user-agent']
        );
      }
    }
  };
}
else if (chalk) {
  logger = {
    info: function () {},
    request: function () {}
  };
}

if (version) {
  logger.info('v' + require('../package.json').version);
  process.exit();
}

if (!port) {
  portfinder.basePort = 8080;
  portfinder.getPort(function (err, port) {
    if (err) { throw err; }
    listen(port);
  });
}
else {
  listen(port);
}

function listen(port) {
  var options = {
    root: argv._[0],
    cache: argv.c,
    timeout: argv.t,
    showDir: argv.d,
    autoIndex: argv.i,
    gzip: argv.g || argv.gzip,
    brotli: argv.b || argv.brotli,
    robots: argv.r || argv.robots,
    ext: argv.e || argv.ext,
    logFn: logger.request,
    proxy: proxy,
    proxyOptions: proxyOptions,
    showDotfiles: argv.dotfiles,
    mimetypes: argv.mimetypes,
    username: argv.username || process.env.NODE_HTTP_SERVER_USERNAME,
    password: argv.password || process.env.NODE_HTTP_SERVER_PASSWORD
  };

  if (argv.cors) {
    options.cors = true;
    if (typeof argv.cors === 'string') {
      options.corsHeaders = argv.cors;
    }
  }

  if (proxy) {
    try {
      new url.URL(proxy)
    }
    catch (err) {
      logger.info(chalk.red('Error: Invalid proxy url'));
      process.exit(1);
    }
  }

  if (tls) {
    options.https = {
      cert: argv.C || argv.cert || 'cert.pem',
      key: argv.K || argv.key || 'key.pem',
      passphrase: sslPassphrase,
    };
    try {
      fs.lstatSync(options.https.cert);
    }
    catch (err) {
      logger.info(chalk.red('Error: Could not find certificate ' + options.https.cert));
      process.exit(1);
    }
    try {
      fs.lstatSync(options.https.key);
    }
    catch (err) {
      logger.info(chalk.red('Error: Could not find private key ' + options.https.key));
      process.exit(1);
    }
  }

  var server = httpServer.createServer(options);
  server.listen(port, host, function () {
    var protocol = tls ? 'https://' : 'http://';

    logger.info([
      chalk.yellow('Starting up http-server, serving '),
      chalk.cyan(server.root),
      tls ? (chalk.yellow(' through') + chalk.cyan(' https')) : ''
    ].join(''));

    logger.info([chalk.yellow('\nhttp-server version: '), chalk.cyan(require('../package.json').version)].join(''));

    logger.info([
      chalk.yellow('\nhttp-server settings: '),
      ([chalk.yellow('CORS: '), argv.cors ? chalk.cyan(argv.cors) : chalk.red('disabled')].join('')),
      ([chalk.yellow('Cache: '), argv.c ? (argv.c === '-1' ? chalk.red('disabled') : chalk.cyan(argv.c + ' seconds')) : chalk.cyan('3600 seconds')].join('')),
      ([chalk.yellow('Connection Timeout: '), argv.t === '0' ? chalk.red('disabled') : (argv.t ? chalk.cyan(argv.t + ' seconds') : chalk.cyan('120 seconds'))].join('')),
      ([chalk.yellow('Directory Listings: '), argv.d ? chalk.red('not visible') : chalk.cyan('visible')].join('')),
      ([chalk.yellow('AutoIndex: '), argv.i ? chalk.red('not visible') : chalk.cyan('visible')].join('')),
      ([chalk.yellow('Serve GZIP Files: '), argv.g || argv.gzip ? chalk.cyan('true') : chalk.red('false')].join('')),
      ([chalk.yellow('Serve Brotli Files: '), argv.b || argv.brotli ? chalk.cyan('true') : chalk.red('false')].join('')),
      ([chalk.yellow('Default File Extension: '), argv.e ? chalk.cyan(argv.e) : (argv.ext ? chalk.cyan(argv.ext) : chalk.red('none'))].join(''))
    ].join('\n'));

    logger.info(chalk.yellow('\nAvailable on:'));

    if (argv.a && host !== '0.0.0.0') {
      logger.info(`  ${protocol}${host}:${chalk.green(port.toString())}`);
    } else {
      Object.keys(ifaces).forEach(function (dev) {
        ifaces[dev].forEach(function (details) {
          if (details.family === 'IPv4') {
            logger.info(('  ' + protocol + details.address + ':' + chalk.green(port.toString())));
          }
        });
      });
    }

    if (typeof proxy === 'string') {
      if (proxyOptions) {
        logger.info('Unhandled requests will be served from: ' + proxy + '. Options: ' + JSON.stringify(proxyOptions));
      }
      else {
        logger.info('Unhandled requests will be served from: ' + proxy);
      }
    }

    logger.info('Hit CTRL-C to stop the server');
    if (argv.o) {
      const openHost = host === '0.0.0.0' ? '127.0.0.1' : host;
      let openUrl = `${protocol}${openHost}:${port}`;
      if (typeof argv.o === 'string') {
        openUrl += argv.o[0] === '/' ? argv.o : '/' + argv.o;
      }
      logger.info('Open: ' + openUrl);
      opener(openUrl);
    }

    // Spacing before logs
    if (!argv.s) logger.info();
  });
}

if (process.platform === 'win32') {
  require('readline').createInterface({
    input: process.stdin,
    output: process.stdout
  }).on('SIGINT', function () {
    process.emit('SIGINT');
  });
}

process.on('SIGINT', function () {
  logger.info(chalk.red('http-server stopped.'));
  process.exit();
});

process.on('SIGTERM', function () {
  logger.info(chalk.red('http-server stopped.'));
  process.exit();
});
