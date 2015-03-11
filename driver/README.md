<!--
Copyright 2015 realglobe, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
-->


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
