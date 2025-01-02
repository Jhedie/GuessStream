pipeline {
    agent any
    stages {
        stage ('Determine Changes') {
            steps {
                script {
                    echo "Test"
                    def deploymentOption = input(
                        id: 'deploymentOption', message:'Select the parts to deploy:', parameters: [
                            choice(name: 'DEPLOY_PART', choices: ['Frontend', 'Backend', 'Both'], description: 'Choose which part to deploy')
                        ]
                    )
                    env.DEPLOY_PART = deploymentOption
                }
            }
        }
        stage ('Deploy Frontend') {
            when {
                expression { env.DEPLOY_PART == 'Frontend' || env.DEPLOY_PART == 'Both' }
            }
            steps {
                dir('Frontend') {
                    sh(script: "ls")
                }
            }
        }
        stage ('Deploy Backend') {
            when {
                expression { env.DEPLOY_PART == 'Backend' || env.DEPLOY_PART == 'Both' }
            }
            steps {
                dir('Backend') {
                    sh(script: "ls")
                }
            }
        }
    }
}
