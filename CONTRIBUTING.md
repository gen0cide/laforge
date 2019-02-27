# Contributing
### Making changes to fork
Make sure remote is set to fork: `git remote set-url origin https://github.com/jace0x21/laforge.git`
1. `cd ~/go/src/github.com/jace0x21/`
2. `rm -rf /laforge`
3. `git clone git@github.com:jace0x21/laforge.git`
4. Make changes to `~/go/src/github.com/jace0x21/laforge`
5. `cd ~/go/src/github.com/jace0x21/laforge/cmd/laforge`
6. `go build -a && go install`

