# rapGO.io : an experimental platform for generating a rap from any input voice.

When I initially started this project, I wanted to dig deeper in the audio manipulation of sounds. 
Then I wanted to add a scalable infrastructure layer using Kubernetes and a distributed filesystem.

This project is a kind of training for deploying a complex kubernetes application with several microservices
powered by several languages like Go, C, Python and Reactjs for the frontend services.

For now, the `src/rapgeneratorservice/lib` which is the core of the application is a bit simple and the results are very far from being audible. However, I'm improving this core library over the time to create a better user experience.

To learn more about the architecture of rapGO.io, please refer to the [docs](/docs/README.md) section.

## CephFS instructions (obselete)

`git clone https://github.com/rook/rook.git`
`cd k8s/rook/cluster/example/kubernetes/ceph`
`kubectl create -f common.yaml`
`kubectl create -f operator.yaml`

Wait that the `rook-operator``is in running state and then :
`kubectl create -f cluster-test.yaml`
Once everything is running, you can deploy the `rook-toolbox` with :
`kubectl create -f toolbox.yaml`
Once the toolbox is running, do :
`kubectl -n rook-ceph exec -it $(kubectl -n rook-ceph get pod -l "app=rook-ceph-tools" -o jsonpath='{.items[0].metadata.name}') bash` to enter the pod and launch the command : `ceph status` to see if the ceph cluster is healthy. 

Then, you may want to see this in a nice dashboard. So you can start a kubernetes proxy with `kubectl proxy`, go in your browser and type in `http://localhost:8001/api/v1/namespaces/rook-ceph/services/https:rook-ceph-mgr-dashboard:https-dashboard/proxy/#/login`

After you connect to the dashboard you will need to login for secure access. Rook creates a default user named `admin` and generates a secret called `rook-ceph-dashboard-admin-password` in the namespace where the Rook Ceph cluster is running. To retrieve the generated password, you can run the following:
`kubectl -n rook-ceph get secret rook-ceph-dashboard-password -o jsonpath="{['data']['password']}" | base64 --decode && echo`

## google cloud bucket

It was so painful to deploy glusterFS or cephFS that I decided I will setup a simple google cloud bucket, and thus add an extra bill for this simpler solution...

## Kafka

* https://success.docker.com/article/getting-started-with-kafka  (for testing)
* https://rmoff.net/2018/08/02/kafka-listeners-explained/ (an other interesting link for debugging kafka networking)
* kafka official Helm chart (for production)