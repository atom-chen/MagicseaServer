package login

import (
	"MagicseaServer/GAServer/service"
	"fmt"
	"net/http"
	"MagicseaServer/GAServer/config"
	"MagicseaServer/GAServer/log"
	"strings"
	"strconv"
	"MagicseaServer/Server/cluster"
	"MagicseaServer/gameproto/msgs"
	"MagicseaServer/gameproto"
	"github.com/gogo/protobuf/proto"
)

type LoginService struct{
	service.ServiceData
}

//Service 获取服务对象
func Service() service.IService {
	return new(LoginService)
}

func Type() string {
	return "login"
}
//以下为接口函数
func (s *LoginService) OnReceive(context service.Context) {
	fmt.Println("center.OnReceive:", context.Message())
}
func (s *LoginService) OnInit() {

}

func (s *LoginService) OnStart(as *service.ActorService) {


	go func() {
		//开启http服务
		http.HandleFunc("/login", login)

		http.HandleFunc("/regist", regist)
		httpAddr := config.GetServiceConfigString(s.Name, "httpAddr")
		log.Println("login listen http:", s.Name, "  ", httpAddr)
		http.ListenAndServe(httpAddr, nil)
	}()
}

//注册
func regist(w http.ResponseWriter, req *http.Request) {

	log.Info("=========打开  注册 ============")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	req.ParseForm()
	if req.Form["a"] ==nil || req.Form["p"] == nil{
		log.Error("a,p is empty :",req.Form)
		return
	}
	//账号
	acc := ""
	if al,ok := req.Form["a"];ok{
		acc = al[0]
	}
	//密码
	pwd:= ""
	if al, ok := req.Form["p"]; ok {
		pwd = al[0]
	}
	if len(acc) < 1 || len(pwd) < 1 {
		registBackError(w, "账号密码都不能为空", nil)
		return
	}
	log.Info("reg account:acc=%s,pwd=%s", acc, pwd)
	//key := "User:nameindex:" + acc
	//todo db
	w.Write([]byte("success"))

}
func registBackError(w http.ResponseWriter, val string, e error) {
	log.Error("create user db error:%s,%v", val, e)
	w.Write([]byte(val))
}

func login(w http.ResponseWriter, req *http.Request) {

	req.ParseForm()
	if req.Form["a"] == nil || req.Form["p"] == nil {
		log.Error("a,p is empty:", req.Form)
		return
	}
	//账号
	acc := ""
	if al, ok := req.Form["a"]; ok {
		acc = al[0]
	}
	log.Println("login account:", acc)
	strs := strings.Split(acc, "_")
 	id, _ := strconv.Atoi(strs[0])

	resp, err := cluster.GetServicePID("session").Ask(&msgs.UserLogin{acc, uint64(id)})
	if err == nil {
		var s, _ = resp.(*gameproto.UserLoginResult).Marshal()
		//var s, _ = json.Marshal(resp)
		w.Write(s)
		log.Info("login ok:msg=%v", resp)
	} else {
		loginBackError(w, err)
		log.Println("login error:", acc, err)
	}
}

func loginBackError(w http.ResponseWriter, e error) {
	log.Error("create user db :%v", e)
	d, _ := proto.Marshal(&gameproto.UserLoginResult{Result: int32(msgs.Error)})
	w.Write(d)
}
