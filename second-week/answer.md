
## Question: 
我们在数据库操作的时候，比如 dao 层中当遇到一个 sql.ErrNoRows 的时候，是否应
该 Wrap 这个 error，抛给上层。为什么，应该怎么做请写出代码？

## Answer:
应该Wrap这个error，抛给上层。
需要将堆栈信息给到上层调用者，在dao层warp。

详见: [demo](./main.go)