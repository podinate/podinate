package cmd

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/mail"

	"github.com/johncave/podinate/api-backend/account"
	"github.com/johncave/podinate/api-backend/config"
	"github.com/johncave/podinate/api-backend/iam"
	lh "github.com/johncave/podinate/api-backend/loghandler"
	"github.com/johncave/podinate/api-backend/user"
	"gopkg.in/yaml.v3"

	api "github.com/johncave/podinate/api-backend/go"

	"github.com/matthewhartstonge/argon2"
	"github.com/spf13/cobra"
	"go.tmthrgd.dev/passit"
)

var email string
var ip string
var dnsNames []string

func init() {
	initCmd.PersistentFlags().StringVar(&email, "email", "", "Email address of the initial admin user")
	initCmd.PersistentFlags().StringVar(&ip, "ip", "", "IP addresses to add to server certificate")
	//initCmd.MarkFlagRequired("email")
	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialise the configuration",
	Long:  `Initial configuration of the podinate server. This will create an initial administrator user and print out the crdeentials.`,
	// Disable default help
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {

		if !validMailAddress(email) {
			return errors.New("Invalid email address")
		}

		if ip == "" {
			return errors.New("IP address required")
		}

		// //Check if any users exist in the database
		// var uuid string
		// err := config.DB.QueryRow("SELECT uuid FROM \"user\" LIMIT 1").Scan(&uuid)
		// if err == nil {
		// 	return errors.New("Configuration already initialised")
		// }
		// if err != sql.ErrNoRows {
		// 	return err
		// }

		// Generate a random username in the format word-word8888
		username, err := passit.Repeat(passit.EFFLargeWordlist, "-", 2).Password(rand.Reader)
		if err != nil {
			return err
		}
		code, err := passit.Repeat(passit.Digit, "", 4).Password(rand.Reader)
		if err != nil {
			return err
		}
		username = username + code

		fmt.Print("Username: ", username, "\n")

		password, err := passit.Repeat(passit.EFFLargeWordlist, " ", 5).Password(rand.Reader)
		if err != nil {
			lh.Log.Fatalw("Error generating initial administrator password", "error", err)
		}
		argon := argon2.DefaultConfig()
		hash, err := argon.HashEncoded([]byte(password))
		if err != nil {
			lh.Log.Fatalw("Error hashing default administraor password", "error", err)
		}
		store := base64.StdEncoding.EncodeToString(hash)
		var adminID string
		err = config.DB.QueryRow("INSERT INTO \"user\" (id, display_name, email, password_hash) VALUES ($1, 'Administrator', $2, $3) returning uuid", username, email, store).Scan(&adminID)
		if err != nil {
			lh.Log.Fatalw("Error creating initial administrator account", "error", err)
		}

		lh.Log.Infow("Created initial administrator account", "username", "administrator", "password", password)

		// Issue the new user an api key
		u, err := user.GetByUUID(adminID)
		if err != nil {
			return err
		}

		// Create the `default` account
		defaultAccount := api.Account{
			Id:   "default",
			Name: "default",
		}
		lh.Log.Debugw("Creating default account", "account", defaultAccount, "owner", u, "error", err)
		newAcc, apierr := account.Create(defaultAccount, u)
		if apierr != nil {
			return apierr
		}

		// Add initial policies to the account
		superAdminPolicyDocument := `
version: 2023.1
statements:
	- effect: allow
	actions: ["**"]
	resources: ["**"]`
		superAdminPolicy, err := iam.CreatePolicyForAccount(&newAcc, "super-administrator", superAdminPolicyDocument, "Default policy created during initial account creation")
		apierr = superAdminPolicy.AttachToRequestor(u, u)
		if apierr != nil {
			// We can pass this error directly to the API response
			lh.Log.Fatalw("Error attaching super-administrator policy to initial default account", "error", apierr)
			return apierr
		}

		apiKey, err := u.IssueAPIKey("Initial Credentials")
		if err != nil {
			lh.Log.Fatalw("Error generating initial administrator API key", "error", err)
		}

		profile := struct {
			Name   string `yaml:"name"`
			ApiKey string `yaml:"api_key"`
			ApiUrl string `yaml:"api_url"`
		}{
			Name:   "default",
			ApiKey: apiKey,
			ApiUrl: "http://" + ip + ":31443",
		}

		// write the yaml to a file
		yamlData, err := yaml.Marshal(profile)
		if err != nil {
			lh.Log.Fatalw("Error marshaling profile yaml", "error", err)
		}
		err = ioutil.WriteFile("/profile.yaml", yamlData, 0644)
		if err != nil {
			lh.Log.Fatalw("Error writing profile yaml", "error", err)
		}

		log.Println("Great success, profile generated")

		// Intend to use https keypair security eventually but for alpha this will do

		// Create a CA keypair
		// caSerial, err := generateSerial()
		// if err != nil {
		// 	return err
		// }
		// ca := &x509.Certificate{
		// 	SerialNumber:          caSerial,
		// 	Subject:               pkix.Name{CommonName: "Podinate CA"},
		// 	NotBefore:             time.Now(),
		// 	NotAfter:              time.Now().AddDate(10, 0, 0),
		// 	IsCA:                  true,
		// 	KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign | x509.KeyUsageDigitalSignature,
		// 	BasicConstraintsValid: true,
		// }

		// caPubKey, caPrivateKey, err := ed25519.GenerateKey(rand.Reader)
		// if err != nil {
		// 	return err
		// }
		// caBytes, err := x509.CreateCertificate(rand.Reader, ca, ca, caPubKey, caPrivateKey)

		// // Generate the server certificate
		// serverSerial, err := generateSerial()
		// if err != nil {
		// 	return err
		// }
		// serverCert := x509.Certificate{
		// 	SerialNumber: serverSerial,
		// 	Subject:      pkix.Name{CommonName: "Podinate Controller"},

		// 	NotBefore:   time.Now(),
		// 	NotAfter:    time.Now().AddDate(1, 0, 0),
		// 	KeyUsage:    x509.KeyUsageDigitalSignature,
		// 	ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		// }

		// err = config.DB.QueryRow("INSERT INTO \"user\" (id, display_name, email) VALUES ($1, $2, $3) returning uuid", username, "Administrator", email).Scan(&uuid)
		// if err != nil {
		// 	return err
		// }

		return nil
	},
}

func validMailAddress(address string) bool {
	// if address == nil {
	// 	return false
	// }
	_, err := mail.ParseAddress(address)
	if err != nil {
		return false
	}
	return true
}

func generateSerial() (*big.Int, error) {
	bi := big.NewInt(0)
	bi.UnmarshalText([]byte("0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF"))
	return rand.Int(rand.Reader, bi)
}
