package patch

import (
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
	"strings"
)

// NotFound 用来保证和第三方库解耦，这里的第三方库是 sql 包
var NotFound = errors.New("not found")
var notFoundCode = 40001
var systemErr = 50001

func Biz() error {
	err := Dao("select user_name from users where id = 666;")
	if errors.Is(err, NotFound) {
		// 要站在业务的角度考虑，未找到，是否是正常的
		// 这里假设是正常的
		return nil
	}
	if err != nil {
		// 出现了其他查询问题，可以直接向上传递
		return  err
	}
	return nil
}

func Dao(query string) error {

	err := mockError()
	if err == sql.ErrNoRows {
		// 封装好查询参数，带上堆栈信息
		return errors.Wrapf(NotFound, fmt.Sprintf("data not found, sql %s", query))
	}
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("db query some system error sql: %s", query))
	}
	// do something
	return nil

}

// Another handler

func Biz1() error {
	err := Dao1("select user_name from users where id = 666;")

	if err != nil {
		// 出现了其他查询问题，可以直接向上传递
		return err
	}
	return nil
}

func Dao1(query string) error {
	err := mockError()

	if err != nil {
		// 我们没有仔细区别 err 是什么，反正就是告诉上游，出错了
		return errors.Wrapf(NotFound, fmt.Sprintf("data not found, sql %s", query))
	}
	// do something
	return nil

}

// Another...return error code

func Biz2() error {
	err := Dao2("select user_name from users where id = 666;")

	if IsNoRow(err) {
		// 不管怎么说，出现了数据库查询的问题，可以转为业务领域错误，也可以继续向上传递
		return err
	} else if err != nil {
		return err
	}
	return  nil
}

func Dao2(query string) error {
	err := mockError()

	if err == sql.ErrNoRows {
		// 我们没有仔细区别 err 是什么，反正就是告诉上游，出错了
		return fmt.Errorf("%d, not found", notFoundCode)
	}
	if err != nil {
		return fmt.Errorf("%d, not found", systemErr)
	}
	return nil


}

func IsNoRow(err error) bool {
	return strings.HasPrefix(err.Error(), fmt.Sprintf("%d", notFoundCode))
}

func mockError() error {
	return sql.ErrNoRows
}