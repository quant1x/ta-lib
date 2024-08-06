# Changelog
All notable changes to this project will be documented in this file.

## [Unreleased]

## [0.7.28] - 2024-08-06
### Changed
- 更新engine版本到1.8.45

## [0.7.27] - 2024-08-06
### Changed
- 更新engine版本到1.8.44
- 调整chart库的引用
- update changelog

## [0.7.26] - 2024-08-06
### Changed
- 更新engine版本到1.8.43
- update changelog

## [0.7.25] - 2024-07-31
### Changed
- 更新engine版本到1.8.41

## [0.7.24] - 2024-07-14
### Changed
- 更新engine版本到1.8.40
- update changelog

## [0.7.23] - 2024-07-05
### Changed
- 新增pzzy指标
- 新增macd用series计算的函数
- 更新依赖库版本
- 调整部分测试代码
- update changelog

## [0.7.22] - 2024-06-27
### Changed
- 更新engine版本到1.8.38
- update changelog

## [0.7.21] - 2024-06-27
### Changed
- 更新依赖库版本
- update changelog

## [0.7.20] - 2024-06-27
### Changed
- 更新engine版本到1.8.37
- update changelog

## [0.7.19] - 2024-06-26
### Changed
- 更新engine版本到1.8.35
- 更新engine版本到1.8.36
- update changelog

## [0.7.18] - 2024-06-26
### Changed
- 修复与engine循环引入的问题
- update changelog

## [0.7.17] - 2024-06-26
### Changed
- 更新engine版本到1.8.34
- update changelog

## [0.7.16] - 2024-06-24
### Changed
- 更新engine版本到1.8.30
- update changelog

## [0.7.15] - 2024-06-21
### Changed
- 删除测试数据
- update changelog

## [0.7.14] - 2024-06-21
### Changed
- 更新engine版本到1.8.29
- update changelog

## [0.7.13] - 2024-06-20
### Changed
- 更新engine版本到1.8.28
- update changelog

## [0.7.12] - 2024-06-20
### Changed
- 更新engine版本到1.8.27
- update changelog

## [0.7.11] - 2024-06-16
### Changed
- 新增SAR指标计算方法
- 更新依赖库engine版本到1.8.22
- sar算法增加注释
- sar算法增加增量计算方法
- 结构体增加说明文字
- 调整sar实现的部分函数名
- sar特征组合结构增加当前趋势的周期数, 上涨趋势, 周期数大于0, 下跌趋势, 周期数小于0, 绝对值就是已过多少天
- sar特征组合结构周期数字段Period增加用途描述
- 更新依赖库engine版本到1.8.23
- 调整sar测试代码
- 更新依赖库版本
- 参数输出保留2位小数点
- 调整sar代码, 剔除废弃的代码
- update changelog

## [0.7.10] - 2024-06-02
### Changed
- 更新依赖库engine版本到1.8.20
- update changelog

## [0.7.9] - 2024-05-27
### Changed
- 更新engine版本到1.8.18
- update changelog
- update changelog
- update changelog

## [0.7.8] - 2024-05-27
### Changed
- 更新依赖库engine版本到1.8.17
- update changelog

## [0.7.7] - 2024-05-24
### Changed
- embed资源路径是*nix格式, 否则在windows无法打开
- update changelog

## [0.7.6] - 2024-05-21
### Changed
- 更新engine版本到1.8.13
- update changelog

## [0.7.5] - 2024-05-20
### Changed
- 更新engine版本到1.8.12
- update changelog

## [0.7.4] - 2024-05-20
### Changed
- 更新engine版本到1.8.11
- update changelog

## [0.7.3] - 2024-05-18
### Changed
- 更新engine版本到1.8.10
- update changelog

## [0.7.2] - 2024-05-15
### Changed
- 更新engine版本到1.8.9
- update changelog

## [0.7.1] - 2024-05-15
### Changed
- 更新engine版本到1.8.8
- update changelog

## [0.7.0] - 2024-05-13
### Changed
- 更新engine版本到1.8.7
- update changelog

## [0.6.9] - 2024-05-11
### Changed
- 更新engine版本到1.8.6
- update changelog

## [0.6.8] - 2024-05-02
### Changed
- 更新engine版本到1.8.4
- update changelog

## [0.6.7] - 2024-05-02
### Changed
- 更新engine版本到1.8.2
- 更新engine版本到1.8.3
- update changelog

## [0.6.6] - 2024-04-28
### Changed
- 拟增加新的波浪推导方式
- 更新engine版本到1.8.0
- 调整部分测试代码
- update changelog

## [0.6.5] - 2024-04-21
### Changed
- 新增v3版本的波浪处理逻辑
- 更新engine版本到1.7.9
- 拆分不同的wave用法
- update changelog

## [0.6.4] - 2024-04-19
### Changed
- 更新依赖库版本
- update changelog

## [0.6.3] - 2024-04-14
### Changed
- 更新依赖库engine版本到1.7.7
- 新增wedge趋势线交叉cross的方法
- update changelog

## [0.6.2] - 2024-04-12
### Changed
- 更新依赖库engine版本到1.7.6
- update changelog

## [0.6.1] - 2024-04-12
### Changed
- 补充K线样本结构的注释
- 优化chart
- 收敛chart引用
- 形态Pattern接口增加输出图表的方法
- 补充方法注释
- 调整waves初始化方法
- 调整waves用法
- 更新依赖库版本
- 更新依赖库engine版本到1.7.5
- update changelog

## [0.6.0] - 2024-04-04
### Changed
- 更新依赖库版本
- 新增三角形算法
- 根据配置项确定波峰波谷的取值
- waves新增Len方法
- charts新增最后数据的标签展示函数
- 调整测试代码的k线日期
- update changelog

## [0.5.9] - 2024-04-03
### Changed
- 新增数据样本结构
- 收敛图表功能
- 调整图表功能
- 调整图表功能
- 修订测试代码
- 调整结构体方法的接收器
- 调整楔形算法
- update changelog

## [0.5.8] - 2024-04-02
### Changed
- 更新依赖库版本
- 新增楔形趋势算法
- 趋势用当前值, 波峰用最高, 波谷用最低
- update changelog

## [0.5.7] - 2024-03-30
### Changed
- 更新依赖库engine的版本到1.7.0
- 修改通达信公式的源文件扩展名为tdx,支持语法高亮显示
- 更新依赖库版本
- 更新依赖库版本
- 更新依赖库版本
- update changelog

## [0.5.6] - 2024-03-21
### Changed
- 更新依赖库engine的版本到1.6.8
- update changelog

## [0.5.5] - 2024-03-19
### Changed
- 更新依赖库engine的版本到1.6.7
- update changelog

## [0.5.4] - 2024-03-19
### Changed
- 更新依赖库engine的版本到1.6.6
- update changelog

## [0.5.3] - 2024-03-18
### Changed
- 更新依赖库engine的版本到1.6.5
- update changelog

## [0.5.2] - 2024-03-17
### Changed
- 更新依赖库engine的版本到1.6.4
- update changelog

## [0.5.1] - 2024-03-17
### Changed
- 更新依赖库
- update changelog

## [0.5.0] - 2024-03-16
### Changed
- 优化部分代码,删除对gonum.org/v1/plot的依赖
- update changelog

## [0.4.9] - 2024-03-12
### Changed
- 更新依赖库版本及go版本
- update changelog

## [0.4.8] - 2024-03-12
### Changed
- 更新依赖库版本及go版本
- update changelog

## [0.4.7] - 2024-03-12
### Changed
- 更新依赖库版本
- update changelog

## [0.4.6] - 2024-03-11
### Changed
- 更新依赖库版本
- update changelog

## [0.4.5] - 2024-03-11
### Changed
- 更新依赖库engine版本
- update changelog

## [0.4.4] - 2024-03-11
### Changed
- 补充注释
- 更新依赖库num版本
- update changelog

## [0.4.3] - 2024-03-10
### Changed
- 更新依赖库num版本
- 调整颜色
- 新增字体默认值函数, 不返回错误
- update changelog

## [0.4.2] - 2024-03-03
### Changed
- 补充ta-lib基本信息
- 修复测试代码没有适配pandas的问题
- 抽象部分常用的go-chart用法
- 实验代码,新增html样式的K线图
- 新增红色定义, go-chart颜色中红色不准确的
- update changelog

## [0.4.1] - 2024-02-26
### Changed
- 更新依赖库版本
- update changelog

## [0.4.0] - 2024-02-25
### Changed
- 新增点的style
- 图表新增右侧Y轴的提示
- 新增打开默认浏览器的函数
- 更新依赖库版本
- 调整测试代码
- update changelog

## [0.3.9] - 2024-02-21
### Changed
- 增加变形W底的计算测试代码
- update changelog

## [0.3.8] - 2024-02-19
### Changed
- 适配新版本pandas
- update changelog

## [0.3.7] - 2024-02-19
### Changed
- 调整series函数
- update changelog

## [0.3.6] - 2024-02-18
### Changed
- 更新pandas版本
- update changelog

## [0.3.5] - 2024-02-12
### Changed
- 更新依赖库engine版本
- update changelog

## [0.3.4] - 2024-02-12
### Changed
- 新增默认的字体SimHei
- 新增series整型索引函数
- 调试绘图用法
- 测试PeekDetect功能
- 测试PeekDetect功能
- 新增FindPeeks查找波峰波谷功能
- 调整m头和w底测试代码
- 调整波峰波谷结构体名
- 更新依赖库版本
- 更新依赖库版本
- update changelog

## [0.3.3] - 2024-01-28
### Changed
- 更新pandas版本
- update changelog

## [0.3.2] - 2024-01-13
### Changed
- 适配新版本的engine
- update changelog

## [0.3.1] - 2023-12-24
### Changed
- 适配新版本的engine
- update changelog

## [0.3.0] - 2023-12-20
### Changed
- 适配新版本的engine
- update changelog

## [0.2.1] - 2023-12-19
### Changed
- 更新依赖库版本
- update changelog

## [0.2.0] - 2023-12-17
### Changed
- 更新依赖库版本
- update changelog

## [0.1.9] - 2023-12-17
### Changed
- 更新依赖库版本
- 新增试验性质的波峰检测工具
- update changelog

## [0.1.8] - 2023-12-12
### Changed
- 更新依赖库版本
- update changelog

## [0.1.7] - 2023-12-05
### Changed
- 更新依赖库版本
- update changelog

## [0.1.6] - 2023-11-06
### Changed
- 更新依赖库版本
- update changelog

## [0.1.5] - 2023-10-30
### Changed
- 测试代码的K线数据改成由gotdx获取
- update changelog

## [0.1.4] - 2023-10-29
### Changed
- 更新pandas版本
- update changelog

## [0.1.3] - 2023-10-08
### Changed
- 更新pandas版本
- update changelog

## [0.1.2] - 2023-10-08
### Changed
- 更新gox版本
- update changelog

## [0.1.1] - 2023-10-08
### Changed
- 增加89K指标源代码
- update changelog

## [0.1.0] - 2023-09-15
### Changed
- add LICENSE
- update changelog

## [0.0.3] - 2023-09-15
### Changed
- 更新依赖库版本
- 增加测试数据
- 调整测试代码
- update changelog

## [0.0.2] - 2023-09-13
### Changed
- 剔除测试文件
- add Files
- update changelog

## [0.0.1] - 2023-09-12
### Changed
- 第一次提交


[Unreleased]: https://gitee.com/quant1x/ta-lib.git/compare/v0.7.28...HEAD
[0.7.28]: https://gitee.com/quant1x/ta-lib.git/compare/v0.7.27...v0.7.28
[0.7.27]: https://gitee.com/quant1x/ta-lib.git/compare/v0.7.26...v0.7.27
[0.7.26]: https://gitee.com/quant1x/ta-lib.git/compare/v0.7.25...v0.7.26
[0.7.25]: https://gitee.com/quant1x/ta-lib.git/compare/v0.7.24...v0.7.25
[0.7.24]: https://gitee.com/quant1x/ta-lib.git/compare/v0.7.23...v0.7.24
[0.7.23]: https://gitee.com/quant1x/ta-lib.git/compare/v0.7.22...v0.7.23
[0.7.22]: https://gitee.com/quant1x/ta-lib.git/compare/v0.7.21...v0.7.22
[0.7.21]: https://gitee.com/quant1x/ta-lib.git/compare/v0.7.20...v0.7.21
[0.7.20]: https://gitee.com/quant1x/ta-lib.git/compare/v0.7.19...v0.7.20
[0.7.19]: https://gitee.com/quant1x/ta-lib.git/compare/v0.7.18...v0.7.19
[0.7.18]: https://gitee.com/quant1x/ta-lib.git/compare/v0.7.17...v0.7.18
[0.7.17]: https://gitee.com/quant1x/ta-lib.git/compare/v0.7.16...v0.7.17
[0.7.16]: https://gitee.com/quant1x/ta-lib.git/compare/v0.7.15...v0.7.16
[0.7.15]: https://gitee.com/quant1x/ta-lib.git/compare/v0.7.14...v0.7.15
[0.7.14]: https://gitee.com/quant1x/ta-lib.git/compare/v0.7.13...v0.7.14
[0.7.13]: https://gitee.com/quant1x/ta-lib.git/compare/v0.7.12...v0.7.13
[0.7.12]: https://gitee.com/quant1x/ta-lib.git/compare/v0.7.11...v0.7.12
[0.7.11]: https://gitee.com/quant1x/ta-lib.git/compare/v0.7.10...v0.7.11
[0.7.10]: https://gitee.com/quant1x/ta-lib.git/compare/v0.7.9...v0.7.10
[0.7.9]: https://gitee.com/quant1x/ta-lib.git/compare/v0.7.8...v0.7.9
[0.7.8]: https://gitee.com/quant1x/ta-lib.git/compare/v0.7.7...v0.7.8
[0.7.7]: https://gitee.com/quant1x/ta-lib.git/compare/v0.7.6...v0.7.7
[0.7.6]: https://gitee.com/quant1x/ta-lib.git/compare/v0.7.5...v0.7.6
[0.7.5]: https://gitee.com/quant1x/ta-lib.git/compare/v0.7.4...v0.7.5
[0.7.4]: https://gitee.com/quant1x/ta-lib.git/compare/v0.7.3...v0.7.4
[0.7.3]: https://gitee.com/quant1x/ta-lib.git/compare/v0.7.2...v0.7.3
[0.7.2]: https://gitee.com/quant1x/ta-lib.git/compare/v0.7.1...v0.7.2
[0.7.1]: https://gitee.com/quant1x/ta-lib.git/compare/v0.7.0...v0.7.1
[0.7.0]: https://gitee.com/quant1x/ta-lib.git/compare/v0.6.9...v0.7.0
[0.6.9]: https://gitee.com/quant1x/ta-lib.git/compare/v0.6.8...v0.6.9
[0.6.8]: https://gitee.com/quant1x/ta-lib.git/compare/v0.6.7...v0.6.8
[0.6.7]: https://gitee.com/quant1x/ta-lib.git/compare/v0.6.6...v0.6.7
[0.6.6]: https://gitee.com/quant1x/ta-lib.git/compare/v0.6.5...v0.6.6
[0.6.5]: https://gitee.com/quant1x/ta-lib.git/compare/v0.6.4...v0.6.5
[0.6.4]: https://gitee.com/quant1x/ta-lib.git/compare/v0.6.3...v0.6.4
[0.6.3]: https://gitee.com/quant1x/ta-lib.git/compare/v0.6.2...v0.6.3
[0.6.2]: https://gitee.com/quant1x/ta-lib.git/compare/v0.6.1...v0.6.2
[0.6.1]: https://gitee.com/quant1x/ta-lib.git/compare/v0.6.0...v0.6.1
[0.6.0]: https://gitee.com/quant1x/ta-lib.git/compare/v0.5.9...v0.6.0
[0.5.9]: https://gitee.com/quant1x/ta-lib.git/compare/v0.5.8...v0.5.9
[0.5.8]: https://gitee.com/quant1x/ta-lib.git/compare/v0.5.7...v0.5.8
[0.5.7]: https://gitee.com/quant1x/ta-lib.git/compare/v0.5.6...v0.5.7
[0.5.6]: https://gitee.com/quant1x/ta-lib.git/compare/v0.5.5...v0.5.6
[0.5.5]: https://gitee.com/quant1x/ta-lib.git/compare/v0.5.4...v0.5.5
[0.5.4]: https://gitee.com/quant1x/ta-lib.git/compare/v0.5.3...v0.5.4
[0.5.3]: https://gitee.com/quant1x/ta-lib.git/compare/v0.5.2...v0.5.3
[0.5.2]: https://gitee.com/quant1x/ta-lib.git/compare/v0.5.1...v0.5.2
[0.5.1]: https://gitee.com/quant1x/ta-lib.git/compare/v0.5.0...v0.5.1
[0.5.0]: https://gitee.com/quant1x/ta-lib.git/compare/v0.4.9...v0.5.0
[0.4.9]: https://gitee.com/quant1x/ta-lib.git/compare/v0.4.8...v0.4.9
[0.4.8]: https://gitee.com/quant1x/ta-lib.git/compare/v0.4.7...v0.4.8
[0.4.7]: https://gitee.com/quant1x/ta-lib.git/compare/v0.4.6...v0.4.7
[0.4.6]: https://gitee.com/quant1x/ta-lib.git/compare/v0.4.5...v0.4.6
[0.4.5]: https://gitee.com/quant1x/ta-lib.git/compare/v0.4.4...v0.4.5
[0.4.4]: https://gitee.com/quant1x/ta-lib.git/compare/v0.4.3...v0.4.4
[0.4.3]: https://gitee.com/quant1x/ta-lib.git/compare/v0.4.2...v0.4.3
[0.4.2]: https://gitee.com/quant1x/ta-lib.git/compare/v0.4.1...v0.4.2
[0.4.1]: https://gitee.com/quant1x/ta-lib.git/compare/v0.4.0...v0.4.1
[0.4.0]: https://gitee.com/quant1x/ta-lib.git/compare/v0.3.9...v0.4.0
[0.3.9]: https://gitee.com/quant1x/ta-lib.git/compare/v0.3.8...v0.3.9
[0.3.8]: https://gitee.com/quant1x/ta-lib.git/compare/v0.3.7...v0.3.8
[0.3.7]: https://gitee.com/quant1x/ta-lib.git/compare/v0.3.6...v0.3.7
[0.3.6]: https://gitee.com/quant1x/ta-lib.git/compare/v0.3.5...v0.3.6
[0.3.5]: https://gitee.com/quant1x/ta-lib.git/compare/v0.3.4...v0.3.5
[0.3.4]: https://gitee.com/quant1x/ta-lib.git/compare/v0.3.3...v0.3.4
[0.3.3]: https://gitee.com/quant1x/ta-lib.git/compare/v0.3.2...v0.3.3
[0.3.2]: https://gitee.com/quant1x/ta-lib.git/compare/v0.3.1...v0.3.2
[0.3.1]: https://gitee.com/quant1x/ta-lib.git/compare/v0.3.0...v0.3.1
[0.3.0]: https://gitee.com/quant1x/ta-lib.git/compare/v0.2.1...v0.3.0
[0.2.1]: https://gitee.com/quant1x/ta-lib.git/compare/v0.2.0...v0.2.1
[0.2.0]: https://gitee.com/quant1x/ta-lib.git/compare/v0.1.9...v0.2.0
[0.1.9]: https://gitee.com/quant1x/ta-lib.git/compare/v0.1.8...v0.1.9
[0.1.8]: https://gitee.com/quant1x/ta-lib.git/compare/v0.1.7...v0.1.8
[0.1.7]: https://gitee.com/quant1x/ta-lib.git/compare/v0.1.6...v0.1.7
[0.1.6]: https://gitee.com/quant1x/ta-lib.git/compare/v0.1.5...v0.1.6
[0.1.5]: https://gitee.com/quant1x/ta-lib.git/compare/v0.1.4...v0.1.5
[0.1.4]: https://gitee.com/quant1x/ta-lib.git/compare/v0.1.3...v0.1.4
[0.1.3]: https://gitee.com/quant1x/ta-lib.git/compare/v0.1.2...v0.1.3
[0.1.2]: https://gitee.com/quant1x/ta-lib.git/compare/v0.1.1...v0.1.2
[0.1.1]: https://gitee.com/quant1x/ta-lib.git/compare/v0.1.0...v0.1.1
[0.1.0]: https://gitee.com/quant1x/ta-lib.git/compare/v0.0.3...v0.1.0
[0.0.3]: https://gitee.com/quant1x/ta-lib.git/compare/v0.0.2...v0.0.3
[0.0.2]: https://gitee.com/quant1x/ta-lib.git/compare/v0.0.1...v0.0.2

[0.0.1]: https://gitee.com/quant1x/ta-lib.git/releases/tag/v0.0.1
