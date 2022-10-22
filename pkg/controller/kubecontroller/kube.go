package kubecontroller

import (
	"bufio"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/qingwave/weave/pkg/authorization"
	"github.com/qingwave/weave/pkg/common"
	"github.com/qingwave/weave/pkg/controller"
	"github.com/qingwave/weave/pkg/library/kubernetes"
	"github.com/qingwave/weave/pkg/model"
	"github.com/qingwave/weave/pkg/service"

	"github.com/gin-gonic/gin"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/remotecommand"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	prefix        = "weave.io"
	platformLabel = prefix + "/platform"
	creator       = prefix + "/creator"
	updater       = prefix + "/updater"

	createAction = "create"
	updateAction = "update"
)

const (
	KubeNamespace   = "namespaces"
	KubeDeployment  = "deployments"
	KubeStatefulset = "statefulsets"
	KubeDaemonset   = "daemonsets"
	KubePod         = "pods"
	KubeService     = "services"
	KubeIngress     = "ingresses"
)

var (
	resourceMap = map[string]KubeResource{
		KubeNamespace: {
			Name:         KubeNamespace,
			Object:       func() client.Object { return &corev1.Namespace{} },
			ObjectList:   func() client.ObjectList { return &corev1.NamespaceList{} },
			IsNamespaced: false,
		},
		KubeDeployment: {
			Name:         KubeDeployment,
			Object:       func() client.Object { return &appsv1.Deployment{} },
			ObjectList:   func() client.ObjectList { return &appsv1.DeploymentList{} },
			IsNamespaced: true,
		},
		KubeStatefulset: {
			Name:         KubeStatefulset,
			Object:       func() client.Object { return &appsv1.StatefulSet{} },
			ObjectList:   func() client.ObjectList { return &appsv1.StatefulSetList{} },
			IsNamespaced: true,
		},
		KubeDaemonset: {
			Name:         KubeDaemonset,
			Object:       func() client.Object { return &appsv1.DaemonSet{} },
			ObjectList:   func() client.ObjectList { return &appsv1.DaemonSetList{} },
			IsNamespaced: true,
		},
		KubePod: {
			Name:         KubePod,
			Object:       func() client.Object { return &corev1.Pod{} },
			ObjectList:   func() client.ObjectList { return &corev1.PodList{} },
			IsNamespaced: true,
		},
		KubeService: {
			Name:         KubeService,
			Object:       func() client.Object { return &corev1.Service{} },
			ObjectList:   func() client.ObjectList { return &corev1.ServiceList{} },
			IsNamespaced: true,
		},
		KubeIngress: {
			Name:         KubeIngress,
			Object:       func() client.Object { return &networkingv1.Ingress{} },
			ObjectList:   func() client.ObjectList { return &networkingv1.IngressList{} },
			IsNamespaced: true,
		},
	}
)

type KubeResource struct {
	Name         string
	Object       func() client.Object
	ObjectList   func() client.ObjectList
	IsNamespaced bool
}

func NewKubeControllers(client *kubernetes.KubeClient, groupService service.GroupService) controller.Controller {
	kc := &KubeController{client: client, groupService: groupService}
	if client != nil {
		kc.proxy = ProxyKubeAPIServer(client.GetConfig())
	}
	return kc
}

type KubeController struct {
	client       *kubernetes.KubeClient
	proxy        gin.HandlerFunc
	groupService service.GroupService
}

// @Summary Create kube resource
// @Description Create kube resource
// @Accept json
// @Produce json
// @Tags k8s
// @Security JWT
// @Param app body map[string]interface{} true "k8s json info"
// @Success 200 {object} common.Response{data=map[string]interface{}}
// @Router /api/v1/namespaces/:namespace/:resource [post]
func (kc *KubeController) Create(c *gin.Context) {
	ok, kubeResource := kc.checkParams(c)
	if !ok {
		return
	}

	obj := kubeResource.Object()
	if err := c.BindJSON(obj); err != nil {
		common.ResponseFailed(c, http.StatusBadRequest, err)
		return
	}

	namespace := c.Param("namespace")
	if namespace != "" {
		obj.SetNamespace(namespace)
	}

	kc.setDefaultParams(obj, common.GetUser(c), createAction)

	if err := kc.client.Create(context.Background(), obj); err != nil {
		common.ResponseFailed(c, http.StatusBadRequest, err)
		return
	}

	common.ResponseSuccess(c, obj)
}

// @Summary List kube resource
// @Description List kube resource
// @Accept json
// @Produce json
// @Tags k8s
// @Security JWT
// @Success 200 {object} common.Response{data=map[string]interface{}}
// @Router /api/v1/namespaces/:namespace/:resource [get]
func (kc *KubeController) List(c *gin.Context) {
	ok, kubeResource := kc.checkParams(c)
	if !ok {
		return
	}

	list := kubeResource.ObjectList()
	if err := kc.client.List(context.Background(), list, client.InNamespace(c.Param("namespace"))); err != nil {
		common.ResponseFailed(c, http.StatusInternalServerError, err)
		return
	}

	user := common.GetUser(c)
	switch kubeResource.Name {
	case KubeNamespace:
		if !authorization.IsClusterAdmin(user) {
			if err := kc.listNamespacesByGroup(list.(*corev1.NamespaceList), user.Name); err != nil {
				common.ResponseFailed(c, http.StatusInternalServerError, err)
				return
			}
		}
	}

	common.ResponseSuccess(c, list)
}

func (kc *KubeController) listNamespacesByGroup(list *corev1.NamespaceList, user string) error {
	groups, err := kc.groupService.List()
	if err != nil {
		return err
	}

	inGroup := func(name string) bool {
		for _, g := range groups {
			if g.Name == name {
				return true
			}
		}
		return false
	}

	namespaces := make([]corev1.Namespace, 0)
	for _, namespace := range list.Items {
		if inGroup(namespace.Name) {
			namespaces = append(namespaces, namespace)
		}
	}

	list.Items = namespaces
	return nil
}

// @Summary Update kube resource
// @Description Update kube resource
// @Accept json
// @Produce json
// @Tags k8s
// @Security JWT
// @Param app body map[string]interface{} true "k8s json info"
// @Success 200 {object} common.Response{data=map[string]interface{}}
// @Router /api/v1/namespaces/:namespace/:resource/:name [put]
func (kc *KubeController) Update(c *gin.Context) {
	ok, kubeResource := kc.checkParams(c)
	if !ok {
		return
	}

	obj := kubeResource.Object()
	if err := c.BindJSON(obj); err != nil {
		common.ResponseFailed(c, http.StatusBadRequest, err)
		return
	}

	namespace := c.Param("namespace")
	name := c.Param("name")
	if kubeResource.Name == KubeNamespace {
		name = namespace
	}

	if obj.GetName() != name || kubeResource.IsNamespaced && obj.GetNamespace() != namespace {
		common.ResponseFailed(c, http.StatusBadRequest, fmt.Errorf("request resource not matched"))
		return
	}

	kc.setDefaultParams(obj, common.GetUser(c), createAction)

	if err := kc.client.Update(context.Background(), obj); err != nil {
		common.ResponseFailed(c, http.StatusBadRequest, err)
		return
	}

	common.ResponseSuccess(c, nil)
}

// @Summary Get kube resource
// @Description Get kube resource
// @Accept json
// @Produce json
// @Tags k8s
// @Security JWT
// @Success 200 {object} common.Response{data=map[string]interface{}}
// @Router /api/v1/namespaces/:namespace/:resource/:name [get]
func (kc *KubeController) Get(c *gin.Context) {
	ok, kubeResource := kc.checkParams(c)
	if !ok {
		return
	}

	obj := kubeResource.Object()
	var namespacedName types.NamespacedName
	namespace := c.Param("namespace")
	name := c.Param("name")
	if kubeResource.IsNamespaced {
		namespacedName.Namespace = namespace
		namespacedName.Name = name
	} else if namespace != "" {
		namespacedName.Name = namespace
	} else {
		namespacedName.Name = name
	}

	if err := kc.client.Get(context.Background(), namespacedName, obj); err != nil {
		common.ResponseFailed(c, http.StatusBadRequest, err)
		return
	}

	common.ResponseSuccess(c, obj)
}

// @Summary Delete kube resource
// @Description Delete kube resource
// @Accept json
// @Produce json
// @Tags k8s
// @Security JWT
// @Success 200 {object} common.Response
// @Router /api/v1/namespaces/:namespace/:resource/:name [delete]
func (kc *KubeController) Delete(c *gin.Context) {
	ok, kubeResource := kc.checkParams(c)
	if !ok {
		return
	}

	obj := kubeResource.Object()
	namespace := c.Param("namespace")
	name := c.Param("name")
	if kubeResource.IsNamespaced {
		obj.SetName(name)
		obj.SetNamespace(namespace)
	} else if namespace != "" {
		obj.SetName(namespace)
	} else {
		obj.SetName(name)
	}

	if err := kc.client.Delete(context.Background(), obj); err != nil {
		common.ResponseFailed(c, http.StatusBadRequest, err)
		return
	}

	common.ResponseSuccess(c, obj)
}

func (kc *KubeController) checkParams(c *gin.Context) (bool, *KubeResource) {
	user := common.GetUser(c)
	if user == nil {
		common.ResponseFailed(c, http.StatusUnauthorized, nil)
		return false, nil
	}

	value, _ := c.Get(common.KubeResourceContextKey)
	kubeResource, ok := value.(*KubeResource)
	if !ok {
		common.ResponseFailed(c, http.StatusBadRequest, fmt.Errorf("failed to get resource from context"))
		return false, nil
	}

	if kubeResource.IsNamespaced && c.Param("namespace") == "" {
		common.ResponseFailed(c, http.StatusBadRequest, fmt.Errorf("namespace must be set"))
		return false, nil
	}

	return true, kubeResource
}

func (kc *KubeController) setDefaultParams(obj client.Object, user *model.User, action string) {
	// set labels
	labels := obj.GetLabels()
	if labels == nil {
		labels = make(map[string]string)
	}
	switch action {
	case createAction:
		labels[creator] = user.Name
	case updateAction:
		labels[updater] = user.Name
	}
	labels[platformLabel] = "true"
	obj.SetLabels(labels)

	// set annotations
	annotations := obj.GetAnnotations()
	if annotations == nil {
		annotations = make(map[string]string)
	}
	annotations[platformLabel] = "true"
	obj.SetAnnotations(annotations)
}

func (kc *KubeController) PodLogs(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	if namespace == "" || name == "" {
		common.ResponseFailed(c, http.StatusBadRequest, fmt.Errorf("invaild request info"))
		return
	}

	ctx := c.Request.Context()
	reader, err := kc.client.Log(ctx, namespace, name, "", c.Query("follow") == "true")
	if err != nil {
		common.ResponseFailed(c, http.StatusBadRequest, err)
		return
	}
	defer reader.Close()

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		n, err := c.Writer.Write(append(scanner.Bytes(), '\n'))
		c.Writer.Flush()
		if err != nil {
			return
		}
		if n == 0 {
			break
		}
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// Allow connections from any Origin
	CheckOrigin: func(r *http.Request) bool { return true },
}

func (kc *KubeController) PodExec(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	if namespace == "" || name == "" {
		common.ResponseFailed(c, http.StatusBadRequest, fmt.Errorf("invaild request info"))
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		common.ResponseFailed(c, http.StatusInternalServerError, err)
		return
	}

	session := &session{conn, 10 * time.Second}

	tty := c.Query("tty") == "true"
	err = kc.client.Exec(namespace, name, c.Query("container"),
		c.QueryArray("command"),
		tty,
		remotecommand.StreamOptions{
			Stdin:  session,
			Stdout: session,
			Stderr: session,
			Tty:    tty,
		})
	if err != nil {
		session.Close(websocket.CloseAbnormalClosure, err.Error())
		return
	}

	session.Close(websocket.CloseNormalClosure, "closed")
}

func (kc *KubeController) RegisterRoute(api *gin.RouterGroup) {
	for kind := range resourceMap {
		res := resourceMap[kind]
		resourcePath := fmt.Sprintf("/%s", kind)
		resourceDetail := resourcePath + "/:name"
		if kind == KubeNamespace {
			resourceDetail = resourcePath + "/:namespace"
		} else if res.IsNamespaced {
			resourcePath = fmt.Sprintf("/namespaces/:namespace/%s", kind)
			resourceDetail = resourcePath + "/:name"
		}

		api.GET(resourcePath, wrap(&res, kc.List))
		api.POST(resourcePath, wrap(&res, kc.Create))
		api.GET(resourceDetail, wrap(&res, kc.Get))
		api.PUT(resourceDetail, wrap(&res, kc.Update))
		api.DELETE(resourceDetail, wrap(&res, kc.Delete))
	}
	api.GET("/namespaces/:namespace/pods/:name/log", kc.proxy)
	api.Any("/namespaces/:namespace/pods/:name/exec", kc.PodExec)
}

func (kc *KubeController) Name() string {
	return "Kubernetes"
}

func wrap(res *KubeResource, h gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(common.KubeResourceContextKey, res)
		h(c)
	}
}

type session struct {
	conn      *websocket.Conn
	writeWait time.Duration
}

func (s *session) Read(p []byte) (int, error) {
	messageType, data, err := s.conn.ReadMessage()
	if err != nil {
		return 0, err
	}

	switch messageType {
	case websocket.TextMessage, websocket.BinaryMessage:
		return copy(p, data), nil
	case websocket.CloseMessage:
		return 0, nil
	}
	return 0, fmt.Errorf("unknown message type %d", messageType)
}

func (s *session) Write(p []byte) (int, error) {
	if s.writeWait > 0 {
		s.conn.SetWriteDeadline(time.Now().Add(s.writeWait))
	}
	err := s.conn.WriteMessage(websocket.TextMessage, p)
	if err != nil {
		return 0, err
	}
	return len(p), nil
}

func (s *session) Close(code int, reason string) {
	s.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(
		code,
		reason,
	))
	s.conn.Close()
}
