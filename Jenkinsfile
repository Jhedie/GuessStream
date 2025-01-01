pipeline {
    agent any
    stages {
        stage ('Determine Changes') {
            steps {
                script {
                    echo "Test"
                    def changes = sh(script: "git diff --name-only HEAD~1", returnStdout: true).trim()
                    env.FRONTEND_CHANGED = changes.contains('frontend/')
                    env.BACKEND_CHANGED = changes.contains('backend/')
                }
            }
        }
        stage ('Deploy Frontend') {
            when {
                expression {env.FRONTEND_CHANGED == 'true'}
            }
            steps {
                dir('Frontend')
                sh(script: "ls")
            }
        }
        stage ('Deploy Backend') {
            when {
                expression {env.BACKEND_CHANGED == 'true'}
            }
            steps {
                dir('Frontend')
                sh(script: "ls")
            }
        }
    }
}
