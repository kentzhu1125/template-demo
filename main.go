package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type Student struct {
	ID     int
	Gender string
	Name   string
	Sno    string
}
type LogBody struct {
	Timestamp   string `json:"timestamp"`
	AppName     string `json:"app_name"`
	LogLevel    string `json:"log_level"`
	Thread      string `json:"thread"`
	Tid         string `json:"tid"`
	ClassName   string `json:"class_name"`
	ClassMethod string `json:"class_method"`
	LineNumber  string `json:"line_number"`
	Message     string `json:"message"`
}

func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}
func simulateLog(logBody LogBody, logfile string) {
	logBodyStr, err := json.Marshal(logBody)
	if !checkFileIsExist(logfile) { //文件不存在
		fmt.Println("日志文件不存在，创建文件")
		os.Create(logfile)
	}
	fileWriter, err := os.OpenFile(logfile, os.O_RDWR|os.O_APPEND, 0660) //打开日志文件
	defer fileWriter.Close()                                             //关闭文件
	if err != nil {
		fmt.Println("打开日志文件错误")
		return
	}
	//_,err=io.WriteString(fileWriter,string(logBodyStr))
	_, err = fileWriter.WriteString(string(logBodyStr) + "\n")
	if err != nil {
		fmt.Println("写入日志文件错误", err.Error())
		return
	}

}

func getLocalIpAddress() []string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println("获取本地IP地址失败")
		return nil
	}
	ipAddress := []string{}
	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ipAddress = append(ipAddress, ipnet.IP.String())
			}
		}
	}
	return ipAddress
}

var nacosContent string

//update nacosContent Globally
func updateContentRetrievedFromNacos(data string) {
	nacosContent = data
}

//get local ip address
func getIpAddress() string {
	addrs, err := net.InterfaceAddrs()
	var localIPAddress string
	if err != nil {
		fmt.Println("获取本地IP 地址失败", err)
		return ""
	}
	//only return the first non loopback ip address
	for _, address := range addrs {
		// check ip address is loopback or not
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				//fmt.Println(ipnet.IP.String())
				localIPAddress = ipnet.IP.String()
				break
			}
		}
	}
	return localIPAddress
}

func main() {
	app := gin.Default()
	appName := "APIDemo"
	thread := "Main"
	//tid := "1000"
	className := "Main_Class"
	classMethod := "Main_Method"
	lineNumber := "1000"
	logDir := "/logs/apidemo"
	err := os.MkdirAll(logDir, os.ModePerm)
	if err != nil {
		fmt.Println("创建日志目录失败")
	}
	logFile := logDir + "/apidemo_json.log"
	//logFile ="apidemo_json.log"
	app.GET("/master/simulate/error", func(c *gin.Context) {
		//logLevel := c.Param("logLevel")
		logLevel := "error"
		logBody := LogBody{
			time.Now().String(),
			appName,
			logLevel,
			thread,
			strconv.Itoa(os.Getpid()),
			className,
			classMethod,
			lineNumber,
			"This is simulated log ,level:" + logLevel,
		}
		simulateLog(logBody, logFile)
		c.JSON(http.StatusOK, gin.H{
			"code":      200,
			"msg":       "simulate log successfully",
			"IPAddress": getLocalIpAddress()})
	})

	app.GET("/master/get", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code":      200,
			"msg":       "api address: /get",
			"IPAddress": getLocalIpAddress()})

	})
	app.GET("/master/simulate409", func(c *gin.Context) {
		c.JSON(http.StatusConflict, gin.H{
			"code":      409,
			"msg":       "api address: /simulate409",
			"IPAddress": getLocalIpAddress()})

	})
	app.HEAD("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code":      200,
			"msg":       "api address: /",
			"IPAddress": getLocalIpAddress()})

	})
	app.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code":      200,
			"msg":       "Thanks for accessing apidemo.api address: /",
			"IPAddress": getLocalIpAddress()})

	})
	app.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code":      200,
			"msg":       "api address: /",
			"IPAddress": getLocalIpAddress()})

	})
	app.GET("/master/error", func(c *gin.Context) {
		fmt.Println(time.Now(), " [error] this is simulated log ,level:info")
		c.JSON(http.StatusOK, gin.H{
			"code":      200,
			"msg":       "api address: /error",
			"IPAddress": getLocalIpAddress()})
	})
	app.GET("/master/sleep", func(c *gin.Context) {
		fmt.Println(time.Now(), " [sleep] this is simulated sleep ,level:info")
		time.Sleep(70 * time.Second)
		c.JSON(http.StatusOK, gin.H{
			"code":      200,
			"msg":       "api address: /sleep",
			"IPAddress": getLocalIpAddress()})
	})
	app.GET("/master/info", func(c *gin.Context) {
		fmt.Println(time.Now(), " [info] this is simulated log ,level:info")
		c.JSON(http.StatusOK, gin.H{
			"code":      200,
			"msg":       "api address: /info",
			"IPAddress": getLocalIpAddress()})
	})
	app.GET("/master/release", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code":      200,
			"msg":       "release: /release",
			"IPAddress": getLocalIpAddress()})
	})
	app.POST("/master/post", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code":      200,
			"msg":       "api address: /post",
			"IPAddress": getLocalIpAddress()})
	})
	//get nacos configuration from OS environment variables
	nacosAddr := os.Getenv("NACOS_ADDR")
	nacosNamespace := os.Getenv("NACOS_NAMESPACE")
	nacosUsername := os.Getenv("NACOS_USERNAME")
	nacosPassword := os.Getenv("NACOS_PASSWORD")
	//check if environment is nil
	fmt.Println("nacos地址是", nacosAddr)
	//channel used to control nacos registry
	var registerCompleted chan int = make(chan int)
	if nacosAddr != "" && nacosNamespace != "" && nacosUsername != "" && nacosPassword != "" {
		nacosIpAddr := strings.Split(nacosAddr, ":")[0]
		nacosPort, _ := strconv.ParseUint(strings.Split(nacosAddr, ":")[1], 10, 64)
		fmt.Println("Nacos连接信息如下：")
		fmt.Println(nacosAddr, " ", nacosNamespace, " ", nacosUsername, " ", nacosPassword, " ", nacosIpAddr, " ", nacosPort)
		//below are code related to nacos
		sc := []constant.ServerConfig{
			{
				IpAddr:      nacosIpAddr,
				Port:        nacosPort,
				Scheme:      "http",
				ContextPath: "/nacos",
			},
		}
		//or a more graceful way to create ServerConfig

		//initial config for nacos server connection
		//sc := []constant.ServerConfig{
		//	*constant.NewServerConfig(
		//		nacosIpAddr,
		//		nacosPort,
		//		constant.WithScheme("http"),
		//		constant.WithContextPath("/nacos")),
		//}
		//initial config for nacos client connnection
		cc := constant.ClientConfig{
			NamespaceId:         nacosNamespace, //namespace id
			TimeoutMs:           5000,
			NotLoadCacheAtStart: true,
			//LogDir:              "/tmp/nacos/log",
			//CacheDir:            "/tmp/nacos/cache",
			LogLevel: "debug",
			Username: nacosUsername,
			Password: nacosPassword,
		}
		//or a more graceful way to create ClientConfig
		//_ = *constant.NewClientConfig(
		//	constant.WithNamespaceId("e525eafa-f7d7-4029-83d9-008937f9d468"),
		//	constant.WithTimeoutMs(5000),
		//	constant.WithNotLoadCacheAtStart(true),
		//	constant.WithLogDir("/tmp/nacos/log"),
		//	constant.WithCacheDir("/tmp/nacos/cache"),
		//	constant.WithLogLevel("debug"),
		//)

		// a more graceful way to create config client
		client, err := clients.NewConfigClient(
			vo.NacosClientParam{
				ClientConfig:  &cc,
				ServerConfigs: sc,
			},
		)
		if err != nil {
			fmt.Println("Connect to nacos failed")
			//panic(err)
		}
		//fetch config from nacos server
		content, err := client.GetConfig(vo.ConfigParam{
			DataId: "application.yaml",
			Group:  "apidemo",
		})
		//assign initial value to global content parameter
		nacosContent = content
		if err != nil {
			fmt.Println("get content in nacos failed")
		}

		//listen to config
		err = client.ListenConfig(vo.ConfigParam{
			DataId: "application.yaml",
			Group:  "apidemo",
			OnChange: func(namespace, group, dataId, data string) {
				fmt.Println("configuration has been changed,new data is as below \n", "group:"+group+", dataId:"+dataId+", data:\n"+data)
				updateContentRetrievedFromNacos(data)

			},
		})
		if err != nil {
			fmt.Println("listen config change failed")
		}
		// register service to nacos registry

		// Another way of create naming client for service discovery (recommend)
		namingClient, err := clients.NewNamingClient(
			vo.NacosClientParam{
				ClientConfig:  &cc,
				ServerConfigs: sc,
			},
		)
		// get host ip address,used for instance identity in nacos discovery
		ipAddr := getIpAddress()

		go func() {
			//try to register to nacos 5 times,then exit
			for i := 0; i < 5; i++ {
				success, err := namingClient.RegisterInstance(vo.RegisterInstanceParam{
					Ip:          ipAddr,
					Port:        8888,
					ServiceName: "apidemo",
					Weight:      10,
					Enable:      true,
					Healthy:     true,
					Ephemeral:   true,
					Metadata:    map[string]string{"company": "cec"},
					ClusterName: "cluster-a", // default value is DEFAULT
					GroupName:   "group-a",   // default value is DEFAULT_GROUP
				})
				if success {
					fmt.Println("Register service to nacos successfully")
					break
				} else {
					fmt.Println("Register service to nacos failed", "error message:\n", err)
					time.Sleep(5 * time.Second)
				}
			}
			registerCompleted <- 0

		}()
	} else {
		fmt.Println("获取到的Nacos环境变量为空")
	}

	//display the content of config
	app.GET("/master/nacosConfig", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"msg":     "/master/nacosConfig",
			"content": nacosContent})
	})

	//success, err := namingClient.RegisterInstance(vo.RegisterInstanceParam{
	//	Ip:          ipAddr,
	//	Port:        8888,
	//	ServiceName: "apidemo",
	//	Weight:      10,
	//	Enable:      true,
	//	Healthy:     true,
	//	Ephemeral:   true,
	//	Metadata:    map[string]string{"company": "cec"},
	//	ClusterName: "cluster-a", // default value is DEFAULT
	//	GroupName:   "group-a",   // default value is DEFAULT_GROUP
	//})
	//if success {
	//	fmt.Println("Register service to nacos successfully")
	//} else {
	//	fmt.Println("Register service to nacos failed", "error message:\n", err)
	//}

	//waiting for nacos registration finish
	if nacosAddr != "" && nacosNamespace != "" && nacosUsername != "" && nacosPassword != "" {
		<-registerCompleted
	}

	//start gin server
	app.Run(":8888")

}
