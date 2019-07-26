package main

import (
	"aif.io/k8s/protal/pkg/server"
	"fmt"
	"html/template"
	"k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"net/http"
	"sort"
	"time"
)
var(
    stop chan struct{}
)
func init(){
	stop=make(chan struct{})
}
func main(){
 fmt.Println("ok")
 store:= Start("","")
 http.Handle("/html/",http.StripPrefix("/html/",http.FileServer(http.Dir("html"))))
 http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
	 services:=store.List()
	 serverInfos:=[]server.ServiceInfo{}
	 for _,v:=range services{
		 svc:=v.(*v1.Service)
		 ports:=[]server.Port{}
		 for _,p:=range svc.Spec.Ports{
		 	 pName:=p.Name
		 	 if pName == ""{
		 	 	pName="https"
			 }
			 if svc.Spec.Type == "NodePort"{
				 ports=append(ports,server.Port{Name:pName,Target:p.TargetPort.String(),Protocol:string(p.Protocol),Url:fmt.Sprintf("%s://%s:%d",pName,"iseex.picp.io",p.NodePort)})
			 }else if svc.Spec.Type == "LoadBalancer"{
				 ports=append(ports,server.Port{Name:pName,Target:p.TargetPort.String(),Protocol:string(p.Protocol),Url:fmt.Sprintf("%s://%s:%d",pName,"iseex.picp.io",p.Port)})
			 }
		 }
		 if len(ports)>0 {
			 serverInfos=append(serverInfos,server.ServiceInfo{Name:svc.Name,ServerIp:"iseex.picp.io",Ports:ports})
		 }
	 }
	 t, err :=template.ParseFiles("html/index.html")
	 if err != nil {
	 	fmt.Fprintf(w,err.Error())
	 }else{
		 sort.Slice(serverInfos, func(i, j int) bool {
			 return serverInfos[i].Name < serverInfos[j].Name
		 })
		 t.ExecuteTemplate(w, "layout", serverInfos)
	 }
 })
 http.ListenAndServe("127.0.0.1:8000", nil)
 <-stop
}

func Start(kubeconfig string,apiServerAddress string) cache.Store{
	var config *rest.Config
	var err error
	fmt.Println(kubeconfig,"-----",apiServerAddress)
	kubeconfig="config"
	if kubeconfig == "" {
		log.Printf("using in-cluster configuration")
		config, err = rest.InClusterConfig()
	} else {
		log.Printf("using configuration from '%s'", kubeconfig)
		config, err = clientcmd.BuildConfigFromFlags(apiServerAddress, kubeconfig)
	}
	if err != nil {
		panic(err)
	}
	client, err :=kubernetes.NewForConfig(config)
	sharedInformers := informers.NewSharedInformerFactoryWithOptions(client, 3*time.Second, informers.WithNamespace(""))

	svcInformer := sharedInformers.Core().V1().Services().Informer()
	go svcInformer.Run(stop)

	createCacheHandler(svcInformer,"Services")
	return svcInformer.GetStore()
}

func createCacheHandler(informer cache.SharedIndexInformer, otype string)  {
	informer.AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {

			},
			UpdateFunc: func(old, cur interface{}) {

			},
			DeleteFunc: func(obj interface{}) {

			},
		})
}
