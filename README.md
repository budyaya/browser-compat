# 说明

将 browser compat data 数据导入 sqlite 数据库，方便查询。

## 使用

```sql
SELECT * FROM browser_compat_data WHERE browser='firefox' and browser_version like '130%'
```