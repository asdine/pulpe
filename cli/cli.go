package cli

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/blankrobot/pulpe"
	"github.com/blankrobot/pulpe/http"
	"github.com/blankrobot/pulpe/http/api"
	"github.com/blankrobot/pulpe/mongo"
	"github.com/spf13/cobra"
)

// New returns the pulpe CLI application.
func New() *cobra.Command {
	cmd := cobra.Command{
		Use: "pulpe",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Usage()
		},
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	cmd.AddCommand(NewServerCmd())
	cmd.AddCommand(NewVersionCmd())
	return &cmd
}

// NewVersionCmd returns a command that displays the pulpe version number.
func NewVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:  "version",
		Long: "Display the version number",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println(pulpe.Version)
			return nil
		},
		SilenceErrors: true,
		SilenceUsage:  true,
	}
}

// NewServerCmd returns a ServerCmd.
func NewServerCmd() *cobra.Command {
	var s ServerCmd

	cmd := cobra.Command{
		Use:           "server",
		RunE:          s.Run,
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	cmd.Flags().StringVar(&s.addr, "http", ":4000", "HTTP address")
	cmd.Flags().StringVar(&s.mongoURI, "mongo", "mongodb://localhost:27017/pulpe", "MongoDB uri")
	cmd.Flags().StringVar(&s.assetsPath, "assets", "./web/build", "Assets directory")

	return &cmd
}

// ServerCmd is a command the runs the pulpe server.
type ServerCmd struct {
	addr       string
	mongoURI   string
	assetsPath string
}

// Run creates a bolt client and runs the HTTP server.
func (c *ServerCmd) Run(cmd *cobra.Command, args []string) error {
	client := mongo.NewClient(c.mongoURI)
	err := client.Open()
	if err != nil {
		return err
	}
	defer client.Close()

	client.Authenticator = new(mongo.Authenticator)

	connect := http.NewCookieConnector(client)

	mux := http.NewServeMux()

	api.Register(mux, connect)
	if c.assetsPath != "" {
		http.RegisterStaticHandler(mux, c.assetsPath)
		http.RegisterPageHandler(mux, c.assetsPath)
	}

	srv := http.NewServer(c.addr, mux)
	err = srv.Open()
	if err != nil {
		return err
	}

	log.Printf("Serving HTTP on address %s\n", c.addr)

	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)

	<-ch
	fmt.Println()
	log.Println("Stopping server...")
	err = srv.Close()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("OK")
	return nil
}
