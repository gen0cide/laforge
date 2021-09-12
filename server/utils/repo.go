package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	gssh "github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"golang.org/x/crypto/ssh"
)

func CloneGit(repoURL, repoPath, privateKey, branchName string) (string, error) {

	_, err := os.Stat(privateKey)
	if err != nil {
		err := fmt.Errorf("read file %s failed %s", privateKey, err.Error())
		return "", err
	}

	publicKeys, err := gssh.NewPublicKeysFromFile("git", privateKey, "")

	if err != nil {
		err := fmt.Errorf("generate publickeys failed: %s", err.Error())
		return "", err
	}

	branch := fmt.Sprintf("refs/heads/%s", branchName)
	repo, err := git.PlainClone(repoPath, false, &git.CloneOptions{
		Auth:          publicKeys,
		URL:           repoURL,
		ReferenceName: plumbing.ReferenceName(branch),
	})
	if err != nil {
		err := fmt.Errorf("unable to clone repo: %s", err.Error())
		return "", err
	}

	// Print the latest commit that was just pulled
	ref, err := repo.Head()
	if err != nil {
		err := fmt.Errorf("unable to access git head: %s", err.Error())
		return "", err
	}

	commit, err := repo.CommitObject(ref.Hash())
	if err != nil {
		err := fmt.Errorf("unable to get commit hash: %s", err.Error())
		return "", err
	}

	return commit.String(), err
}

func PullGit(repoPath, privateKey, branchName string) (string, error) {

	_, err := os.Stat(privateKey)
	if err != nil {
		err := fmt.Errorf("read file %s failed %s", privateKey, err.Error())
		return "", err
	}

	publicKeys, err := gssh.NewPublicKeysFromFile("git", privateKey, "")

	if err != nil {
		err := fmt.Errorf("generate publickeys failed: %s", err.Error())
		return "", err
	}

	// We instantiate a new repository targeting the given path (the .git folder)
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		err := fmt.Errorf("opening git repo failed: %s", err.Error())
		return "", err
	}
	// Get the working directory for the repository
	w, err := repo.Worktree()
	if err != nil {
		err := fmt.Errorf("getting git working directory failed: %s", err.Error())
		return "", err
	}
	branch := fmt.Sprintf("refs/heads/%s", branchName)
	if err = w.Pull(&git.PullOptions{
		ReferenceName: plumbing.ReferenceName(branch),
		SingleBranch:  true,
		Force:         true,
		Auth:          publicKeys,
	}); err != nil && err != git.NoErrAlreadyUpToDate {
		return "", err
	}

	// Print the latest commit that was just pulled
	ref, err := repo.Head()
	if err != nil {
		err := fmt.Errorf("unable to access git head: %s", err.Error())
		return "", err
	}

	commit, err := repo.CommitObject(ref.Hash())
	if err != nil {
		err := fmt.Errorf("unable to get commit hash: %s", err.Error())
		return "", err
	}

	return commit.String(), err
}

// MakeSSHKeyPair make a pair of public and private keys for SSH access.
// Public key is encoded in the format for inclusion in an OpenSSH authorized_keys file.
// Private Key generated is PEM encoded
func MakeSSHKeyPair(privateKeyPath string) error {
	_, fileCheck := os.Stat(privateKeyPath)
	if fileCheck == nil {
		return nil
	}
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return err
	}

	// generate and write private key as PEM
	privateKeyFile, err := os.Create(privateKeyPath)
	if err != nil {
		return err
	}
	defer privateKeyFile.Close()

	privateKeyPEM := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)}
	if err := pem.Encode(privateKeyFile, privateKeyPEM); err != nil {
		return err
	}

	// generate and write public key
	pub, err := ssh.NewPublicKey(&privateKey.PublicKey)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(privateKeyPath+".pub", ssh.MarshalAuthorizedKey(pub), 0655)
}
