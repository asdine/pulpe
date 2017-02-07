function buildConfig(env) {
  return require(`./config/${env}.js`); // eslint-disable-line
}

module.exports = buildConfig;
