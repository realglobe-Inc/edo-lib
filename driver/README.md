データ読み書き補助
===


RawDataStore
---

文字列キーで、バイト列の値を管理。


KeyValueStore
---

文字列キーで、指定されたデータ型の値を管理。
RawDataStore を包んで実装しているものが多い。


TimeLimitedKeyValueStore
---

文字列キーで、指定されたデータ型の値を管理。
有効期限が来たら勝手に消える。


Memory*
---

メモリ上に保存する。主にデバッグ用。


File*
---

ファイルシステムに保存する。


Mongo*
---

Mongodb に保存する。


Web*
---

HTTP(S) でリモートに保存する。
