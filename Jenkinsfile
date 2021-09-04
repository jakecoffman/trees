pipeline {
    agent any

    stages {
        stage('Build') {
            steps {
                sh '''
                cd server
                go build server.go
                '''
                nodejs(nodeJSInstallationName: '13') {
                    sh '''
                    cd ui
                    npm ci
                    npm run build
                    '''
                }
            }
        }
        stage('Deploy') {
            steps {
                sh '''
scp server/server deploy@stldevs.com:~
scp -r ui/dist deploy@stldevs.com:~
ssh deploy@stldevs.com << EOF
   sudo service trees stop
   mv -f ~/server /opt/trees/server
   rm -rf /opt/trees/dist
   mv ~/dist /opt/trees
   cd /opt/trees
   chmod +x trees
   sudo service trees start
'''
            }
        }
    }
}
