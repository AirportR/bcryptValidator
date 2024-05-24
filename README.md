# bcryptValidator
v2board 数据库中用户密码的部分验证
input 文件格式示例：
```
6190491826@qq.com,$2y$10$k2PgxO78FEzQSZrkSWd5H.wHitOUAkEhwZiNs/E8xFNBjeA17oG1S
chajsdja9652@gmail.com,$2y$10$NZj11bMh0DagsBjSqOZjfOXUq3Hns.42IY/2lxlfNSa2Ey73vD0KW
sdadasda@qq.com,$2y$10$WUcT5IbG99FhpIhKtmQTN.U3nCA363fhy3va3pwOU4/nA548Sgs9W
4095774031@qq.com,$2y$10$wM6UCnbI88M0CWd8Z5JwJu9wAWv9YYOLfz4TiuTYzdxn8zFXnH68C
787263404@qq.com,$2y$10$JO4P/zo34rLGojtZGwtAt.l20JKL9pQxULpbDWsrEMlWUmRPU1t9u
```
## 结果演示
![image](https://github.com/AirportR/bcryptValidator/assets/104753221/08aa3f32-f39f-4f41-8ce6-55c3cae7862b)

## 编译
```bash
go mod tidy
go build .
```

## 使用
```bash
chmod +x bcryptValidator && ./bcryptValidator -h
```

## 示例
```bash
./bcryptValidator -input v2_user.txt -password 123456789 -count 1000
```

