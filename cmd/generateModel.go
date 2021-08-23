package cmd

/**
 * @Description: 生成model
 * @author ganJinBiao
 * @return unc
 * date 2021-04-28 09:32:41
 */
import (
	"cli/converter"
	"fmt"
	"github.com/spf13/cobra"
)

/**
 * @Description:定义配置结构体
 * @author ganJinBiao
 */
type config struct {
	Prefix     string //表前缀
	SavePath   string //生成文件保存路径,请写对路径，别问为什么，问就是我菜。
	DbUserName string //数据库账号名
	DbConfig   string //设置
	DbName     string //数据库名
	DbPassword string //密码
	Host       string //ip
}

//请在这里定义配置
var (
	configObj = config{
		Prefix:     "prefix",                          //表前缀
		SavePath:   "D:\\programming\\GO\\cli\\model", //生成文件保存路径,请写对路径，别问为什么，问就是我菜。
		DbUserName: "root",                            //数据库账号名
		DbConfig:   "?charset=utf8mb4",                //设置
		DbName:     "demo",                            //数据库名
		DbPassword: "root",                            //密码
		Host:       "localhost:13306",                 //ip
	}
)

// generateModelCmd represents the generateModel command
var generateModelCmd = &cobra.Command{
	Use:   "generateModel",
	Short: "生成对应的表结构的model",
	Long:  `使用前请查看readme.md文档，配置好后在使用`,
	Run: func(cmd *cobra.Command, args []string) {

		tableName, _ := cmd.Flags().GetString("tableName")
		fmt.Println("将生成", tableName, "表的model结构体")
		//定义dsn
		dsn := configObj.DbUserName + ":" + configObj.DbPassword + "@tcp(" + configObj.Host + ")/" + configObj.DbName + configObj.DbConfig
		//定义保存路径
		path := configObj.SavePath + "/" + tableName + ".go"

		// 初始化
		obj := converter.NewTableToStruct()
		// 个性化配置
		obj.Config(&converter.T2tConfig{
			// 如果字段首字母本来就是大写, 就不添加tag, 默认false添加, true不添加
			RmTagIfUcFirsted: false,
			// tag的字段名字是否转换为小写, 如果本身有大写字母的话, 默认false不转
			TagToLower: false,
			// 字段首字母大写的同时, 是否要把其他字母转换为小写,默认false不转换
			UcFirstOnly: false,
			//// 每个struct放入单独的文件,默认false,放入同一个文件(暂未提供)
			//SeperatFile: false,
		})
		// 开始迁移转换
		err := obj.
			// 指定某个表,如果不指定,则默认全部表都迁移
			Table(tableName).
			// 表前缀
			Prefix(configObj.Prefix).
			// 是否添加json tag
			EnableJsonTag(true).
			// 生成struct的包名(默认为空的话, 则取名为: package model)
			PackageName("model").
			// tag字段的key值,默认是orm
			TagKey("gorm").
			// 是否添加结构体方法获取表名
			RealNameMethod("TableName").
			// 生成的结构体保存路径
			SavePath(path).
			// 数据库dsn
			Dsn(dsn).
			// 执行
			Run()
		if err != nil {
			fmt.Println(tableName, "模型生成失败，请重试", err.Error())
			return
		}
		fmt.Println(tableName, "模型生成成功")

	},
}

func init() {
	rootCmd.AddCommand(generateModelCmd)

	//表名参数
	generateModelCmd.PersistentFlags().StringP("tableName", "n", "", "指定表名，省略前缀")
}
