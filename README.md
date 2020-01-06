### 已支持变量
变量名 | 结果 | 描述
--- | --- | ---
success | 1 | 该变量永远返回1
rand  | 随机值 | 获取[1-100]内的随机值
ip | 当前用户IP | 需要先将用户IP注入的Context中
country | 获取注入IP所在国家 | 需要先将用户IP注入的Context中
province | 获取注入IP所在省份 | 需要先将用户IP注入的Context中
city | 获取注入IP所在城市 | 需要先将用户IP注入的Context中
timestamp | 获取当前时间时间戳（单位秒）| number
ts_simple | 获取当前时间的时间戳int64类型| number
second | 获取当前时间的秒 | number
minute | 获取当前时间的分钟 | number
hour | 获取当前时间的小时 | number
day | 获取当前时间在当年中的第几天 | number
month | 获取当前时间的月份 | number
year | 获取当前时间的年份 | number
wday | 获取当时时间为本周的第几天 | 周日为0，周六为6
date | 获取当前时间的的日期（具体到天） | string(2006-01-02)
time | 获取当前时间 | string(2006-01-02 15:04:05)
ua | 获取user_agent信息 | 需要先注入的Context中
referer| 获取refer信息 | 需要先注入的Context中
is_login | 获取当前用户是否登录 | 需要先蒋user_id注入的Context中
version | 获取版本信息 | 需要先注入的Context中
platform | 获取平台信息 | 需要先注入的Context中
channel | 获取渠道信息 | 需要先注入的Context中
uid | 获取用户ID | 需要先注入的Context中
device | 获取设备信息| 需要先注入的Context中
user_tag | 获取用户标签 | 需要先注入的Context中
get.xxx | 获取注入到context中form信息 | 需要先注入的Context中
data.xxx | 获取传入的对象的字段值 | 传入的对象可以是ptr或者struct或者是map
calc.experation | 获取计算表达式 | 计算experation 表达式例如（calc.__value1 * __value2）需要业务方实现CalcFactorGet接口获取变量的值，需要返回float64类型的值
freq.xxx | 获取xxx对应的频次 | 用户频次控制；需要业务方实现FrequencyGetter接口
ctx.xxx | 获取对应Context中的值 | 需要先注入的Context中


### 已支持运算符
 运算符 | 描述 | 逻辑
 :---:  | 			  :---:  | :---:
   =  | 比较运算符 		| 比较字符串或者数字是否相等
   != | 比较运算符 		| 与 = 逻辑互斥
   <> | 比较运算符 		| 与 != 等价
   >  | 比较运算符 		| 比较字符串或者数字是否大于另外一个值
   >= | 比较运算符 		| 比较字符串或者数字是否大于或者等于另外一个值
   <  | 比较运算符 		| 比较字符串或者数字是否小于另外一个值
   <= | 比较运算符 		| 比较字符串或者数字是否小于或者等于另外一个值
   ~  | 正则或者字符串匹配| 以/开头且以/结尾的字符串表示正则；否则会判断目标字符串是否包含当前字符串（内部会转化为小写之后在判断）
   !~ | 正则或者字符串匹配 | 与 ~ 逻辑相反
   ~* | 正则或者字符串匹配 | 匹配目标对象中的
   !~* | 正则或者字符串匹配 | 
   between | 区间变量运算符 |
   in | 判断数组中是否存在某个值 |
   nin | 判断数组中是否不存在某个值| 与 in 逻辑相反（not in）
   any | 数组中的任意一个值匹配即可 | 
   has | 数组中必须存在某一个值 |
   none | 数组中不能存在某个值 |
   vgt | 比较版本号是否大于 | 版本值为 xx.xx.xx
   vgte | 比较版本号是否大于或者等于 | 版本值为 xx.xx.xx
   vlt | 比较本班好是否小于 | 版本值为 xx.xx.xx
   vlte | 比较版本号是否小于或者等于 | 版本值为 xx.xx.xx
   iir | 判断IP是否在某个IP段中 | in ip range
   niir | 判断IP是否不在某个IP段中 | not in ip range 与iir 逻辑相反