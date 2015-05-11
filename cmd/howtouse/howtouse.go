package main

import (
	"flag"
	"github.com/jkl1337/go-mediainfo"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		log.Fatal("Not enough arguments.")
	}

	mi := mediainfo.New()
	if err := mi.Open(args[0]); err != nil {
		log.Fatal(err)
	}
	s := make([]string, 0)
	s = append(s, mi.Option("Info_Version"))
	s = append(s, "\nInfo_Parameters\n", mi.Option("Info_Parameters"))
	s = append(s, "\nInfo_Codecs\n", mi.Option("Info_Codecs"))

	s = append(s, "\nOpen\n")
	mi.SetOption("Complete", "")
	s = append(s, "\nInform with Complete=false\n", mi.Inform())
	mi.SetOption("Complete", "1")
	s = append(s, "\nInform with Complete=true\n", mi.Inform())

	mi.SetOption("Inform", "General;Example : FileSize=%FileSize%")
	s = append(s, "\nCustom Inform\n", mi.Inform())

	s = append(s, "\nGet with Stream=General and Parameter=\"FileSize\"\n",
		mi.GetKind(mediainfo.StreamGeneral, 0, "FileSize", mediainfo.InfoText))

	s = append(s, "\nCount with StreamKind=Audio\n",
		strconv.FormatInt(int64(mi.Count(mediainfo.StreamAudio)), 10))

	s = append(s, "\nGet with Stream=General and Parameter=46\n",
		mi.GetI(mediainfo.StreamGeneral, 0, 46, mediainfo.InfoText))

	s = append(s, "\nGet with Stream=General and Parameter=\"AudioCount\"\n",
		mi.GetKind(mediainfo.StreamGeneral, 0, "AudioCount", mediainfo.InfoText))

	s = append(s, "\nGet with Stream=Audio and Parameter=\"StreamCount\"\n",
		mi.GetKind(mediainfo.StreamAudio, 0, "StreamCount", mediainfo.InfoText))

	s = append(s, "\nClose\n")
	mi.Close()

	os.Stdout.WriteString(strings.Join(s, "\n"))
}
