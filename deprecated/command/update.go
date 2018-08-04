package command

import (
	"archive/zip"
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"runtime"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/gen0cide/laforge/deprecated/competition"
	"github.com/google/go-github/github"
	update "github.com/inconshreveable/go-update"
)

func CmdUpdate(c *cli.Context) {
	ctx := context.Background()

	client := github.NewClient(nil)

	repRel, _, err := client.Repositories.GetLatestRelease(ctx, "gen0cide", "laforge")
	if err != nil {
		competition.LogFatal("Github Error: " + err.Error())
	}

	tagName := strings.TrimLeft(repRel.GetTagName(), "v")

	if c.App.Version == tagName {
		competition.Log("Running Lastest Version (no update available): " + c.App.Version)
		return
	}

	competition.Log("New Version Found: " + tagName)

	fileName := fmt.Sprintf("laforge_%s_%s_amd64.zip", tagName, runtime.GOOS)

	assetID := 0
	fileSize := 0

	for _, asset := range repRel.Assets {
		if asset.GetName() == fileName {
			assetID = asset.GetID()
			fileSize = asset.GetSize()
			competition.Log(fmt.Sprintf("Downloading New Binary: %s", asset.GetBrowserDownloadURL()))
		}
	}

	if assetID == 0 {
		competition.LogFatal("No release found for your OS! WTF? Report this.")
	}

	_, redirectURL, err := client.Repositories.DownloadReleaseAsset(ctx, "gen0cide", "laforge", assetID)
	if err != nil {
		competition.LogFatal("Github Error: " + err.Error())
	}

	if len(redirectURL) == 0 {
		competition.LogFatal("There was an error retriving the release from Github.")
	}

	resp, err := http.Get(redirectURL)
	if err != nil {
		competition.LogFatal("Error Retrieving Binary: " + err.Error())
	}
	defer resp.Body.Close()
	compressedFile, err := ioutil.ReadAll(resp.Body)

	competition.Log("Uncompressing Binary...")

	var binary bytes.Buffer
	writer := bufio.NewWriter(&binary)
	compressedReader := bytes.NewReader(compressedFile)

	r, err := zip.NewReader(compressedReader, int64(fileSize))
	if err != nil {
		competition.LogFatal("Error Buffering Zip File: " + err.Error())
	}

	for _, zf := range r.File {
		if zf.Name != "laforge" {
			continue
		}
		src, err := zf.Open()
		if err != nil {
			competition.LogFatal("Error Unzipping File: " + err.Error())
		}
		defer src.Close()

		io.Copy(writer, src)
	}

	reader := bufio.NewReader(&binary)

	err = update.Apply(reader, update.Options{})
	if err != nil {
		if rerr := update.RollbackError(err); rerr != nil {
			competition.LogFatal(fmt.Sprintf("Failed to rollback from bad update: %v", rerr))
		}
		competition.LogFatal("Update Failed - original version rolled back successfully.")
	}

	competition.Log("Successfully updated to laforge v" + tagName)

	return
}
