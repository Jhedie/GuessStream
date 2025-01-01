pipeline {
    agent any
    stages {
        stage ('Determine Changes') {
            steps {
                script {
                    echo "Test"
                    def changes = sh(script: "git diff --name-only HEAD~1", returnStdout: true).trim()
                    env.FRONTEND_CHANGED = changes.contains('Frontend/')
                    env.BACKEND_CHANGED = changes.contains('Backend/')
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
                dir('Backend')
                sh(script: "ls")
            }
        }
    }
}
