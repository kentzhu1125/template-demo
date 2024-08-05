pipeline {
  agent {
    node {
      label 'go'
    }

  }
  stages {
    stage('build & push') {
      agent none
      steps {
        container('go') {
          sh '''export GO111MODULE=on
export GOPROXY=https://goproxy.cn
go build -o app -v -a  main.go'''
          sh '''export TAG=master-$BUILD_NUMBER
docker login -u "${DOCKER_USERNAME}" -p "${DOCKER_PASSWORD}" "${REGISTRY}"

docker build -t $REGISTRY/$DOCKER_NAMESPACE/apidemo:$TAG  .
docker push $REGISTRY/$DOCKER_NAMESPACE/apidemo:$TAG'''
        }

      }
    }

    stage('审核') {
      agent none
      steps {
        input(message: ' 部署到测试环境? @admin   ', submitter: 'admin')
      }
    }

    stage('deploy-prod') {
      agent none
      steps {
        container('go') {
          withCredentials([kubeconfigContent(credentialsId : 'kubeconfig' ,variable : 'KUBE_CONFIG' ,)]) {
            sh '''mkdir $HOME/.kube && echo "${KUBE_CONFIG}" > $HOME/.kube/config
export TAG=master-$BUILD_NUMBER
export NAMESPACE="prod-apidemo"
pwd && ls -ltr
envsubst < yaml/deployment.yaml | kubectl apply -f -'''
          }

        }

      }
    }

    stage('send success email') {
      agent none
      steps {
        mail(to: 'xiezw@csxjy.com.cn', subject: '部署生产环境成功', body: '部署生产环境成功')
      }
    }

  }
  environment {
    DOCKER_USERNAME = 'admin'
    DOCKER_PASSWORD = 'Galaxyclouds2021'
    REGISTRY = 'harbor.galaxyclouds.com'
    DOCKER_NAMESPACE = 'demo'
  }
}
