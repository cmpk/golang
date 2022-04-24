const path = require("path");
const webpack = require("webpack");
const dotenv = require("dotenv");

module.exports = () => {
  const env = dotenv.config().parsed;
  console.log(JSON.stringify(env));

  return {
    mode: "development",
    devtool: "inline-source-map", //エラー発生時に js を纏める前のソースコードをポイントするため
    entry: "./src/index.jsx",
    output: {
      filename: "main.js",
      path: path.resolve(__dirname, "dist"), // js を纏めたファイルはここに出力する
    },
    module: {
      rules: [
        {
          // Babel のローダーの設定
          test: /\.(js|mjs|jsx)$/, //対象のファイルの拡張子
          exclude: /node_modules/,
          use: [
            {
              loader: "babel-loader",
              options: {
                presets: ["@babel/preset-env", "@babel/preset-react"],
              },
            },
          ],
        },
      ],
    },
    //webpack-dev-server の設定
    devServer: {
      liveReload: true,
      watchFiles: ["src/*"],

      server: "http",
      allowedHosts: [env.APP_DOMAIN_NAME],
      port: 3000,
      static: "./dist",

      setupExitSignals: true,

      client: {
        logging: "log",
        overlay: true,
        progress: true,

        webSocketURL: "wss://0.0.0.0/ws",
      },

      historyApiFallback: true,
      //バンドルされたファイルを出力する（実際に書き出す）には以下のコメントアウトを外す
      //writeToDisk: true
    },

    plugins: [
      //jsx ファイル内の環境変数の置換
      new webpack.DefinePlugin({
        "process.env": JSON.stringify(env),
      }),
    ],
  };
};
