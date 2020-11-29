package cmd

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"mq/bootstrapers"
	"os"
)

var app = &bootstrapers.App{
	Bootstrapers: []bootstrapers.Boot{
		&bootstrapers.ConfigLoader{},
		&bootstrapers.DBLoader{},
		&bootstrapers.RedisLoader{},
	},
}

var rootCmd = &cobra.Command{
	Use:   "app",
	Short: "mq event bus",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initLogger)
}

func initLogger() {
	log.SetLevel(log.InfoLevel)
	fmt.Print(`
          _____                   _______         
         /\    \                 /::\    \        
        /::\____\               /::::\    \       
       /::::|   |              /::::::\    \      
      /:::::|   |             /::::::::\    \     
     /::::::|   |            /:::/~~\:::\    \    
    /:::/|::|   |           /:::/    \:::\    \   
   /:::/ |::|   |          /:::/    / \:::\    \  
  /:::/  |::|___|______   /:::/____/   \:::\____\ 
 /:::/   |::::::::\    \ |:::|    |     |:::|    |
/:::/    |:::::::::\____\|:::|____|     |:::|____|
\::/    / ~~~~~/:::/    / \:::\   _\___/:::/    / 
 \/____/      /:::/    /   \:::\ |::| /:::/    /  
             /:::/    /     \:::\|::|/:::/    /   
            /:::/    /       \::::::::::/    /    
           /:::/    /         \::::::::/    /     
          /:::/    /           \::::::/    /      
         /:::/    /             \::::/____/       
        /:::/    /               |::|    |        
        \::/    /                |::|____|        
         \/____/                  ~~		@2020.11 by duc.
`)
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
}
