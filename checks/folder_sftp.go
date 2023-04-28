package checks

import (
	"fmt"

	"github.com/flanksource/canary-checker/api/context"
	v1 "github.com/flanksource/canary-checker/api/v1"
	"github.com/flanksource/canary-checker/pkg"
	"github.com/flanksource/canary-checker/pkg/db"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

func CheckSFTP(ctx *context.Context, check v1.FolderCheck) pkg.Results {
	result := pkg.Success(check, ctx.Canary)
	var results pkg.Results
	results = append(results, result)

	k8sClient, err := ctx.Kommons.GetClientset()
	if err != nil {
		return results.Failf("error getting k8s client from kommons client: %v", err)
	}

	foundConn, err := check.SFTPConnection.PopulateFromConnection(ctx, db.Gorm, k8sClient, ctx.Namespace)
	if err != nil {
		return results.Failf("failed to populate SFTP connection: %v", err)
	}

	auth := check.SFTPConnection.Auth
	if !foundConn {
		auth, err = GetAuthValues(check.SFTPConnection.Auth, ctx.Kommons, ctx.Canary.Namespace)
		if err != nil {
			return results.ErrorMessage(err)
		}
	}

	conn, err := sshConnect(check.SFTPConnection.Host, check.SFTPConnection.GetPort(), auth.GetUsername(), auth.GetPassword())
	if err != nil {
		return results.ErrorMessage(err)
	}
	defer conn.Close()

	client, err := sftp.NewClient(conn)
	if err != nil {
		return results.ErrorMessage(err)
	}
	defer client.Close()

	session := Filesystem(client)
	folders, err := getGenericFolderCheck(session, check.Path, check.Filter)
	if err != nil {
		return results.ErrorMessage(err)
	}
	result.AddDetails(folders)

	if test := folders.Test(check.FolderTest); test != "" {
		return results.Failf(test)
	}
	return results
}

func sshConnect(host string, port int, user string, password string) (*ssh.Client, error) {
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	return ssh.Dial("tcp", fmt.Sprintf("%s:%d", host, port), config)
}
