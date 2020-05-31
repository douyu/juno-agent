![](doc/logo.png)

[![Build Status](https://travis-ci.org/douyu/juno.svg?branch=master)](https://travis-ci.org/douyu/juno)
[![codecov](https://codecov.io/gh/douyu/juno/branch/master/graph/badge.svg)](https://codecov.io/gh/douyu/juno)
[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/douyu/juno?tab=doc)
[![Go Report Card](https://goreportcard.com/badge/github.com/douyu/juno)](https://goreportcard.com/report/github.com/douyu/juno)
![license](https://img.shields.io/badge/license-Apache--2.0-green.svg)

# Juno-Agent
Juno-Agent是一个提供服务代理、应用配置下发、应用配置解析、shell沙箱、探活、消息总线的Agent。

Juno-Agent的设计目标主要是让开发能够通过可插拔的组件，观测和治理自己的系统。

## 最小依赖
* Linux kernel version 2.6.23 or later
* Windows 7 or later
* FreeBSD 11.2 or later
* MacOS 10.11 El Capitan or later

## 快速开始
查看帮助文档
```cmd
Juno-agent --help
```
生成默认配置
```cmd
Juno-agent config > Juno-agent.toml
```
使用文本配置，启用一个``test``组件
```cmd
Juno-agent --config=Juno-agent.toml --test
```
使用文本配置，启用全部组件
```cmd
Juno-agent --config=Juno-agent.toml
```

## 组件
* 代理模块
* 配置模块
* Shell沙箱
* 探活
* 消息总线
* HTTP模块
* 依赖探活(提供HTTP接口，针对应用依赖的组件进行探活)

## api
* [api文档](https://github.com/douyu/juno-agent/tree/master/doc/api/api.md)


## Contact

- DingTalk: 
    ![DingTalk](doc/dingtalk.png)
- Wechat:
    ![Wechat](doc/wechat.png)

