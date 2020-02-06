var path = require('path');
var webpack = require('webpack');

module.exports = {
  mode: 'development',
  entry: './client.js',
  output: {
    path: path.resolve(__dirname, 'dist'),
    filename: 'main.js'
  },
  plugins: [
    new webpack.EnvironmentPlugin(['PUBLISH_KEY', 'SUBSCRIBE_KEY', 'CHANNEL']),
  ]
};