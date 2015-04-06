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


# driver

データ読み書きドライバ。


## 1. インターフェース


### 1.1. RawDataStore

文字列キー、バイト列の値。


### 1.2. KeyValueStore

文字列キー、指定されたデータ型の値。


### 1.3. ListedKeyValueStore

文字列キー、指定されたデータ型の値。
キー一覧メソッド持ち。


### 1.4. VolatileKeyValueStore

文字列キー、指定されたデータ型の値。
期限が来たら勝手に消える。


### 1.5. ConcurrentVolatileKeyValueStore

文字列キー、指定されたデータ型の値。
期限が来たら勝手に消える。
並列アクセス向けメソッド持ち。


## 2. 実装


### 2.1. Memory

メモリ上に保存する。主にデバッグ用。


### 2.2 File

ファイルシステムに保存する。


### 2.3 Mongo

Mongodb に保存する。


### 2.4 Redis

redis に保存する。
