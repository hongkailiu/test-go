package cmd

import (
	"fmt"

	"strconv"

	"github.com/spf13/cobra"

	"github.com/hongkailiu/test-go/pkg/cmd/dgst"
	"github.com/hongkailiu/test-go/pkg/cmd/enc"
	"github.com/hongkailiu/test-go/pkg/cmd/genpkey"
	"github.com/hongkailiu/test-go/pkg/cmd/pkey"
	"github.com/hongkailiu/test-go/pkg/cmd/pkeyutl"
	"github.com/hongkailiu/test-go/pkg/cmd/rand"
	"github.com/hongkailiu/test-go/pkg/cmd/req"
	"github.com/hongkailiu/test-go/pkg/cmd/verify"
	"github.com/hongkailiu/test-go/pkg/cmd/x509"
)

var rootCmd = &cobra.Command{
	Use: "gossl",
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}

func Execute() error {
	return rootCmd.Execute()
}

var (
	genpkeyOption genpkey.Option
	pkeyOption    pkey.Option
	reqOption     req.Option
	x509Option    x509.Option
	verifyOption  verify.Option
	randOption    rand.Option
	pkeyutlOption pkeyutl.Option
	dgstOption    dgst.Option
	encOption     enc.Option
)

func init() {

	genpkeyCmd.Flags().StringVarP(&genpkeyOption.Algorithm, "algorithm", "", "", "The public key algorithm")
	genpkeyCmd.Flags().StringSliceVarP(&genpkeyOption.PKeyOpt, "pkeyopt", "", []string{}, "Set the public key algorithm option as opt:value")
	genpkeyCmd.Flags().StringVarP(&genpkeyOption.Out, "out", "", "", "Output (private key) file")

	pkeyCmd.Flags().StringVarP(&pkeyOption.In, "in", "", "", "Input key")
	pkeyCmd.Flags().BoolVarP(&pkeyOption.PubOut, "pubout", "", false, "Restrict encoded output to public components")
	pkeyCmd.Flags().BoolVarP(&pkeyOption.Text, "text", "", false, "Output key components in plaintext")
	pkeyCmd.Flags().BoolVarP(&pkeyOption.NoOut, "noout", "", false, "Do not output the key in encoded form")
	pkeyCmd.Flags().StringVarP(&pkeyOption.Out, "out", "", "", "Output file for encoded and/or text output")

	reqCmd.Flags().BoolVarP(&reqOption.New, "new", "", false, "New request")
	reqCmd.Flags().StringVarP(&reqOption.Key, "key", "", "", "Key for signing, and to include unless -in given")
	reqCmd.Flags().StringVarP(&reqOption.Out, "out", "", "", "Output file for encoded and/or text output")
	reqCmd.Flags().BoolVarP(&reqOption.Sha256, "sha256", "", false, "This specifies the message digest to sign the request")
	reqCmd.Flags().BoolVarP(&reqOption.X509, "x509", "", false, "This option outputs a self signed certificate instead of a certificate request")
	reqCmd.Flags().BoolVarP(&reqOption.NoEnc, "noenc", "", false, "Key for signing, and to include unless -in given")
	reqCmd.Flags().IntVarP(&reqOption.Days, "days", "", 30, "When the -x509 option is being used this specifies the number of days to certify the certificate for, otherwise it is ignored. n should be a positive integer. The default is 30 days.")

	x509Cmd.Flags().BoolVarP(&x509Option.Req, "req", "", false, "Input is a CSR file (rather than a certificate)")
	x509Cmd.Flags().BoolVarP(&x509Option.Sha256, "sha256", "", false, "This specifies the message digest to sign the request")
	x509Cmd.Flags().BoolVarP(&x509Option.Text, "text", "", false, "Output key components in plaintext")
	x509Cmd.Flags().BoolVarP(&x509Option.NoOut, "noout", "", false, "Do not output the key in encoded form")
	x509Cmd.Flags().IntVarP(&x509Option.Days, "days", "", 30, "Number of days until newly generated certificate expires")
	x509Cmd.Flags().StringVarP(&x509Option.CA, "CA", "", "", "Use the given CA certificate")
	x509Cmd.Flags().StringVarP(&x509Option.CAKey, "CAkey", "", "", "The corresponding CA key")
	x509Cmd.Flags().BoolVarP(&x509Option.CACreateSerial, "CAcreateserial", "", false, "Create CA serial number file if it does not exist")
	x509Cmd.Flags().StringVarP(&x509Option.In, "in", "", "", "Certificate input, or CSR input file with -req")
	x509Cmd.Flags().StringVarP(&x509Option.Out, "out", "", "", "Output file")
	x509Cmd.Flags().BoolVarP(&x509Option.PubKey, "pubkey", "", false, "Print the public key in PEM format")

	verifyCmd.Flags().StringVarP(&verifyOption.CAFile, "CAfile", "", "", "A file of trusted certificates")

	randCmd.Flags().BoolVarP(&randOption.Base64, "base64", "", false, "Base64 encode output")
	randCmd.Flags().StringVarP(&randOption.Out, "out", "", "", "Output file")

	pkeyutlCmd.Flags().BoolVarP(&pkeyutlOption.Encrypt, "encrypt", "", false, "Encrypt input data with public key")
	pkeyutlCmd.Flags().BoolVarP(&pkeyutlOption.Decrypt, "decrypt", "", false, "Decrypt input data with private key")
	pkeyutlCmd.Flags().BoolVarP(&pkeyutlOption.PubIn, "pubin", "", false, "Input key is a public key")
	pkeyutlCmd.Flags().StringVarP(&pkeyutlOption.In, "in", "", "", "Input file")
	pkeyutlCmd.Flags().StringVarP(&pkeyutlOption.InKey, "inkey", "", "", "Input key")
	pkeyutlCmd.Flags().StringVarP(&pkeyutlOption.Out, "out", "", "", "Output file")

	dgstCmd.Flags().BoolVarP(&dgstOption.Sha1, "sha1", "", false, "This specifies the message digest to sign the request")
	dgstCmd.Flags().StringVarP(&dgstOption.Out, "out", "", "", "Output to filename")
	dgstCmd.Flags().StringVarP(&dgstOption.Sign, "sign", "", "", "Sign digest using private key")
	dgstCmd.Flags().StringVarP(&dgstOption.Verify, "verify", "", "", "Verify a signature using public key")
	dgstCmd.Flags().StringVarP(&dgstOption.Signature, "signature", "", "", "File with signature to verify")

	encCmd.Flags().BoolVarP(&encOption.AES256CBC, "aes-256-cbc", "", false, "Use 256 bit AES cipher in CBC mode")
	encCmd.Flags().BoolVarP(&encOption.P, "p", "", false, "Print the iv/key")
	encCmd.Flags().BoolVarP(&encOption.D, "d", "", false, "Decrypt")
	encCmd.Flags().StringVarP(&encOption.Pass, "pass", "", "", "Passphrase source")
	encCmd.Flags().StringVarP(&encOption.Md, "md", "", "", "Use specified digest to create a key from the passphrase")
	encCmd.Flags().StringVarP(&encOption.In, "in", "", "", "Input file")
	encCmd.Flags().StringVarP(&encOption.Out, "out", "", "", "Output file")

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(genpkeyCmd)
	rootCmd.AddCommand(pkeyCmd)
	rootCmd.AddCommand(reqCmd)
	rootCmd.AddCommand(x509Cmd)
	rootCmd.AddCommand(verifyCmd)
	rootCmd.AddCommand(randCmd)
	rootCmd.AddCommand(pkeyutlCmd)
	rootCmd.AddCommand(dgstCmd)
	rootCmd.AddCommand(encCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of gossl",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Version is not set yet")
	},
}

var genpkeyCmd = &cobra.Command{
	Use:   "genpkey",
	Short: "The genpkey command generates a private key.",
	RunE: func(cmd *cobra.Command, args []string) error {
		return genpkey.Run(genpkeyOption)
	},
}

var pkeyCmd = &cobra.Command{
	Use:   "pkey",
	Short: "The pkey command processes public or private keys.",
	RunE: func(cmd *cobra.Command, args []string) error {
		return pkey.Run(pkeyOption)
	},
}

var reqCmd = &cobra.Command{
	Use:   "req",
	Short: "This command primarily creates and processes certificate requests (CSRs) in PKCS#10 format.",
	RunE: func(cmd *cobra.Command, args []string) error {
		return req.Run(reqOption)
	},
}

var x509Cmd = &cobra.Command{
	Use:   "x509",
	Short: "The x509 command is a multi purpose certificate utility.",
	RunE: func(cmd *cobra.Command, args []string) error {
		return x509.Run(x509Option)
	},
}

var verifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "The verify command verifies certificate chains.",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return verify.Run(verifyOption, args)
	},
}

var randCmd = &cobra.Command{
	Use:   "rand",
	Short: "This command generates num random bytes using a cryptographically secure pseudo random number generator (CSPRNG)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		i, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("failed to convert %s to int: %w", args[0], err)
		}
		return rand.Run(randOption, i)
	},
}

var pkeyutlCmd = &cobra.Command{
	Use:   "pkeyutl",
	Short: "The pkeyutl command can be used to perform low-level public key operations using any supported algorithm.",
	RunE: func(cmd *cobra.Command, args []string) error {
		return pkeyutl.Run(pkeyutlOption)
	},
}

var dgstCmd = &cobra.Command{
	Use:   "dgst",
	Short: "The digest functions output the message digest of a supplied file or files in hexadecimal. The digest functions also generate and verify digital signatures using message digests.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return dgst.Run(dgstOption, args[0])
	},
}

var encCmd = &cobra.Command{
	Use:   "enc",
	Short: "The symmetric cipher commands allow data to be encrypted or decrypted using various block and stream ciphers using keys based on passwords or explicitly provided.",
	RunE: func(cmd *cobra.Command, args []string) error {
		return enc.Run(encOption)
	},
}
