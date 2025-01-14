pipeline {
    agent { label 'kthw-agent' }
    environment {
        JENKINS_NODE_COOKIE = 'do_not_kill'
        ANSIBLE_PRIVATE_KEY = credentials('admin-agent')
    }
    stages {
        stage('Run kubectl-access.yml') {
            steps {
                script {
                    // Run another playbook
                    def result = sh(script: "ansible-playbook -i hosts.hosts --private-key=$ANSIBLE_PRIVATE_KEY kubectl-access.yml", returnStatus: true)

                    // Check if the third playbook ran successfully
                    if (result != 0) {
                        currentBuild.result = 'FAILURE'
                        error "kubectl-access.yml failed."
                    }
                }
            }
        }
    }

    post {
        success {
            echo "All playbooks executed successfully!"
        }
        failure {
            echo "One or more playbooks failed."
        }
    }
}
