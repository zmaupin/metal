const path = require('path');

module.exports = {
  entry: path.resolve('src', 'index.tsx'),
  mode: 'development',
  output: {
    filename: 'bundle.js',
    path: path.resolve(__dirname, 'dist')
  },
  // Enable sourcemaps for debugging webpack's output.
  devtool: 'source-map',
  devServer: {
   contentBase: path.join(__dirname, 'dist'),
   compress: true,
   port: 9000
 }
  resolve: {
    // Add '.ts' and '.tsx' as resolvable extensions.
    extensions: ['.ts', '.tsx', '.js', '.json']
  },
  module: {
    rules: [
      // All files with a '.ts' or '.tsx' extension will be handled by 'awesome-typescript-loader'.
      { test: /\.tsx?$/, loader: 'awesome-typescript-loader' },
      // All output '.js' files will have any sourcemaps re-processed by 'source-map-loader'.
      { enforce: 'pre', test: /\.js$/, loader: 'source-map-loader' },
      // Load CSS
      {
        test: /\.(css|scss)/,
        use: [
          { loader: "style-loader" },
          { loader: "css-loader" },
          { loader: "sass-loader" }
        ]
      },
      {
        test: /\.(woff|woff2|eot|ttf|otf)$/,
        loader: "file-loader"
      },
      { test: /\.tsx$/, enforce: 'pre', use: [ { loader: 'tslint-loader' } ] },
    ]
  },
};
