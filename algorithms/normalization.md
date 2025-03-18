# 数据归一化算法说明文档

## 1. 线性归一化方法

### 1.1 最小-最大归一化 (Min-Max Scaling)
**公式**  
$$ x' = \frac{x - \text{dataMin}}{\text{dataMax} - \text{dataMin}} $$

**参数说明**
- $\text{dataMin}$：数据最小值（通过 `slices.Min` 计算）
- $\text{dataMax}$：数据最大值（通过 `slices.Max` 计算）

**适用场景**
- 数据分布已知且无明显异常值
- 适配神经网络激活函数（如 Sigmoid/ReLU）

---

## 2. 标准化方法

### 2.1 Z-Score 标准化
**公式**  
$$ x' = \frac{x - \mu}{\sigma} $$

**参数说明**
- $\mu$：样本均值
- $\sigma$：样本标准差（分母为 $n-1$）

**特点**
- 处理后数据均值为 0，标准差为 1
- 受极端值影响可能较大

### 2.2 鲁棒标准化 (Robust Scaling)
**公式**  
$$ x' = \frac{x - \text{median}}{\text{IQR}} $$

**参数说明**
- $\text{median}$：数据中位数
- $\text{IQR}$：四分位距（Q3 - Q1）

---

## 3. 向量归一化方法

### 3.1 L2 正则化
**公式**  
$$ x' = \frac{x}{\|x\|_2} \quad \text{其中} \quad \|x\|_2 = \sqrt{\sum_{i=1}^n x_i^2} $$

**适用场景**
- 文本分类（如 TF-IDF 向量）
- 需要单位向量输入的模型

---

## 4. 非线性归一化方法

### 4.1 对数变换
**公式**  
$$ x' = \ln\left(\frac{x}{\text{dataMax}} + \epsilon\right) $$

**参数说明**
- $\epsilon = 1e-9$：防止 $\ln(0)$ 的极小值

### 4.2 分位数变换
**公式**  
$$
x' =
\begin{cases}
F^{-1}_U(\text{rank}(x)/n) & \text{均匀分布} \\
F^{-1}_N(\text{rank}(x)/n) & \text{正态分布}
\end{cases}
$$

**Probit 函数近似公式**  
$$
\text{probit}(p) =
\begin{cases}
t - \frac{2.515517 + 0.802853t + 0.010328t^2}{1 + 1.432788t + 0.189269t^2 + 0.001308t^3} & p > 0.5 \\
-(t - \frac{2.515517 + 0.802853t + 0.010328t^2}{1 + 1.432788t + 0.189269t^2 + 0.001308t^3}) & p \leq 0.5
\end{cases}
$$  
其中 $t = \sqrt{-2 \ln(\min(p, 1-p))}$

---

## 5. 选择建议
| 场景                  | 推荐方法                     |
|-----------------------|-----------------------------|
| 数据含异常值          | 鲁棒标准化、分位数变换       |
| 稀疏数据处理          | L2 正则化、绝对最大值归一化 |
| 深度学习输入          | 最小-最大归一化             |
| 正态分布假设          | Z-Score、分位数正态变换     |

---

## 6. 注意事项
1. **JetBrains 渲染设置**
    - 安装插件：`File → Settings → Plugins → Markdown`
    - 启用公式渲染：`Languages & Frameworks → Markdown → Enable LaTeX`
    - 推荐使用 **MathJax** 渲染引擎

2. **代码实现一致性**
    - 所有算法实现均包含零方差防御（返回原始数据或报错）
    - 大数据集建议预分配内存并行计算  