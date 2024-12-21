# 说明

将 browser compat data 数据导入 sqlite 数据库，方便查询。

```shell
git clone https://github.com/mdn/browser-compat-data.git
git clone https://github.com/budyaya/browser-compat.git
cd browser-compat
go mod tidy
go run . parse -d ../browser-compat-data/
```

将在当前目录生成 `browser-compat.db`

查询

```sql
select * from browser_compat_data where browser='chrome' and browser_version in ('122','123','126')
```
