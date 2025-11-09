# ani-generator

複数の画像から Windows アニメーション カーソル（ANI）を生成するコマンドラインツールです。PNG/JPEG/GIF/BMP/WebP を読み込み、指定したホットスポットとシーケンス・表示時間で `.ani` を出力します。

## 特長
- **複数画像からANI生成**: 任意順序のフレームシーケンスを指定可能
- **透過対応**: 32bit RGBA をそのまま書き出し、アルファマスクも自動生成
- **複数フォーマット入力**: PNG/JPEG/GIF/BMP/WebP をサポート（`golang.org/x/image`）
- **JSON設定対応**: CLI 引数の代わりに設定ファイルで一括指定可能

## 制約（重要）
- 入力画像サイズは次のいずれかのみ対応: `32` / `64` / `128` / `256` ピクセル（幅・高さとも）
- フレームシーケンス（`frameIndex`）と表示時間（`rate`）は同じ長さである必要があります
- 画像枚数と各フレームのホットスポット数（`hotSpotX`/`hotSpotY`）は同数である必要があります
- `rate` の単位は「1/60 秒」です（例: `6` は約 100ms）

## ビルド
前提: Go 1.25 以上（`go.mod` 参照）

```bash
go build -o ani-generator
```

## 使い方
### 1) 直接オプションで指定
```bash
ani-generator ^
	--input img/F00.png ^
	--input img/F01.png ^
	--input img/F02.png ^
	--output img/export.ani ^
	--hotSpotX 16 --hotSpotX 16 --hotSpotX 16 ^
	--hotSpotY 16 --hotSpotY 16 --hotSpotY 16 ^
	--frameIndex 0 --frameIndex 1 --frameIndex 2 --frameIndex 1 ^
	--rate 6 --rate 6 --rate 6 --rate 6
```

### 2) JSON 設定ファイルで指定
`settings.json` の例:
```json
{
	"images": [
		{ "path": "img/F00.png", "hotSpotX": 16, "hotSpotY": 16 },
		{ "path": "img/F01.png", "hotSpotX": 16, "hotSpotY": 16 },
		{ "path": "img/F02.png", "hotSpotX": 16, "hotSpotY": 16 }
	],
	"frameIndexes": [0, 1, 2, 1],
	"rates": [6, 6, 6, 6],
	"output": "img/export.ani"
}
```
実行:
```bash
ani-generator --json settings.json
```

## オプション一覧
- **--input, -i（複数可）**: 入力画像ファイル（例: `--input img/F00.png`）
- **--output, -o**: 出力 ANI ファイル（例: `--output img/export.ani`）
- **--hotSpotX, -hx（複数可）**: 各フレームのホットスポット X 座標
- **--hotSpotY, -hy（複数可）**: 各フレームのホットスポット Y 座標
- **--frameIndex, -f（複数可）**: 表示シーケンスで使うフレーム番号（0 始まり）
- **--rate, -r（複数可）**: 各ステップの表示時間（単位: 1/60 秒）
- **--json**: 上記一式を JSON ファイルから読み込む

ヒント:
- フレームシーケンスの例: フレームが 3 枚で `0 → 1 → 2 → 1` の往復にしたい場合は、`--frameIndex 0 --frameIndex 1 --frameIndex 2 --frameIndex 1`
- 各表示時間をすべて 100ms にする場合は、`--rate 6` をシーケンス長分だけ指定します

## サンプル
リポジトリの `img/` ディレクトリに素材が入っています。例えば下記のように実行できます。
```bash
ani-generator ^
	--input img/F00.png --input img/F01.png --input img/F02.png ^
	--output img/export.ani ^
	--hotSpotX 16 --hotSpotX 16 --hotSpotX 16 ^
	--hotSpotY 16 --hotSpotY 16 --hotSpotY 16 ^
	--frameIndex 0 --frameIndex 1 --frameIndex 2 --frameIndex 1 ^
	--rate 6 --rate 6 --rate 6 --rate 6
```

## エラーになったとき
- 「unsupported image width/height」: 入力画像サイズが未対応（`32/64/128/256` 以外）です
- 「images/frameIndexes/rates が空」: 必須の配列が空です
- 「配列長の不一致」: 画像数とホットスポット数、またはフレームシーケンスと rate の数が一致していません
- 実行時にエラーが出た場合、簡易的な Usage が表示されます

## 実装メモ
- 出力は RIFF 構造で `anih`/`LIST(fram)`/`rate`/`seq ` チャンクを持ち、各フレームは `icon` チャンクとして格納されます
- `rate` は 1 要素 4byte の Little Endian、単位は 1/60 秒
- ホットスポットは ICO/CUR ヘッダの予約領域に書き込みます（カーソルで使用）
