# Contributing
### Making changes to fork 
Set the origin to our private fork and change to the branch you want to work on. 
- `cd ~/go/src/github.com/gen0cide/laforge`
- `git remote set-url origin https://github.com/jace0x21/laforge.git`
- `git checkout <branch>`
#####  If you make changes to `~/go/src/github.com/gen0cide/laforge`
Make sure you're on the right branch and rebuild the project and install.
- `git checkout <branch>`
- `go build -a && go install`
##### If you want to grab latest changes from fork
- `cd ~/go/src/github.com/gen0cide/laforge`
- `git pull`

