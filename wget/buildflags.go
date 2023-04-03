package wget

import (
	"flag"
	"os"
	"strings"
)

var Flags WgetFlags

type WgetFlags struct {
	H_Flag         bool   // flag for help
	O_Flag         string // flag for filename output
	Mirror_Flag    bool   // flag for mirroring website
	RateLimit_Flag string // flag for rame-limiting
	I_Flag         string // flag for downloading multiple files "./wget -i=downloads.txt"
	B_Flag         bool   // flag for logging output instead printing it out
	X_Flag         string // flag for excluding folders (ex. skipping/cutting out the /img folder --mirror https://trypap.com/)
	Reject_Flag    string // flag for excluding file types
	P_Flag         string   // flag for setting download directory
	Links          []string // links
}

func BuildFlags() {
	//boolean flags
	flag.BoolVar(&Flags.H_Flag, "h", false, "Helpful information for the user")
	flag.BoolVar(&Flags.Mirror_Flag, "mirror", false, "Mirror flag for CLONING webpages")
	flag.BoolVar(&Flags.B_Flag, "B", false, "BACKGROUND downloading")
	//string flags
	flag.StringVar(&Flags.O_Flag, "O", "", "Flag for downloaded file Name")
	flag.StringVar(&Flags.RateLimit_Flag, "rate-limit", "", "SET Maximum download speed for downloads")
	flag.StringVar(&Flags.I_Flag, "i", "", "Downloading MULTIPLE files from file")
	flag.StringVar(&Flags.X_Flag, "x", "", "EXCLUDE folders from being downloaded")
	flag.StringVar(&Flags.Reject_Flag, "reject", "", "EXCLUDE files from being downloaded")

	flag.StringVar(&Flags.P_Flag, "P", "", "SET download directory")
	flag.Parse() // parsing built flags to flags variable

	//links flag
	//if I_Flag tag used:
	if Flags.I_Flag != "" {
		file_data, err := os.ReadFile(Flags.I_Flag)
		errorHandler(err, true)

		file_array := strings.Split(string(file_data), "\n")
		for i := 0; i < len(file_array); i++ {
			if file_array[i] != "" {
				Flags.Links = append(Flags.Links, file_array[i])
			}
		}
	} else {
		Flags.Links = append(Flags.Links, flag.Arg(0))
	}

}
