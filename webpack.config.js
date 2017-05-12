// `CheckerPlugin` is optional. Use it if you want async error reporting.
// We need this plugin to detect a `--watch` mode. It may be removed later
// after https://github.com/webpack/webpack/issues/3460 will be resolved.
const webpack = require('webpack');
const { CheckerPlugin } = require('awesome-typescript-loader');
const path = require('path');

const PRODUCTION = 'production';

let config = {
    entry: {
        index: "./src/frontend/index.tsx"
    },
    output: {
        path: path.resolve(__dirname, 'views', 'dist'),
        filename: "[name].js"
    },

    // Currently we need to add '.ts' to the resolve.extensions array.
    resolve: {
        extensions: ['.ts', '.tsx', '.js', '.jsx']
    },

    // Source maps support ('inline-source-map' also works)
    //devtool: 'cheap-module-source-map',

    // Add the loader for .ts files.
    module: {
        loaders: [
            {
                test: /\.css$/,
                use: [
                    { loader: 'style-loader' },
                    { loader: 'css-loader',
                        options: {
                            minimize: process.env.NODE_ENV === PRODUCTION,
                            sourceMap: process.env.NODE_EN !== PRODUCTION
                        }
                    }
                ]
            },
            {
                test: /\.(png|jpg|gif|svg|eot|ttf|woff|woff2)$/,
                loader: 'url-loader',
                options: {
                    limit: 10000
                }
            },
            {
                test: /\.tsx?$/,
                loader: 'awesome-typescript-loader',
                exclude: /node_modules/,
            }
        ]
    },
    plugins: [
        new CheckerPlugin(),
        new webpack.optimize.CommonsChunkPlugin({
            name: 'vendor',
            minChunks: function (module) {
                // this assumes your vendor imports exist in the node_modules directory
                return module.context && module.context.indexOf('node_modules') !== -1;
            }
        })
    ]
};

if (process.env.NODE_ENV === PRODUCTION) {
    delete config.devtool;
    config.plugins.push(
        new webpack.optimize.UglifyJsPlugin()
    )
}

module.exports = config;
