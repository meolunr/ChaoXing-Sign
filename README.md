## 超星学习通自动签到

- 支持学习通的所有签到

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;（普通签到、手势签到、拍照签到、位置签到、二维码签到）

- 支持自定义签到开始时间和刷新间隔

- 支持 Server 酱推送签到消息到微信

### 使用方法
在 [Releases](/releases "Releases") 页面下载适合您系统的二进制可执行文件。

新建 **profile.json** 文件并和上一步下载的程序文件放置在同一目录。

根据需要在 **profile.json** 文件中填入适当的字段。

运行程序。
> Linux 下如需要后台运行，可使用 nohup ./chaoxingsign >> chaoxingsign.log 2>&1 &

#### profile.json 字段说明
```
{
  // 除 username 和 password 必填外，其他均为可选字段

  "username": "user",         // 用户名
  "password": "passwd",       // 密码
  "interval": 60,             // 刷新间隔时间，默认 60，单位：秒
  "startTime": "07:00",       // 开始时间，每天这个时间会自动开始签到，24 小时制。不填即为一直处于运行状态
  "endTime": "19:00",         // 结束时间
  "serverChan": "url",        // Server 酱推送地址
  "excludeCourse": ["1","2"]  // 在此列表内的课程 ID 不会自动签到，课程 ID 可以在程序开始运行时获得
}
```

#### 自定义拍照签到图片
默认签到图片是一张 654x872 尺寸的纯黑色背景。  
如需自定义拍照签到图片，可将图片重命名为 **photo.jpg** 放置在和程序文件同一目录。

#### Server 酱是什么
[点这儿](http://sc.ftqq.com "点这儿")