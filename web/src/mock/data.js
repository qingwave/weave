import { setUser } from '@/utils';

const data = [
    {
        path: "/api/v1/auth/token",
        method: "post",
        do: (params) => {
            let user = JSON.parse(params.body);
            setUser({
                name: user.name ? user.name : 'admin',
                id: 1,
                email: "admin@weave.com",
                avatar: "",
                authInfos: [],
                groups: [{
                    id: 1,
                    name: "root",
                    kind: "system",
                    describe: "weave+system+group",
                    roles: []
                }]
            });
            return {}
        }
    },
    {
        path: "/api/v1/users",
        method: "get",
        data: [
            {
                name: "admin",
                id: '1',
                email: "admin@weave.com",
            },
            {
                name: "alice",
                id: '2',
                email: "alice@weave.com",
            },
            {
                name: "bob",
                id: '3',
                email: "bob@weave.com",
            },
        ]
    },
    {
        path: /\/api\/v1\/users\/\d+\/groups/,
        method: "get",
        data: [
            {
                "id": 1,
                "name": "root",
                "kind": "system",
                "describe": "weave system group",
                "createdAt": "2023-05-11T01:38:23.500389Z",
                "updatedAt": "0001-01-01T00:00:00Z"
            }
        ]
    },
    {
        path: /\/api\/v1\/users\/\d+/,
        method: "get",
        data: {
            name: "admin",
            id: '1',
            email: "admin@weave.com",
        }
    },
    {
        path: '/api/v1/containers',
        method: "get",
        data: [
            {
                id: "abcd123",
                image: "nginx",
                name: "app-nginx",
                port: 80,
                startAt: "2023-05-11T02:35:33",
                status: "running"
            },
            {
                id: "abcd124",
                image: "server:v1.0",
                name: "server",
                port: 80,
                startAt: "2023-04-10T12:42:51",
                status: "stop"
            },
            {
                id: "abcd125",
                image: "app:v1.0",
                name: "app1",
                port: 80,
                startAt: "2023-02-08T09:12:00",
                status: "dead"
            },
        ]
    },
    {
        path: '/api/v1/groups',
        method: "get",
        data: [
            {
                "id": 1,
                "name": "root",
                "kind": "system",
                "describe": "weave system group",
                "createdAt": "2023-05-11T01:38:23.500389Z",
                "updatedAt": "0001-01-01T00:00:00Z"
            },
            {
                "id": 2,
                "name": "weave",
                "kind": "custom",
                "describe": "weave custom group",
                "createdAt": "2023-05-11T01:38:23.500389Z",
                "updatedAt": "0001-01-01T00:00:00Z"
            },
        ]
    },
    {
        path: /\/api\/v1\/groups\/.*/,
        method: "get",
        data: {
            "id": 1,
            "name": "root",
            "kind": "system",
            "describe": "weave system group",
            "createdAt": "2023-05-11T01:38:23.500389Z",
            "updatedAt": "0001-01-01T00:00:00Z"
        }
    },
    {
        path: '/api/v1/roles',
        method: 'get',
        data: [
            {"id":1,"name":"cluster-admin","scope":"cluster","namespace":"","rules":[{"resource":"*","operation":"*"}]},
            {"id":2,"name":"authenticated","scope":"cluster","namespace":"","rules":[{"resource":"users","operation":"*"},{"resource":"auth","operation":"*"}]},
            {"id":3,"name":"unauthenticated","scope":"cluster","namespace":"","rules":[{"resource":"auth","operation":"create"}]}]
    },
    {
        path: '/api/v1/namespaces',
        method: "get",
        data: {
            "items": [
                { "kind": "Namespace", "apiVersion": "v1", "metadata": { "name": "default", "creationTimestamp": "2023-05-11T02:49:51Z", "labels": { "kubernetes.io/metadata.name": "default" } }, "spec": { "finalizers": ["kubernetes"] }, "status": { "phase": "Active" } },
                { "kind": "Namespace", "apiVersion": "v1", "metadata": { "name": "kube-system", "creationTimestamp": "2023-05-11T02:49:49Z", "labels": { "kubernetes.io/metadata.name": "kube-system" } }, "spec": { "finalizers": ["kubernetes"] }, "status": { "phase": "Active" } },
                { "kind": "Namespace", "apiVersion": "v1", "metadata": { "name": "kube-public", "creationTimestamp": "2023-05-11T02:49:49Z", "labels": { "kubernetes.io/metadata.name": "kube-public" } }, "spec": { "finalizers": ["kubernetes"] }, "status": { "phase": "Active" } },
                { "kind": "Namespace", "apiVersion": "v1", "metadata": { "name": "kube-node-lease", "creationTimestamp": "2023-05-11T02:49:49Z", "labels": { "kubernetes.io/metadata.name": "kube-node-lease" } }, "spec": { "finalizers": ["kubernetes"] }, "status": { "phase": "Active" } }
            ]
        }
    },
    {
        path: /\/api\/v1\/namespaces\/.*\/deployment/,
        method: "get",
        data: {"items": [
            {"kind":"Deployment","apiVersion":"apps/v1","metadata":{"name":"weave-server","namespace":"kube-system","uid":"c021cca3-4ddf-40a6-a631-639c58966cbf","resourceVersion":"531","generation":1,"creationTimestamp":"2023-05-11T02:49:52Z","labels":{"k8s-app":"weave-server"}},"spec":{"replicas":2,"selector":{"matchLabels":{"k8s-app":"weave-server"}},"template":{"metadata":{"creationTimestamp":null,"labels":{"k8s-app":"weave-server"}},"spec":{"containers":[{"name":"weave","image":"k8s.gcr.io/weave/weave:v1.8.6","ports":[{"name":"dns","containerPort":53,"protocol":"UDP"},{"name":"dns-tcp","containerPort":53,"protocol":"TCP"},{"name":"metrics","containerPort":9153,"protocol":"TCP"}],"resources":{"limits":{"memory":"170Mi"},"requests":{"cpu":"100m","memory":"70Mi"}},"volumeMounts":[{"name":"config-volume","readOnly":true,"mountPath":"/etc/weave"}],"livenessProbe":{"httpGet":{"path":"/health","port":8080,"scheme":"HTTP"},"initialDelaySeconds":60,"timeoutSeconds":5,"periodSeconds":10,"successThreshold":1,"failureThreshold":5},"readinessProbe":{"httpGet":{"path":"/ready","port":8181,"scheme":"HTTP"},"timeoutSeconds":1,"periodSeconds":10,"successThreshold":1,"failureThreshold":3},"terminationMessagePath":"/dev/termination-log","terminationMessagePolicy":"File","imagePullPolicy":"IfNotPresent","securityContext":{"capabilities":{"add":["NET_BIND_SERVICE"],"drop":["all"]},"readOnlyRootFilesystem":true,"allowPrivilegeEscalation":false}}],"restartPolicy":"Always","terminationGracePeriodSeconds":30,"dnsPolicy":"Default","nodeSelector":{"kubernetes.io/os":"linux"},"serviceAccountName":"weave","serviceAccount":"weave","securityContext":{},"schedulerName":"default-scheduler","tolerations":[{"key":"CriticalAddonsOnly","operator":"Exists"},{"key":"node-role.kubernetes.io/master","effect":"NoSchedule"},{"key":"node-role.kubernetes.io/control-plane","effect":"NoSchedule"}],"priorityClassName":"system-cluster-critical"}},"strategy":{"type":"RollingUpdate","rollingUpdate":{"maxUnavailable":1,"maxSurge":"25%"}},"revisionHistoryLimit":10,"progressDeadlineSeconds":600},"status":{"observedGeneration":1,"replicas":2,"updatedReplicas":2,"readyReplicas":2,"availableReplicas":2,"conditions":[{"type":"Available","status":"True","lastUpdateTime":"2023-05-11T02:50:15Z","lastTransitionTime":"2023-05-11T02:50:15Z","reason":"MinimumReplicasAvailable","message":"Deployment has minimum availability."},{"type":"Progressing","status":"True","lastUpdateTime":"2023-05-11T02:50:16Z","lastTransitionTime":"2023-05-11T02:50:06Z","reason":"NewReplicaSetAvailable","message":"ReplicaSet \"weave-64897985d\" has successfully progressed."}]}},
            {"kind":"Deployment","apiVersion":"apps/v1","metadata":{"name":"weave-frontend","namespace":"kube-system","uid":"c021cca3-4ddf-40a6-a631-639c58966cbf","resourceVersion":"531","generation":1,"creationTimestamp":"2023-05-11T02:49:52Z","labels":{"k8s-app":"weave-server"}},"spec":{"replicas":2,"selector":{"matchLabels":{"k8s-app":"weave-server"}},"template":{"metadata":{"creationTimestamp":null,"labels":{"k8s-app":"weave-server"}},"spec":{"containers":[{"name":"weave","image":"k8s.gcr.io/weave/weave:v1.8.6","ports":[{"name":"dns","containerPort":53,"protocol":"UDP"},{"name":"dns-tcp","containerPort":53,"protocol":"TCP"},{"name":"metrics","containerPort":9153,"protocol":"TCP"}],"resources":{"limits":{"memory":"170Mi"},"requests":{"cpu":"100m","memory":"70Mi"}},"volumeMounts":[{"name":"config-volume","readOnly":true,"mountPath":"/etc/weave"}],"livenessProbe":{"httpGet":{"path":"/health","port":8080,"scheme":"HTTP"},"initialDelaySeconds":60,"timeoutSeconds":5,"periodSeconds":10,"successThreshold":1,"failureThreshold":5},"readinessProbe":{"httpGet":{"path":"/ready","port":8181,"scheme":"HTTP"},"timeoutSeconds":1,"periodSeconds":10,"successThreshold":1,"failureThreshold":3},"terminationMessagePath":"/dev/termination-log","terminationMessagePolicy":"File","imagePullPolicy":"IfNotPresent","securityContext":{"capabilities":{"add":["NET_BIND_SERVICE"],"drop":["all"]},"readOnlyRootFilesystem":true,"allowPrivilegeEscalation":false}}],"restartPolicy":"Always","terminationGracePeriodSeconds":30,"dnsPolicy":"Default","nodeSelector":{"kubernetes.io/os":"linux"},"serviceAccountName":"weave","serviceAccount":"weave","securityContext":{},"schedulerName":"default-scheduler","tolerations":[{"key":"CriticalAddonsOnly","operator":"Exists"},{"key":"node-role.kubernetes.io/master","effect":"NoSchedule"},{"key":"node-role.kubernetes.io/control-plane","effect":"NoSchedule"}],"priorityClassName":"system-cluster-critical"}},"strategy":{"type":"RollingUpdate","rollingUpdate":{"maxUnavailable":1,"maxSurge":"25%"}},"revisionHistoryLimit":10,"progressDeadlineSeconds":600},"status":{"observedGeneration":1,"replicas":2,"updatedReplicas":2,"readyReplicas":1,"availableReplicas":1,"conditions":[{"type":"Available","status":"True","lastUpdateTime":"2023-05-11T02:50:15Z","lastTransitionTime":"2023-05-11T02:50:15Z","reason":"MinimumReplicasAvailable","message":"Deployment has minimum availability."},{"type":"Progressing","status":"True","lastUpdateTime":"2023-05-11T02:50:16Z","lastTransitionTime":"2023-05-11T02:50:06Z","reason":"NewReplicaSetAvailable","message":"ReplicaSet \"weave-64897985d\" has successfully progressed."}]}}
        ]}
    },
    {   
        path: /\/api\/v1\/namespaces\/.*\/pods/,
        method: "get",
        data: { items: [
            {
                "apiVersion": "v1",
                "kind": "Pod",
                "metadata": {
                    "name": "weave-64897985d-lxcrj",
                    "namespace": "demo",
                },
                "spec": {
                    "containers": [
                        {
                            "image": "k8s.gcr.io/weave/weave:v1.0.0",
                            "name": "weave"
                        }
                    ],
                    "dnsPolicy": "Default",
                    "restartPolicy": "Always",
                    "schedulerName": "default-scheduler"
                },
                "status": {
                    "phase": "Running",
                }
            },
            {
                "apiVersion": "v1",
                "kind": "Pod",
                "metadata": {
                    "name": "weave-64897985d-he27c",
                    "namespace": "demo",
                },
                "spec": {
                    "containers": [
                        {
                            "image": "k8s.gcr.io/weave/weave:v1.0.0",
                            "name": "weave"
                        }
                    ],
                    "dnsPolicy": "Default",
                    "restartPolicy": "Always",
                    "schedulerName": "default-scheduler"
                },
                "status": {
                    "phase": "Running",
                }
            },
        ]}
    },
    {
        path: '/api/v1/posts',
        method: "get",
        data: [
            {"id":1,"name":"First Post","content":"","summary":"","creatorId":1,"creator":{"id":1,"name":"admin","email":"admin@weave.com","avatar":"","authInfos":null,"groups":null,"roles":null,"createdAt":"2023-05-11T01:38:23.483794Z","updatedAt":"0001-01-01T00:00:00Z"},"tags":[],"categories":[],"comments":null,"views":23,"likes":0,"userLiked":false,"createdAt":"2023-05-11T01:38:23.533917Z","updatedAt":"0001-01-01T00:00:00Z"},
            {"id":1,"name":"Weave Document","content":"","summary":"","creatorId":1,"creator":{"id":1,"name":"admin","email":"admin@weave.com","avatar":"","authInfos":null,"groups":null,"roles":null,"createdAt":"2023-05-11T01:38:23.483794Z","updatedAt":"0001-01-01T00:00:00Z"},"tags":[],"categories":[],"comments":null,"views":23,"likes":0,"userLiked":false,"createdAt":"2023-05-11T01:38:23.533917Z","updatedAt":"0001-01-01T00:00:00Z"}
        ]
    },
    {
        path: /\/api\/v1\/posts\/.*/,
        method: "get",
        data: {
            id: 1,
            createdAt: "2023-04-28T01:38:23.500389Z",
            name: "Weave Post Demo",
            creator: {"id":1,"name":"admin"},
            views: "108",
            likes: "16",
            comments: [],
            content: "\n## h2 Heading\n### h3 Heading\n#### h4 Heading\n##### h5 Heading\n###### h6 Heading\n\n## Emphasis\n\n**This is bold text**\n\n__This is bold text__\n\n*This is italic text*\n\n_This is italic text_\n\n~~Strikethrough~~\n\n\n## Horizontal Rules\n\n___\n\n---\n\n***\n\n\n## Typographic replacements\n\nEnable typographer option to see result.\n\n(c) (C) (r) (R) (tm) (TM) (p) (P) +-\n\ntest.. test... test..... test?..... test!....\n\n!!!!!! ???? ,,  -- ---\n\n\"Smartypants, double quotes\" and single quotes\n\n\n## Blockquotes\n\n\n\u003e Blockquotes can also be nested...\n\u003e\u003e ...by using additional greater-than signs right next to each other...\n\u003e \u003e \u003e ...or with spaces between arrows.\n\n\n## Lists\n\nUnordered\n\n+ Create a list by starting a line with `+`, `-`, or `*`\n+ Sub-lists are made by indenting 2 spaces:\n  - Marker character change forces new list start:\n    * Ac tristique libero volutpat at\n    + Facilisis in pretium nisl aliquet\n    - Nulla volutpat aliquam velit\n+ Very easy!\n\nOrdered\n\n1. Lorem ipsum dolor sit amet\n2. Consectetur adipiscing elit\n3. Integer molestie lorem at massa\n\n\n1. You can use sequential numbers...\n1. ...or keep all the numbers as `1.`\n\nStart numbering with offset:\n\n57. foo\n1. bar\n\n\n## Code\n\nInline `code`\n\nIndented code\n\n    // Some comments\n    line 1 of code\n    line 2 of code\n    line 3 of code\n\n\nBlock code \"fences\"\n\n```\nSample text here...\n```\n\nSyntax highlighting\n\n``` js\nvar foo = function (bar) {\n  return bar++;\n};\n\nconsole.log(foo(5));\n```\n\n## Tables\n\n| Option | Description |\n| ------ | ----------- |\n| data   | path to data files to supply the data that will be passed into templates. |\n| engine | engine to be used for processing templates. Handlebars is the default. |\n| ext    | extension to be used for dest files. |\n\nRight aligned columns\n\n| Option | Description |\n| ------:| -----------:|\n| data   | path to data files to supply the data that will be passed into templates. |\n| engine | engine to be used for processing templates. Handlebars is the default. |\n| ext    | extension to be used for dest files. |\n\n\n## Links\n\n[link text](http://qingwave.github.io)\n\n[link with title](http://qingwave.github.io \"title text!\")\n\nAutoconverted link https://qingwave.github.io (enable linkify to see)\n\n"
        }
    },
    {
        path: /\/api\/v1\/.*/,
        method: "",
        data: {}
    }
]

export default data;
