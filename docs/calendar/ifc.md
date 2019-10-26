Bus 类型：session bus

服务名称：com.deepin.daemon.Calendar

对象路径：/com/deepin/daemon/Calendar/Scheduler

接口名称：com.deepin.daemon.Calendar.Scheduler


# 方法

## 获取指定范围内的日程
GetJobs(startYear int32, startMonth int32, startDay int32,
endYear int32, endMonth int32, endDay int32) -> (string)


指定开始日期和结束日期。


返回 JSON 格式

```json
[
 {
    "Date": "2019-01-01",
    "Jobs": [ job1, job2, ... ],
 }, ...
]
```

## 查询日程
QueryJobs(params string) -> (string)

params 为 JSON 格式：
```json
{
  "Key": "关键字",
  "Start": "2019-09-27T17:00:00+08:00",
  "End": "2019-09-27T18:00:00+08:00"
}
```

params 各字段用途：
Key 是关键字，用于看 Job 的 Title 字段值中是否有此字符串，会忽略两头的空白，如果为空，表示不使用关键字过滤条件。
Start 是查询时间段的开始时间，格式为 RFC3339，比如"2006-01-02T15:04:05+07:00"。
End 是查询时间段的结束时间，格式为 RFC3339，比如"2006-01-02T15:04:05+07:00"。

返回数据格式同 GetJobs。

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

# 信号

JobsUpdated(ids []int64)

只在后端自己修改了 job 数据后发送，前端收到信号后，不用使用 ids 参数，刷新界面所需数据。
