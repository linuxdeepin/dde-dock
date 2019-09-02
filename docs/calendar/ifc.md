Bus 类型：session bus

服务名称：com.deepin.daemon.Calendar

对象路径：/com/deepin/daemon/Calendar/Scheduler

接口名称：com.deepin.daemon.Calendar.Scheduler



## 获取指定范围内的日程
GetJobs(startYear int32, startMonth int32, startDay int32,
endYear int32, endMonth int32, endDay int32) -> (string)


指定开始日期和结束日期。

返回数据格式同 GetJobs。

返回 JSON 格式

```json
[
 {
    "Date": "2019-01-01",
    "Jobs": [ job1, job2, ... ],
 }, ...
]
```

## 根据 id 获取日程
GetJob(jobId int64) -> (string)

根据 id 获取相应的 job。

返回 job 的 json 字符串表示。

## 创建日程

CreateJob(jobInfo string) -> (id int64)

jobInfo 为 job 的字符串表示。

返回新 job 的 id。


## 更新日程

UpdateJob(jobInfo string)

jobInfo 为 job 的字符串表示

## 删除日程

DeleteJob(id int64)

根据 id 删除相应的 job。

## 获取所有类型

GetTypes() -> (string)

返回 job type 列表的 JSON 表示。

## 根据 ID 获取类型

GetType(id int64) -> (string)

返回 job type 的 JSON 表示。

job type 具有的字段：

- ID int
- Name string
- Color string

Name 名称，数据类型：字符串，不能为空。

Color 颜色值，数据类型：字符串，不能为空，为 ”#“ 开头的十六进制颜色。

## 创建类型

CreateType(typeInfo string) -> (id int64)

参数 typeInfo 为 job type 的 JSON 表示。

返回新创建的 job type 的 id。

## 删除类型

DeleteType(id int64) -> ()

根据 id 删除相应的 job type。

## 更新类型

UpdateType(typeInfo string) -> ()

参数 typeInfo 为 job type 的 JSON 表示。
