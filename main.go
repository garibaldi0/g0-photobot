package main

import (
	"context"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/davecgh/go-spew/spew"
	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox/files"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

type pbConfig struct {
	AppKey  string
	AppSec  string
	LogFile string
	DestDir string
	Tokens  map[string]string
}

var log = logrus.New()
var config dropbox.Config
var myConfig pbConfig

var mainCmd = &cobra.Command{
	Use:   "g0-photobot",
	Short: "A CLI tool for moving Camera Uploads to a shared folder",
	Long:  "Use PhotoBot to move your Camera Uploads to a shared folder",
	RunE:  initDbx,
}

func init() {
	flags := mainCmd.Flags()
	flags.BoolP("verbose", "v", false, "Enable verbose logging")
	flags.BoolP("new", "n", false, "Generate new user auth code")

	viper.BindPFlag("verbose", flags.Lookup("verbose"))
	viper.BindPFlag("new", flags.Lookup("new"))

	home, errHome := homedir.Dir()
	if errHome != nil {
		fmt.Println("Couldn't determine Home Directory")
	} else {
		homePath := path.Join(home, "go", "etc")
		viper.AddConfigPath(homePath)
	}
	viper.AddConfigPath("/etc/")   // set the path of your config file
	viper.AddConfigPath("../etc/") // set the path of your config file
	viper.AddConfigPath(".")       // optionally look for config in the working directory

	viper.SetConfigName("g0-photobot")

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if errConfig := viper.ReadInConfig(); errConfig == nil {
		defer log.Println("Using config file:", viper.ConfigFileUsed())
	}
	errUnmarshal := viper.Unmarshal(&myConfig)
	if errUnmarshal != nil {
		fmt.Println("Unable to read config file.")
		os.Exit(1)
	}
	if viper.GetBool("verbose") {
		spew.Dump(myConfig)
	}

	f, errOpen := os.OpenFile(myConfig.LogFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if errOpen != nil {
		fmt.Printf("Failed to open logfile : %s.\n", myConfig.LogFile)
		os.Exit(1)
	}
	log.Out = f
	log.Formatter = &logrus.TextFormatter{
		FullTimestamp:    true,
		DisableTimestamp: false,
	}
}

func main() {
	log.Info("Start Run")
	mainCmd.Execute()
	log.Info("End Run")
}

func validatePath(p string) (path string, err error) {
	path = p

	if !strings.HasPrefix(path, "/") {
		path = fmt.Sprintf("/%s", path)
	}

	path = strings.TrimSuffix(path, "/")

	return
}

func makeRelocationArg(s string, d string) (arg *files.RelocationArg, err error) {
	src, err := validatePath(s)
	if err != nil {
		return
	}
	dst, err := validatePath(d)
	if err != nil {
		return
	}

	arg = files.NewRelocationArg(src, dst)

	return
}

func initDbx(cmd *cobra.Command, args []string) (err error) {
	verbose, _ := cmd.Flags().GetBool("verbose")
	newUser, _ := cmd.Flags().GetBool("new")

	conf := oauth2.Config{
		ClientID:     myConfig.AppKey,
		ClientSecret: myConfig.AppSec,
		Endpoint:     dropbox.OAuthEndpoint(""),
	}

	if len(myConfig.Tokens) == 0 || newUser {
		fmt.Printf("1. Go to %v\n", conf.AuthCodeURL("state"))
		fmt.Printf("2. Click \"Allow\" (you might have to log in first).\n")
		fmt.Printf("3. Copy the authorization code.\n")
		fmt.Printf("Enter the authorization code here: ")

		var code string
		if _, err = fmt.Scan(&code); err != nil {
			return
		}
		var token *oauth2.Token
		ctx := context.Background()
		token, err = conf.Exchange(ctx, code)
		if err != nil {
			return
		}
		fmt.Printf("Add this token to the config file : %s\n", token.AccessToken)
		os.Exit(0)
	}

	logLevel := dropbox.LogOff
	if verbose {
		logLevel = dropbox.LogInfo
	}
	for uName, uToken := range myConfig.Tokens {
		log.Printf("Checking %s's Camera Uploads", uName)
		config = dropbox.Config{uToken, logLevel, nil, "", "", nil, nil, nil}
		mv()
	}

	return
}

// Sends a get_metadata request for a given path and returns the response
func getFileMetadata(c files.Client, path string) (files.IsMetadata, error) {
	arg := files.NewGetMetadataArg(path)

	res, err := c.GetMetadata(arg)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func mv() error {
	var destination = myConfig.DestDir
	var source = "/Camera Uploads"
	var mvErrors []error

	dbx := files.New(config)
	arg := files.NewListFolderArg(source)
	res, err := dbx.ListFolder(arg)
	var entries []files.IsMetadata
	if err != nil {
		switch e := err.(type) {
		case files.ListFolderAPIError:
			if e.EndpointError.Path.Tag == files.LookupErrorNotFolder {
				var metaRes files.IsMetadata
				metaRes, err = getFileMetadata(dbx, source)
				entries = []files.IsMetadata{metaRes}
			} else {
				return err
			}
		default:
			return err
		}
		if err != nil {
			return err
		}
	} else {
		entries = res.Entries
		for res.HasMore {
			arg := files.NewListFolderContinueArg(res.Cursor)
			res, err = dbx.ListFolderContinue(arg)
			if err != nil {
				return err
			}
			entries = append(entries, res.Entries...)
		}
	}
	for _, entry := range entries {
		switch f := entry.(type) {
		case *files.FileMetadata:
			log.Println(f.Metadata.Name)
			mvArg, err := makeRelocationArg(
				source+"/"+f.Metadata.Name,
				destination+"/"+f.Metadata.Name)
			if err != nil {
				relocationError := fmt.Errorf("Error validating move for %s: %v", f.Metadata.Name, err)
				mvErrors = append(mvErrors, relocationError)
			}
			if _, err := dbx.MoveV2(mvArg); err != nil {
				moveError := fmt.Errorf("Move error: %v", arg)
				mvErrors = append(mvErrors, moveError)
			}
		}
	}

	for _, mvError := range mvErrors {
		log.Errorf("%v\n", mvError)
	}

	return nil
}
