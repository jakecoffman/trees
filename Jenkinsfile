pipeline {
    agent any

    stages {
        stage('Build Server') {
            when { changeset "server/*"}
            steps {
                sh '''
                cd server
                go build server.go
                '''
            }
        }
        stage('Build UI') {
            when { changeset "ui/*"}
            steps {
                nodejs(nodeJSInstallationName: '13') {
                    sh '''
                        cd ui
                        npm ci
                        npm run build
                        '''
                    }
                }
            }
        }
        stage('Deploy Server') {
            when { changeset "server/*"}
            steps {
                sh '''
scp server/server deploy@stldevs.com:~
ssh deploy@stldevs.com << EOF
   sudo service trees stop
   mv -f ~/server /opt/trees/server
   cd /opt/trees
   chmod +x server
   sudo service trees start
'''
            }
        }
        stage('Deploy UI') {
            when { changeset "ui/*"}
            steps {
                sh '''
scp -r ui/dist deploy@stldevs.com:~
ssh deploy@stldevs.com << EOF
   rm -rf /opt/trees/dist
   mv ~/dist /opt/trees
'''
        }
    }
}
