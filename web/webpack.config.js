function buildConfig(env) {
  return require(`./config/${env || 'dev'}.js`); // eslint-disable-line
}

module.exports = buildConfig;
