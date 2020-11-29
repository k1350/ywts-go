# 使い方
**(1) Docker compose と yarn が使用可能な環境を用意する**  
説明は省略

**(2) ソースコードをcloneする**  
```
git clone https://github.com/k1350/ywts-go.git
```

**(3) Docker-compose用のenvファイルを作成する。**  
プロジェクト直下の .env.sample を .env にリネームする。  
本件で使うMySQLのルートパスワード、ユーザー名・パスワードを記入する。  
またNginxのタイムゾーンを変更したい場合は変更する。
```
MYSQL_ROOT_PASSWORD=db_root
MYSQL_USER=db_user
MYSQL_PASSWORD=db_pass
NGINX_TZ=Asia/Tokyo
```

**(4) FirebaseUIの使用準備を行う**  
Firebaseプロジェクトを新規作成する。  
次に Authentication を有効にし、Sign-in method で「メール/パスワード」を有効にする。

次にウェブ用のアプリを追加する。

アプリ追加後にコンソールに戻り、プロジェクトの設定から「サービスアカウント」を選択する。  
「新しい秘密鍵の生成」というボタンを押下し、秘密鍵をダウンロードする。  
ダウンロードした秘密鍵を app/private に配置する。

同じ設定画面から「全般」を選択し、ページ下部の「マイアプリ」から先ほど作成したウェブ用アプリの Firebase SDK snippet を表示しておく。

**(5) バックエンドのenvファイルを作成する**  
app/.env.sample を app/.env にリネームする。

MySQLのユーザ名・パスワードは (3) と同じ内容を設定する。  
FIREBASE_JSON_PATH は、パス末尾のjsonファイル名を (4) で app/private に配置した秘密鍵の名称で書き換える。  

SESSION_KEY は適当な推測困難な英数字文字列を記入する。

**(6) フロントエンドのenvファイルを作成する**  
assets/.env.sample を app/.env にリネームする。

VUE_APP_API_BASE_URL は、このREADMEにしたがってDockerで動かすなら変更不要。

VUE_APP_FIREBASE_xxx を (4) で開いておいた Firebase SDK snippet の中身で書き換える。（なお本README記載時のFirebaseバージョンは8.1.1なので、今後バージョンアップに伴いsnippetの内容が変更されたら動かない可能性あり）

**(7) フロントエンド用の依存パッケージをダウンロードする**  
assets 直下で

```
yarn install
```

※本当はDocker composeしたら自動で yarn install するようにしたかったのですが、開発中の時点ではWindows用Dockerの仕様によって yarn install が途中で固まる現象が発生したので事前に yarn install しておく方式にしました。WSL2を有効にしたら問題が解消するかもしれないです。  
もし問題が解消したら、docker-compose.yml の

> command: sh -c "cd assets && yarn serve"

という部分を

> command: sh -c "cd assets && yarn install && yarn serve"

にすれば事前installは不要になります。

**(8) （初回実行時のみ）docker-compose build**  
プロジェクト直下で

```
docker-compose build
```

**(9) 起動する**  
プロジェクト直下で

```
docker-compose up
```

初回は ywts-db が初回設定するので時間がかかります。  
また ywts-go と ywts-vue が起動するのに毎回しばらくかかるので待ってください。

全部起動したら http://localhost を開いてください。ログイン画面が表示されれば成功です。

アカウント登録は、ログイン画面で新しいメールアドレスを入力すれば流れで登録できるようになっています。

# データの保存場所
認証情報は Firebase プロジェクトに保存されます。

その他のデータはDocker の Volume に入っています。

```
docker volume ls
```

コマンドでボリュームの一覧を確認できます。ywts-go_db_data というのがあると思います。

ボリュームのバックアップや削除に関してはDockerのヘルプに書いてあると思うので省略します。
