package kubernetes

import (
	"context"
	"io"
	"time"

	weaveconfig "github.com/qingwave/weave/pkg/config"

	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

var (
	scheme = runtime.NewScheme()

	codec = serializer.NewCodecFactory(scheme)

	neverResync time.Duration = 0
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(apiextensionsv1.AddToScheme(scheme))
}

type KubeClient struct {
	client.Client
	podClient  rest.Interface
	cache      cache.Cache
	config     *rest.Config
	kubeconfig *weaveconfig.KubeConfig
}

func NewClient(kubeconfig *weaveconfig.KubeConfig) (*KubeClient, error) {
	config, err := config.GetConfig()
	if err != nil {
		return nil, err
	}

	podClient, err := apiutil.RESTClientForGVK(schema.GroupVersionKind{Version: "v1", Kind: "Pod"}, false, config, codec)
	if err != nil {
		return nil, err
	}

	cache, err := cache.New(config, cache.Options{
		Scheme:          scheme,
		Resync:          &neverResync,
		DefaultSelector: cache.ObjectSelector{},
	})
	if err != nil {
		return nil, err
	}

	runtimeClient, err := newRuntimeClient(cache, config)
	if err != nil {
		return nil, err
	}

	c := &KubeClient{
		Client:     runtimeClient,
		podClient:  podClient,
		cache:      cache,
		config:     config,
		kubeconfig: kubeconfig,
	}

	for _, resource := range kubeconfig.WatchResources {
		gvk, _ := schema.ParseKindArg(resource)
		if gvk == nil {
			logrus.Warnf("invaild kubernetes gvk [%s]", resource)
			continue
		}

		obj, err := scheme.New(*gvk)
		if err != nil {
			logrus.Warnf("invaild kubernetes gvk [%s]: %v", gvk, err)
			continue
		}

		if err := c.Watch(obj.(client.Object)); err != nil {
			logrus.Warnf("watch kind error: %v", err)
		}
	}

	return c, nil
}

func (c *KubeClient) StartCache() (err error) {
	go func() {
		err = c.cache.Start(context.Background())
		if err != nil {
			logrus.Warnf("failed to start cache: %v", err)
		}
	}()

	if ok := c.cache.WaitForCacheSync(context.Background()); !ok {
		logrus.Warn("failed to sync all cache")
	}

	return
}

func (c *KubeClient) GetConfig() *rest.Config {
	return c.config
}

func (c *KubeClient) Watch(objs ...client.Object) error {
	for _, obj := range objs {
		if _, err := c.cache.GetInformer(context.Background(), obj); err != nil {
			return err
		}
	}
	return nil
}

func (c *KubeClient) Log(ctx context.Context, namespace, pod, container string, follow bool) (io.ReadCloser, error) {
	req := c.podClient.Get().
		Resource("pods").
		Namespace(namespace).
		Name(pod).
		SubResource("log").
		VersionedParams(&corev1.PodLogOptions{
			Container: container,
			Follow:    follow,
		}, clientgoscheme.ParameterCodec)

	return req.Stream(ctx)
}

func (c *KubeClient) Exec(namespace, pod, container string, cmd []string, tty bool, options remotecommand.StreamOptions) error {
	req := c.podClient.Post().
		Resource("pods").
		Namespace(namespace).
		Name(pod).
		SubResource("exec").
		VersionedParams(&corev1.PodExecOptions{
			Command:   cmd,
			Stdin:     true,
			Stdout:    true,
			Stderr:    true,
			TTY:       tty,
			Container: container,
		}, clientgoscheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(c.config, "POST", req.URL())
	if err != nil {
		logrus.Warnf("failed to create executor: %v", err)
		return err
	}

	return exec.Stream(options)
}

func newRuntimeClient(cache cache.Cache, config *rest.Config) (client.Client, error) {
	c, err := client.New(config, client.Options{
		Scheme: scheme,
	})
	if err != nil {
		return nil, err
	}

	return client.NewDelegatingClient(client.NewDelegatingClientInput{
		CacheReader: cache,
		Client:      c,
	})
}
