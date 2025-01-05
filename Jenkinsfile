pipeline {
    agent any
    environment {
        DOCKER_REGISTRY = 'jhedie'
        FRONTEND_IMAGE = "${DOCKER_REGISTRY}/word-game-frontend"
        BACKEND_IMAGE = "${DOCKER_REGISTRY}/word-game-backend"
    }
    stages {
        stage('Checkout') {
            steps {
                checkout scm
            }
        }
        stage('Determine Changes') {
            steps {
                script {
                    echo 'Test'
                    def deploymentOption = input(
                        id: 'deploymentOption', message:'Select the parts to deploy:', parameters: [
                            choice(name: 'DEPLOY_PART', choices: ['Frontend', 'Backend', 'Both'], description: 'Choose which part to deploy')
                        ]
                    )
                    env.DEPLOY_PART = deploymentOption
                }
            }
        }
        stage('Deploy Frontend') {
            when {
                expression { env.DEPLOY_PART == 'Frontend' || env.DEPLOY_PART == 'Both' }
            }
            steps {
                dir('Frontend') {
                    withCredentials([string(credentialsId: '4435865f-e00c-4a92-a656-eaec1939a9da', variable: 'DOCKER_CREDENTIALS')]) {
                        sh(script: "echo $DOCKER_CREDENTIALS")
                        sh 'docker build -t ${FRONTEND_IMAGE}:${BUILD_NUMBER} .'
                    }
                }
            }
        }
        stage('Deploy Backend') {
            when {
                expression { env.DEPLOY_PART == 'Backend' || env.DEPLOY_PART == 'Both' }
            }
            steps {
                dir('Backend') {
                    withCredentials([string(credentialsId: '4435865f-e00c-4a92-a656-eaec1939a9da', variable: 'DOCKER_CREDENTIALS')]) {
                        sh(script: "echo $DOCKER_CREDENTIALS")
                        sh 'docker build -t ${BACKEND_IMAGE}:${BUILD_NUMBER} .'
                    }
                }
            }
        }
        stage('Push Images') {
            steps {
                withCredentials([string(credentialsId: '4435865f-e00c-4a92-a656-eaec1939a9da', variable: 'DOCKER_CREDENTIALS')]) {
                    sh(script: "echo $DOCKER_CREDENTIALS")
                    script {
                        if (env.DEPLOY_PART == 'Frontend' || env.DEPLOY_PART == 'Both') {
                            sh 'docker push ${FRONTEND_IMAGE}:${BUILD_NUMBER}'
                        }
                        if (env.DEPLOY_PART == 'Backend' || env.DEPLOY_PART == 'Both') {
                            sh 'docker push ${BACKEND_IMAGE}:${BUILD_NUMBER}'
                        }
                    }
                }
            }
        }
        stage('Deploy Updated Containers') {
            steps {
                withCredentials([string(credentialsId: '4435865f-e00c-4a92-a656-eaec1939a9da', variable: 'DOCKER_CREDENTIALS')]) {
                    sh(script: "echo $DOCKER_CREDENTIALS")
                    script {
                        if (env.DEPLOY_PART == 'Frontend' || env.DEPLOY_PART == 'Both') {
                            sh 'kubectl set image deployment/frontend frontend=${FRONTEND_IMAGE}:${BUILD_NUMBER}'
                        }
                        if (env.DEPLOY_PART == 'Backend' || env.DEPLOY_PART == 'Both') {
                            sh 'kubectl set image deployment/backend backend=${BACKEND_IMAGE}:${BUILD_NUMBER}'
                        }
                    }
                }
            }
        }
        post {
            always {
                cleanWs(disableDeferredWipeout: true)
            }
        }
        }
    }
